import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';
import Stripe from 'stripe';
import { onCall } from 'firebase-functions/v2/https';

const stripe = new Stripe(functions.config().stripe.secret_key || 'sk_test_demo', {
  apiVersion: '2025-09-30.clover'
});

/**
 * Create a payment intent for auction bidding
 */
export const createPaymentIntent = onCall({cors: true}, async (request) => {
  if (!request.auth) {
    throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
  }

  const { amount, currency = 'usd', auctionId, bidId } = request.data;

  try {
    // Get user data
    const userDoc = await admin.firestore().collection('users').doc(request.auth.uid).get();
    if (!userDoc.exists) {
      throw new functions.https.HttpsError('not-found', 'User not found');
    }

    const userData = userDoc.data();

    if (!userData) {
      throw new functions.https.HttpsError('not-found', 'User data not found');
    }

    // Create payment intent
    const paymentIntent = await stripe.paymentIntents.create({
      amount: Math.round(amount * 100), // Convert to cents
      currency,
      customer: userData.stripeCustomerId,
      metadata: {
        userId: request.auth.uid,
        auctionId,
        bidId,
        type: 'auction_bid'
      },
      automatic_payment_methods: {
        enabled: true
      }
    });

    // Store payment intent in Firestore
    await admin.firestore().collection('payments').doc(paymentIntent.id).set({
      userId: request.auth.uid,
      auctionId,
      bidId,
      amount,
      currency,
      stripePaymentIntentId: paymentIntent.id,
      status: 'requires_payment_method',
      createdAt: admin.firestore.FieldValue.serverTimestamp(),
      updatedAt: admin.firestore.FieldValue.serverTimestamp()
    });

    return {
      success: true,
      clientSecret: paymentIntent.client_secret,
      paymentIntentId: paymentIntent.id
    };
  } catch (error) {
    console.error('Error creating payment intent:', error);
    throw new functions.https.HttpsError('internal', 'Failed to create payment intent');
  }
});

/**
 * Confirm payment and update bid status
 */
export const confirmPayment = onCall({cors: true}, async (request) => {
  if (!request.auth) {
    throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
  }

  const { paymentIntentId, auctionId, bidId } = request.data;

  try {
    // Retrieve payment intent
    const paymentIntent = await stripe.paymentIntents.retrieve(paymentIntentId);

    if (paymentIntent.status === 'succeeded') {
      // Update payment record
      await admin.firestore().collection('payments').doc(paymentIntentId).update({
        status: 'succeeded',
        updatedAt: admin.firestore.FieldValue.serverTimestamp()
      });

      // Update bid status
      if (auctionId && bidId) {
        await admin.firestore()
          .collection('auctions')
          .doc(auctionId)
          .collection('bids')
          .doc(bidId)
          .update({
            paymentStatus: 'paid',
            paidAt: admin.firestore.FieldValue.serverTimestamp(),
            updatedAt: admin.firestore.FieldValue.serverTimestamp()
          });
      }

      // Update user wallet balance
      await admin.firestore().collection('users').doc(request.auth.uid).update({
        walletBalance: admin.firestore.FieldValue.increment(paymentIntent.amount / 100),
        updatedAt: admin.firestore.FieldValue.serverTimestamp()
      });
    }

    return {
      success: true,
      status: paymentIntent.status,
      message: 'Payment processed successfully'
    };
  } catch (error) {
    console.error('Error confirming payment:', error);
    throw new functions.https.HttpsError('internal', 'Failed to confirm payment');
  }
});

/**
 * Handle Stripe webhook events
 */
export const stripeWebhook = functions.https.onRequest(async (req, res) => {
  const sig = req.headers['stripe-signature'];
  const endpointSecret = functions.config().stripe.webhook_secret || 'whsec_demo';

  let event;

  try {
    if (!sig) {
      res.status(400).send('No signature provided');
      return;
    }

    const body = req.rawBody;
    if (!body) {
      res.status(400).send('No body provided');
      return;
    }

    event = stripe.webhooks.constructEvent(body as unknown as string, sig as string, endpointSecret);
  } catch (err) {
    console.error('Webhook signature verification failed:', err);
    res.status(400).send('Invalid signature');
    return;
  }

  try {
    switch (event.type) {
      case 'payment_intent.succeeded':
        const paymentIntent = event.data.object as any;
        console.log('PaymentIntent was successful:', paymentIntent.id);
        // Handle successful payment
        break;

      case 'payment_intent.payment_failed':
        const failedPayment = event.data.object as any;
        console.log('PaymentIntent failed:', failedPayment.id);
        // Handle failed payment
        break;

      default:
        console.log('Unhandled event type:', event.type);
    }

    res.json({ received: true });
  } catch (error) {
    console.error('Webhook handler error:', error);
    res.status(500).send('Internal server error');
  }
});