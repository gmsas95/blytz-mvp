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
exports.sendAuctionUpdate = exports.sendNotification = void 0;
const functions = __importStar(require("firebase-functions"));
const admin = __importStar(require("firebase-admin"));
/**
 * Send push notification to user
 */
exports.sendNotification = functions.https.onCall(async (data, context) => {
    if (!context.auth) {
        throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
    }
    const { userId, title, body, data: notificationData = {} } = data;
    try {
        // Get user's FCM token
        const userDoc = await admin.firestore().collection('users').doc(userId).get();
        const userData = userDoc.data();
        if (!(userData === null || userData === void 0 ? void 0 : userData.fcmToken)) {
            throw new functions.https.HttpsError('failed-precondition', 'User has no FCM token');
        }
        const message = {
            notification: {
                title,
                body
            },
            data: Object.assign(Object.assign({}, data), { timestamp: Date.now().toString() }),
            token: userData.fcmToken
        };
        const response = await admin.messaging().send(message);
        // Store notification in Firestore
        await admin.firestore().collection('notifications').add({
            userId,
            title,
            body,
            data: notificationData,
            status: 'sent',
            messageId: response,
            createdAt: admin.firestore.FieldValue.serverTimestamp()
        });
        return {
            success: true,
            messageId: response,
            message: 'Notification sent successfully'
        };
    }
    catch (error) {
        console.error('Error sending notification:', error);
        throw new functions.https.HttpsError('internal', 'Failed to send notification');
    }
});
/**
 * Send auction update notification
 */
exports.sendAuctionUpdate = functions.https.onCall(async (data, context) => {
    const { auctionId, type, message } = data;
    try {
        // Get auction participants
        const participantsSnapshot = await admin.firestore()
            .collection('auctions')
            .doc(auctionId)
            .collection('participants')
            .get();
        const notifications = [];
        for (const participantDoc of participantsSnapshot.docs) {
            const participant = participantDoc.data();
            if (participant.fcmToken) {
                const notification = {
                    notification: {
                        title: 'Auction Update',
                        body: message
                    },
                    data: {
                        auctionId,
                        type,
                        timestamp: Date.now().toString()
                    },
                    token: participant.fcmToken
                };
                try {
                    const response = await admin.messaging().send(notification);
                    notifications.push({
                        userId: participant.userId,
                        messageId: response,
                        status: 'sent'
                    });
                }
                catch (error) {
                    console.error('Failed to send notification to user:', participant.userId, error);
                    notifications.push({
                        userId: participant.userId,
                        status: 'failed',
                        error: error.message
                    });
                }
            }
        }
        return {
            success: true,
            notificationsSent: notifications.length,
            notifications: notifications
        };
    }
    catch (error) {
        console.error('Error sending auction update:', error);
        throw new functions.https.HttpsError('internal', 'Failed to send auction update');
    }
});
//# sourceMappingURL=notifications.js.map