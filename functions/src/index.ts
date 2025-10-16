import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';

// Initialize Firebase Admin
admin.initializeApp();

// Export functions
export * from './auth';
export * from './payments';
export * from './notifications';
export * from './auction';

// Health check function
export const health = functions.https.onCall(async (data, context) => {
  return {
    status: 'healthy',
    timestamp: new Date().toISOString(),
    service: 'blytz-firebase-functions',
    version: '1.0.0'
  };
});