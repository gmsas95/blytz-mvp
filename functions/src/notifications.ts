import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';

/**
 * Send push notification to user
 */
export const sendNotification = functions.https.onCall(async (data, context) => {
  if (!context.auth) {
    throw new functions.https.HttpsError('unauthenticated', 'User must be authenticated');
  }

  const { userId, title, body, data: notificationData = {} } = data;

  try {
    // Get user's FCM token
    const userDoc = await admin.firestore().collection('users').doc(userId).get();
    const userData = userDoc.data();

    if (!userData?.fcmToken) {
      throw new functions.https.HttpsError('failed-precondition', 'User has no FCM token');
    }

    const message = {
      notification: {
        title,
        body
      },
      data: {
        ...data,
        timestamp: Date.now().toString()
      },
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
  } catch (error: any) {
    console.error('Error sending notification:', error);
    throw new functions.https.HttpsError('internal', 'Failed to send notification');
  }
});

/**
 * Send auction update notification
 */
export const sendAuctionUpdate = functions.https.onCall(async (data, context) => {
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
        } catch (error: any) {
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
  } catch (error: any) {
    console.error('Error sending auction update:', error);
    throw new functions.https.HttpsError('internal', 'Failed to send auction update');
  }
});