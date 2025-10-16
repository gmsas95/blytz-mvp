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
exports.updateProfile = exports.validateToken = exports.createUser = void 0;
const functions = __importStar(require("firebase-functions"));
const admin = __importStar(require("firebase-admin"));
/**
 * Create a new user account with custom claims for auction platform
 */
exports.createUser = functions.https.onCall(async (data, context) => {
    const { email, password, displayName, phoneNumber } = data;
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
    }
    catch (error) {
        console.error('Error creating user:', error);
        throw new functions.https.HttpsError('internal', 'Failed to create user');
    }
});
/**
 * Validate Firebase ID token and return user data
 */
exports.validateToken = functions.https.onCall(async (data, context) => {
    if (!context.auth) {
        throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
    }
    try {
        const userDoc = await admin.firestore().collection('users').doc(context.auth.uid).get();
        if (!userDoc.exists) {
            throw new functions.https.HttpsError('not-found', 'User not found');
        }
        return {
            success: true,
            user: Object.assign(Object.assign({ uid: context.auth.uid }, userDoc.data()), { customClaims: context.auth.token })
        };
    }
    catch (error) {
        console.error('Error validating token:', error);
        throw new functions.https.HttpsError('internal', 'Failed to validate token');
    }
});
/**
 * Update user profile
 */
exports.updateProfile = functions.https.onCall(async (data, context) => {
    if (!context.auth) {
        throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
    }
    const { displayName, phoneNumber } = data;
    try {
        // Update Firebase Auth
        await admin.auth().updateUser(context.auth.uid, {
            displayName,
            phoneNumber
        });
        // Update Firestore
        await admin.firestore().collection('users').doc(context.auth.uid).update({
            displayName,
            phoneNumber,
            updatedAt: admin.firestore.FieldValue.serverTimestamp()
        });
        return {
            success: true,
            message: 'Profile updated successfully'
        };
    }
    catch (error) {
        console.error('Error updating profile:', error);
        throw new functions.https.HttpsError('internal', 'Failed to update profile');
    }
});
//# sourceMappingURL=auth.js.map