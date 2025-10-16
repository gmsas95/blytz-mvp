"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.stripeWebhook = exports.confirmPayment = exports.createPaymentIntent = void 0;
const functions = __importStar(require("firebase-functions"));
const admin = __importStar(require("firebase-admin"));
const stripe_1 = __importDefault(require("stripe"));
const https_1 = require("firebase-functions/v2/https");
const stripe = new stripe_1.default(functions.config().stripe.secret_key || 'sk_test_demo', {
    apiVersion: '2025-09-30.clover'
});
/**
 * Create a payment intent for auction bidding
 */
exports.createPaymentIntent = (0, https_1.onCall)({ cors: true }, async (request) => {
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
    }
    catch (error) {
        console.error('Error creating payment intent:', error);
        throw new functions.https.HttpsError('internal', 'Failed to create payment intent');
    }
});
/**
 * Confirm payment and update bid status
 */
exports.confirmPayment = (0, https_1.onCall)({ cors: true }, async (request) => {
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
    }
    catch (error) {
        console.error('Error confirming payment:', error);
        throw new functions.https.HttpsError('internal', 'Failed to confirm payment');
    }
});
/**
 * Handle Stripe webhook events
 */
exports.stripeWebhook = functions.https.onRequest(async (req, res) => {
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
        event = stripe.webhooks.constructEvent(body, sig, endpointSecret);
    }
    catch (err) {
        console.error('Webhook signature verification failed:', err);
        res.status(400).send('Invalid signature');
        return;
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
    }
    catch (error) {
        console.error('Webhook handler error:', error);
        res.status(500).send('Internal server error');
    }
});
//# sourceMappingURL=payments.js.map