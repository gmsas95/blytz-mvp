import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';
import Stripe from 'stripe';

const stripe = new Stripe(functions.config().stripe.secret_key || 'sk_test_demo', {
  apiVersion: '2023-10-16'
});

/**
 * Create a payment intent for auction bidding
 */
export const createPaymentIntent = functions.https.onCall(async (data, context) => {
  if (!context.auth) {
    throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
  }

  const { amount, currency = 'usd', auctionId, bidId } = data;

  try {
    // Get user data
    const userDoc = await admin.firestore().collection('users').doc(context.auth.uid).get();
    if (!userDoc.exists) {
      throw new functions.https.HttpsError('not-found', 'User not found');
    }

    const userData = userDoc.data();

    // Create payment intent
    const paymentIntent = await stripe.paymentIntents.create({
      amount: Math.round(amount * 100), // Convert to cents
      currency,
      customer: userData.stripeCustomerId,
      metadata: {
        userId: context.auth.uid,
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
      userId: context.auth.uid,
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
export const confirmPayment = functions.https.onCall(async (data, context) => {
  if (!context.auth) {
    throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
  }

  const { paymentIntentId, auctionId, bidId } = data;

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
      await admin.firestore().collection('users').doc(context.auth.uid).update({
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
    event = stripe.webhooks.constructEvent(req.rawBody, sig, endpointSecret);
  } catch (err) {
    console.error('Webhook signature verification failed:', err);
    return res.status(400).send('Invalid signature');
  }

  try {
    switch (event.type) {
      case 'payment_intent.succeeded':
        const paymentIntent = event.data.object;
        console.log('PaymentIntent was successful:', paymentIntent.id);
        // Handle successful payment
        break;

      case 'payment_intent.payment_failed':
        const failedPayment = event.data.object;
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