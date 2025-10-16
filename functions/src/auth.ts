import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';
import { onCall } from 'firebase-functions/v2/https';

/**
 * Create a new user account with custom claims for auction platform
 */
export const createUser = onCall({cors: true}, async (request) => {
  const { email, password, displayName, phoneNumber } = request.data;

  try {
    // Create user in Firebase Auth
    const userRecord = await admin.auth().createUser({
      email,
      password,
      displayName,
      phoneNumber
    });

    // Set custom claims for auction platform
    await admin.auth().setCustomUserClaims(userRecord.uid, {
      role: 'user',
      canBid: true,
      canHost: false,
      createdAt: Date.now()
    });

    // Create user document in Firestore
    await admin.firestore().collection('users').doc(userRecord.uid).set({
      email,
      displayName,
      phoneNumber,
      role: 'user',
      walletBalance: 0,
      totalBids: 0,
      totalAuctions: 0,
      createdAt: admin.firestore.FieldValue.serverTimestamp(),
      updatedAt: admin.firestore.FieldValue.serverTimestamp()
    });

    return {
      success: true,
      uid: userRecord.uid,
      message: 'User created successfully'
    };
  } catch (error) {
    console.error('Error creating user:', error);
    throw new functions.https.HttpsError('internal', 'Failed to create user');
  }
});

/**
 * Validate Firebase ID token and return user data
 */
export const validateToken = onCall({cors: true}, async (request) => {
  if (!request.auth) {
    throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
  }

  try {
    const userDoc = await admin.firestore().collection('users').doc(request.auth.uid).get();

    if (!userDoc.exists) {
      throw new functions.https.HttpsError('not-found', 'User not found');
    }

    return {
      success: true,
      user: {
        uid: request.auth.uid,
        ...userDoc.data(),
        customClaims: request.auth.token
      }
    };
  } catch (error) {
    console.error('Error validating token:', error);
    throw new functions.https.HttpsError('internal', 'Failed to validate token');
  }
});

/**
 * Update user profile
 */
export const updateProfile = onCall({cors: true}, async (request) => {
  if (!request.auth) {
    throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
  }

  const { displayName, phoneNumber } = request.data;

  try {
    // Update Firebase Auth
    await admin.auth().updateUser(request.auth.uid, {
      displayName,
      phoneNumber
    });

    // Update Firestore
    await admin.firestore().collection('users').doc(request.auth.uid).update({
      displayName,
      phoneNumber,
      updatedAt: admin.firestore.FieldValue.serverTimestamp()
    });

    return {
      success: true,
      message: 'Profile updated successfully'
    };
  } catch (error) {
    console.error('Error updating profile:', error);
    throw new functions.https.HttpsError('internal', 'Failed to update profile');
  }
});