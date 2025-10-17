import * as admin from 'firebase-admin';

// Initialize Firebase Admin
admin.initializeApp();

// Export functions
export * from './auth';
export * from './payments';
export * from './notifications';
export * from './auction';

// Health check function
import { onCall } from 'firebase-functions/v2/https';

export const health = onCall({cors: true}, async (request) => {
  return {
    status: 'healthy',
    timestamp: new Date().toISOString(),
    service: 'blytz-firebase-functions',
    version: '2.0.0'
  };
});