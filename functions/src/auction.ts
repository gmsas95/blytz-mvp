import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';
import { onCall } from 'firebase-functions/v2/https';

/**
 * Create a new auction
 */
export const createAuction = onCall({cors: true}, async (request) => {
  if (!request.auth) {
    throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
  }

  const { title, description, startingPrice, duration, category, images } = request.data;

  try {
    // Validate user can create auctions
    const userDoc = await admin.firestore().collection('users').doc(request.auth.uid).get();
    const userData = userDoc.data();

    if (!userData?.canHost) {
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
  } catch (error) {
    console.error('Error creating auction:', error);
    throw new functions.https.HttpsError('internal', 'Failed to create auction');
  }
});

/**
 * Place a bid on an auction
 */
export const placeBid = onCall({cors: true}, async (request) => {
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
    if (auction?.status !== 'active') {
      throw new functions.https.HttpsError('failed-precondition', 'Auction is not active');
    }

    // Check if auction has ended
    if (auction?.endTime?.toDate() < new Date()) {
      throw new functions.https.HttpsError('failed-precondition', 'Auction has ended');
    }

    // Validate bid amount
    if (amount <= auction?.currentPrice) {
      throw new functions.https.HttpsError('failed-precondition', 'Bid must be higher than current price');
    }

    // Check user's wallet balance
    const userDoc = await admin.firestore().collection('users').doc(request.auth.uid).get();
    const userData = userDoc.data();

    if (userData?.walletBalance < amount) {
      throw new functions.https.HttpsError('failed-precondition', 'Insufficient wallet balance');
    }

    const bidId = admin.firestore().collection('auctions').doc(auctionId).collection('bids').doc().id;
    const bidData = {
      id: bidId,
      auctionId,
      userId: request.auth.uid,
      userName: userData?.displayName,
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
  } catch (error) {
    console.error('Error placing bid:', error);
    throw new functions.https.HttpsError('internal', 'Failed to place bid');
  }
});

/**
 * End an auction and declare winner
 */
export const endAuction = onCall({cors: true}, async (request) => {
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
    if (auction?.hostId !== request.auth.uid) {
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
      winningBid: winningBid?.amount || auction?.startingPrice,
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
      winningBid: winningBid?.amount || auction?.startingPrice,
      message: winnerId ? 'Auction ended with winner' : 'Auction ended with no bids'
    };
  } catch (error) {
    console.error('Error ending auction:', error);
    throw new functions.https.HttpsError('internal', 'Failed to end auction');
  }
});

/**
 * Get auction details with bids
 */
export const getAuctionDetails = onCall({cors: true}, async (request) => {
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

    const bids = bidsSnapshot.docs.map(doc => ({
      id: doc.id,
      ...doc.data()
    }));

    return {
      success: true,
      auction: {
        ...auction,
        id: auctionId
      },
      bids,
      totalBids: bids.length
    };
  } catch (error) {
    console.error('Error getting auction details:', error);
    throw new functions.https.HttpsError('internal', 'Failed to get auction details');
  }
});