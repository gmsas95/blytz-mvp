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
Object.defineProperty(exports, "__esModule", { value: true });
exports.getAuctionDetails = exports.endAuction = exports.placeBid = exports.createAuction = void 0;
const functions = __importStar(require("firebase-functions"));
const admin = __importStar(require("firebase-admin"));
const https_1 = require("firebase-functions/v2/https");
/**
 * Create a new auction
 */
exports.createAuction = (0, https_1.onCall)({ cors: true }, async (request) => {
    if (!request.auth) {
        throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
    }
    const { title, description, startingPrice, duration, category, images } = request.data;
    try {
        // Validate user can create auctions
        const userDoc = await admin.firestore().collection('users').doc(request.auth.uid).get();
        const userData = userDoc.data();
        if (!(userData === null || userData === void 0 ? void 0 : userData.canHost)) {
            throw new functions.https.HttpsError('permission-denied', 'User cannot host auctions');
        }
        const auctionId = admin.firestore().collection('auctions').doc().id;
        const endTime = new Date(Date.now() + duration * 60 * 60 * 1000); // duration in hours
        const auctionData = {
            id: auctionId,
            title,
            description,
            startingPrice,
            currentPrice: startingPrice,
            duration,
            category,
            images: images || [],
            hostId: request.auth.uid,
            hostName: userData.displayName,
            status: 'active',
            startTime: admin.firestore.FieldValue.serverTimestamp(),
            endTime: endTime,
            bidCount: 0,
            participantCount: 0,
            winnerId: null,
            createdAt: admin.firestore.FieldValue.serverTimestamp(),
            updatedAt: admin.firestore.FieldValue.serverTimestamp()
        };
        await admin.firestore().collection('auctions').doc(auctionId).set(auctionData);
        return {
            success: true,
            auctionId,
            message: 'Auction created successfully'
        };
    }
    catch (error) {
        console.error('Error creating auction:', error);
        throw new functions.https.HttpsError('internal', 'Failed to create auction');
    }
});
/**
 * Place a bid on an auction
 */
exports.placeBid = (0, https_1.onCall)({ cors: true }, async (request) => {
    var _a;
    if (!request.auth) {
        throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
    }
    const { auctionId, amount } = request.data;
    try {
        const auctionRef = admin.firestore().collection('auctions').doc(auctionId);
        const auctionDoc = await auctionRef.get();
        if (!auctionDoc.exists) {
            throw new functions.https.HttpsError('not-found', 'Auction not found');
        }
        const auction = auctionDoc.data();
        // Validate auction status
        if ((auction === null || auction === void 0 ? void 0 : auction.status) !== 'active') {
            throw new functions.https.HttpsError('failed-precondition', 'Auction is not active');
        }
        // Check if auction has ended
        if (((_a = auction === null || auction === void 0 ? void 0 : auction.endTime) === null || _a === void 0 ? void 0 : _a.toDate()) < new Date()) {
            throw new functions.https.HttpsError('failed-precondition', 'Auction has ended');
        }
        // Validate bid amount
        if (amount <= (auction === null || auction === void 0 ? void 0 : auction.currentPrice)) {
            throw new functions.https.HttpsError('failed-precondition', 'Bid must be higher than current price');
        }
        // Check user's wallet balance
        const userDoc = await admin.firestore().collection('users').doc(request.auth.uid).get();
        const userData = userDoc.data();
        if ((userData === null || userData === void 0 ? void 0 : userData.walletBalance) < amount) {
            throw new functions.https.HttpsError('failed-precondition', 'Insufficient wallet balance');
        }
        const bidId = admin.firestore().collection('auctions').doc(auctionId).collection('bids').doc().id;
        const bidData = {
            id: bidId,
            auctionId,
            userId: request.auth.uid,
            userName: userData === null || userData === void 0 ? void 0 : userData.displayName,
            amount,
            timestamp: admin.firestore.FieldValue.serverTimestamp(),
            status: 'active'
        };
        // Create bid document
        await admin.firestore()
            .collection('auctions')
            .doc(auctionId)
            .collection('bids')
            .doc(bidId)
            .set(bidData);
        // Update auction with new bid
        await auctionRef.update({
            currentPrice: amount,
            bidCount: admin.firestore.FieldValue.increment(1),
            lastBidAt: admin.firestore.FieldValue.serverTimestamp(),
            updatedAt: admin.firestore.FieldValue.serverTimestamp()
        });
        // Update user bid count
        await admin.firestore().collection('users').doc(request.auth.uid).update({
            totalBids: admin.firestore.FieldValue.increment(1),
            updatedAt: admin.firestore.FieldValue.serverTimestamp()
        });
        return {
            success: true,
            bidId,
            message: 'Bid placed successfully'
        };
    }
    catch (error) {
        console.error('Error placing bid:', error);
        throw new functions.https.HttpsError('internal', 'Failed to place bid');
    }
});
/**
 * End an auction and declare winner
 */
exports.endAuction = (0, https_1.onCall)({ cors: true }, async (request) => {
    if (!request.auth) {
        throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
    }
    const { auctionId } = request.data;
    try {
        const auctionRef = admin.firestore().collection('auctions').doc(auctionId);
        const auctionDoc = await auctionRef.get();
        if (!auctionDoc.exists) {
            throw new functions.https.HttpsError('not-found', 'Auction not found');
        }
        const auction = auctionDoc.data();
        // Verify user is the host
        if ((auction === null || auction === void 0 ? void 0 : auction.hostId) !== request.auth.uid) {
            throw new functions.https.HttpsError('permission-denied', 'Only auction host can end the auction');
        }
        // Get the highest bid
        const bidsSnapshot = await admin.firestore()
            .collection('auctions')
            .doc(auctionId)
            .collection('bids')
            .orderBy('amount', 'desc')
            .limit(1)
            .get();
        let winnerId = null;
        let winningBid = null;
        if (!bidsSnapshot.empty) {
            const winningBidDoc = bidsSnapshot.docs[0];
            winningBid = winningBidDoc.data();
            winnerId = winningBid.userId;
        }
        // Update auction status
        await auctionRef.update({
            status: 'ended',
            winnerId: winnerId,
            winningBid: (winningBid === null || winningBid === void 0 ? void 0 : winningBid.amount) || (auction === null || auction === void 0 ? void 0 : auction.startingPrice),
            endedAt: admin.firestore.FieldValue.serverTimestamp(),
            updatedAt: admin.firestore.FieldValue.serverTimestamp()
        });
        // Update winner's auction count
        if (winnerId) {
            await admin.firestore().collection('users').doc(winnerId).update({
                totalAuctions: admin.firestore.FieldValue.increment(1),
                updatedAt: admin.firestore.FieldValue.serverTimestamp()
            });
        }
        return {
            success: true,
            winnerId,
            winningBid: (winningBid === null || winningBid === void 0 ? void 0 : winningBid.amount) || (auction === null || auction === void 0 ? void 0 : auction.startingPrice),
            message: winnerId ? 'Auction ended with winner' : 'Auction ended with no bids'
        };
    }
    catch (error) {
        console.error('Error ending auction:', error);
        throw new functions.https.HttpsError('internal', 'Failed to end auction');
    }
});
/**
 * Get auction details with bids
 */
exports.getAuctionDetails = (0, https_1.onCall)({ cors: true }, async (request) => {
    const { auctionId } = request.data;
    try {
        const auctionRef = admin.firestore().collection('auctions').doc(auctionId);
        const auctionDoc = await auctionRef.get();
        if (!auctionDoc.exists) {
            throw new functions.https.HttpsError('not-found', 'Auction not found');
        }
        const auction = auctionDoc.data();
        // Get recent bids
        const bidsSnapshot = await admin.firestore()
            .collection('auctions')
            .doc(auctionId)
            .collection('bids')
            .orderBy('timestamp', 'desc')
            .limit(50)
            .get();
        const bids = bidsSnapshot.docs.map(doc => (Object.assign({ id: doc.id }, doc.data())));
        return {
            success: true,
            auction: Object.assign(Object.assign({}, auction), { id: auctionId }),
            bids,
            totalBids: bids.length
        };
    }
    catch (error) {
        console.error('Error getting auction details:', error);
        throw new functions.https.HttpsError('internal', 'Failed to get auction details');
    }
});
//# sourceMappingURL=auction.js.map