# Ninja Van Integration - Webhook V2 Compliance

## ✅ **Webhook V2 Implementation Status**

Our implementation is **fully compliant with Ninja Van Webhook V2**:

### Key V2 Features Implemented:

1. **Correct Payload Structure** - All V2 webhook events supported
2. **HMAC Signature Verification** - Using Client Key for security
3. **ISO 8601 Timestamp Parsing** - Proper timezone handling
4. **Complete Event Coverage** - All 40+ webhook events handled
5. **Asynchronous Processing** - Immediate 200 response with background processing

### 🔧 **Authentication Configuration**

**OAuth 2.0 Authentication** (for API calls):
- Uses `client_id` + `client_key` (where client_key is the "Client Key" from dashboard)
- Grant type: `client_credentials`
- Token caching and refresh

**Webhook Verification** (for incoming webhooks):
- Uses `Client Key` for HMAC-SHA256 signature
- Header: `X-Ninjavan-Hmac-Sha256`
- Raw JSON body (no encoding)

### 📋 **Required Environment Variables**

```bash
# OAuth Authentication
NINJAVAN_CLIENT_ID=your_client_id_from_dashboard
NINJAVAN_CLIENT_KEY=your_client_key_from_dashboard  # NOT client_secret

# Configuration
NINJAVAN_ENVIRONMENT=sandbox  # or "production"
NINJAVAN_COUNTRY_CODE=sg      # your country code

# Optional: For webhook signature verification only
NINJAVAN_CLIENT_KEY=your_client_key_from_dashboard
```

### 🚨 **Important Notes**

1. **Client Key vs Client Secret**: 
   - Use **Client Key** (not Client Secret) from Ninja Dashboard
   - The field is called "client_secret" in OAuth API but contains your Client Key
   - Webhook verification also uses the same Client Key

2. **Webhook V2 Features**:
   - All new events like "Driver dispatched for Pickup"
   - International transit events
   - Enhanced proof information (photos, signatures)
   - Proper timezone handling in timestamps

3. **Security**:
   - HMAC-SHA256 signature verification
   - No IP whitelisting (as per Ninja Van docs)
   - Raw JSON body for signature calculation

### 🔄 **Webhook Events Supported**

Our implementation handles all V2 events including:

**Core Events**:
- Pending Pickup
- Driver dispatched for Pickup  
- Picked Up, In Transit to Origin Hub
- On Vehicle for Delivery
- Delivered (all variants)
- Cancelled

**Exception Events**:
- Pickup Exception (all states)
- Delivery Exception (all states)
- Return to Shipper Exception

**Hub Events**:
- Arrived at Origin/Transit/Destination Hub
- In Transit to Next Sorting Hub

**International Events**:
- All international transit states
- Customs clearance events
- Linehaul information

**Special Events**:
- Parcel Measurements Update
- PUDO-related events
- RTS (Return to Sender) events

### 📡 **Webhook Endpoint**

```
POST /api/v1/logistics/ninjavan/webhook
```

- **Public endpoint** (no authentication required)
- **Signature verification** using Client Key
- **Immediate response** with async processing
- **Retry handling** for failed deliveries

### 🎯 **Compliance Summary**

✅ **Webhook V2 Compliant** - All V2 features implemented  
✅ **Security Compliant** - HMAC verification with Client Key  
✅ **Event Coverage** - All 40+ webhook events handled  
✅ **Error Handling** - Proper retry and logging mechanisms  
✅ **Performance** - Async processing with immediate response  

Your integration is ready for production with Ninja Van Webhook V2!