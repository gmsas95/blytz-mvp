# Fiuu Payment Gateway Integration

This document describes the integration of Fiuu payment gateway into the Blytz auction platform.

## Overview

Fiuu is a leading Southeast Asian payment gateway that supports 110+ payment methods across Malaysia, Singapore, Philippines, and other SEA countries. The integration provides:

- **Multi-channel support**: FPX, credit cards, e-wallets (GrabPay, Touch 'n Go, ShopeePay), Buy Now Pay Later
- **Seamless integration**: No redirects to external payment pages
- **Real-time processing**: Instant payment confirmations via webhooks
- **Multi-currency**: Support for MYR, SGD, PHP, THB, USD, etc.

## Architecture

### Backend Integration

The payment service now includes:

1. **Fiuu Client Library** (`pkg/fiuu/`)
   - `client.go`: HTTP client for Fiuu API
   - `types.go`: Request/response types and channel definitions

2. **Updated Payment Service** (`internal/services/payment.go`)
   - Real Fiuu payment processing
   - Webhook handling for payment confirmations
   - Seamless configuration generation

3. **Enhanced API Endpoints**
   - `GET /api/v1/payments/seamless/config`: Get frontend configuration
   - `POST /api/v1/webhooks/fiuu`: Handle Fiuu webhooks

### Frontend Integration

The frontend uses Fiuu's JavaScript library for seamless integration:

```html
<script src="https://pay.fiuu.com/RMS/API/seamless/3.28/js/MOLPay_seamless.deco.js"></script>
```

## Configuration

### Environment Variables

Add these to your payment service environment:

```bash
# Fiuu Payment Gateway Configuration
FIUU_MERCHANT_ID=your_merchant_id
FIUU_VERIFY_KEY=your_verify_key
FIUU_SANDBOX=true
FIUU_RETURN_URL=https://yourdomain.com/payment/return
FIUU_NOTIFY_URL=https://yourdomain.com/api/v1/webhooks/fiuu
FIUU_CALLBACK_URL=https://yourdomain.com/payment/callback
FIUU_CANCEL_URL=https://yourdomain.com/payment/cancel
```

### Getting Fiuu Credentials

1. **Register for Fiuu Account**
   - Visit [Fiuu](https://fiuu.com) and sign up
   - Email support@fiuu.com for merchant account setup

2. **Get Credentials**
   - Merchant ID: Your unique merchant identifier
   - Verify Key: Secret key for generating vcode (verification codes)

3. **Domain Registration**
   - Register your domain with Fiuu support
   - Required for seamless integration

## Supported Payment Channels

### Malaysia (MYR)
- **FPX**: Online banking (Maybank2u, CIMB Clicks, HLB Connect, etc.)
- **Credit/Debit Cards**: Visa, Mastercard, AMEX
- **E-Wallets**: GrabPay, Touch 'n Go, ShopeePay, Boost
- **QR Payments**: DuitNow QR
- **BNPL**: Atome, Rely

### Singapore (SGD)
- **Online Banking**: DBS, OCBC, UOB via eNETS
- **PayNow**: QR code payments

### Philippines (PHP)
- **E-Wallets**: GCash, Maya
- **Online Banking**: BPI, UnionBank

## API Usage

### 1. Get Payment Methods

```bash
GET /api/v1/payments/methods
Authorization: Bearer <jwt_token>
```

Response:
```json
{
  "methods": [
    {
      "id": "fpx",
      "name": "FPX Online Banking",
      "type": "online_banking",
      "description": "Pay with Malaysian online banking",
      "enabled": true,
      "channel": "fpx",
      "currency": "MYR"
    },
    {
      "id": "GrabPay",
      "name": "GrabPay",
      "type": "ewallet",
      "description": "Pay with GrabPay e-wallet",
      "enabled": true,
      "channel": "GrabPay",
      "currency": "MYR"
    }
  ]
}
```

### 2. Get Seamless Configuration

```bash
GET /api/v1/payments/seamless/config?order_id=ORDER123&amount=5000&bill_name=John%20Doe&bill_email=john@example.com&bill_mobile=01234567890&bill_desc=Payment&channel=fpx
Authorization: Bearer <jwt_token>
```

Response:
```json
{
  "mpsmerchantid": "your_merchant_id",
  "mpschannel": "fpx",
  "mpsamount": "50.00",
  "mpsorderid": "ORDER123",
  "mpsbill_name": "John Doe",
  "mpsbill_email": "john@example.com",
  "mpsbill_mobile": "01234567890",
  "mpsbill_desc": "Payment",
  "mpscurrency": "MYR",
  "mpslangcode": "en",
  "vcode": "generated_hash_code",
  "sandbox": true
}
```

### 3. Process Payment

```bash
POST /api/v1/payments/process
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "order_id": "ORDER123",
  "amount": 5000,
  "currency": "MYR",
  "payment_method": "fpx",
  "provider": "fiuu",
  "channel": "fpx",
  "token": "payment_token",
  "bill_name": "John Doe",
  "bill_email": "john@example.com",
  "bill_mobile": "01234567890",
  "bill_desc": "Payment for auction item"
}
```

### 4. Webhook Handling

Fiuu sends payment status updates to your webhook endpoint:

```bash
POST /api/v1/webhooks/fiuu
Content-Type: application/json

{
  "tranID": "123456789",
  "order_id": "ORDER123",
  "amount": "50.00",
  "currency": "MYR",
  "payment_status": "00",
  "payment_channel": "fpx",
  "channel_code": "fpx",
  "payment_ref_id": "REF123",
  "pay_date": "2025-01-01",
  "pay_time": "12:00:00",
  "error_code": "",
  "error_desc": "",
  "signature": "webhook_signature"
}
```

## Frontend Integration

### Basic Implementation

1. **Include Fiuu JavaScript**
```html
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js"></script>
<script src="https://pay.fiuu.com/RMS/API/seamless/3.28/js/MOLPay_seamless.deco.js"></script>
```

2. **Get Configuration from Backend**
```javascript
async function getPaymentConfig(orderId, amount, channel) {
  const response = await fetch(`/api/v1/payments/seamless/config?` + new URLSearchParams({
    order_id: orderId,
    amount: amount.toString(),
    bill_name: 'John Doe',
    bill_email: 'john@example.com',
    bill_mobile: '01234567890',
    bill_desc: 'Payment for auction',
    channel: channel
  }));
  
  return await response.json();
}
```

3. **Initialize Payment**
```javascript
$(document).ready(function() {
  const config = await getPaymentConfig('ORDER123', 5000, 'fpx');
  
  $('#pay-button').MOLPaySeamless({
    mpsmerchantid: config.mpsmerchantid,
    mpschannel: config.mpschannel,
    mpsamount: config.mpsamount,
    mpsorderid: config.mpsorderid,
    mpsbill_name: config.mpsbill_name,
    mpsbill_email: config.mpsbill_email,
    mpsbill_mobile: config.mpsbill_mobile,
    mpsbill_desc: config.mpsbill_desc,
    mpscurrency: config.mpscurrency,
    mpslangcode: config.mpslangcode,
    vcode: config.vcode
  });
});
```

## Testing

### Sandbox Environment

Use sandbox mode for testing:

```bash
FIUU_SANDBOX=true
FIUU_MERCHANT_ID=test_merchant_id
FIUU_VERIFY_KEY=test_verify_key
```

### Test Payments

1. Use the provided HTML example: `examples/fiuu-frontend-integration.html`
2. Select a payment method and click "Get Config"
3. Click "Pay Now" to test the payment flow
4. Check webhook delivery and payment status updates

## Security Considerations

### VCode Generation
- Never generate vcode in frontend JavaScript
- Always generate vcode on backend using the verify key
- The verify key must never be exposed to clients

### Webhook Validation
- Validate webhook signatures using Fiuu's method
- Ensure webhook endpoints are secure
- Use HTTPS for all webhook URLs

### Data Protection
- Encrypt sensitive payment data
- Comply with PCI DSS requirements
- Log payment activities for audit trails

## Error Handling

### Common Error Codes

| Error Code | Description | Action |
|------------|-------------|--------|
| `00` | Success | Payment completed |
| `11` | Pending | Payment processing |
| `22` | Invalid parameters | Check request data |
| `44` | Payment failed | Retry with different method |
| `99` | System error | Contact support |

### Error Response Format

```json
{
  "error": "payment_failed",
  "message": "Payment processing failed: Invalid channel",
  "details": {
    "code": "22",
    "description": "Invalid payment channel"
  }
}
```

## Monitoring and Logging

### Key Metrics
- Payment success rate by channel
- Average processing time
- Webhook delivery success rate
- Error rates by payment method

### Logging
The integration provides comprehensive logging:

```go
logger.Info("Fiuu payment created successfully", 
    zap.String("transaction_id", resp.TransactionID),
    zap.String("order_id", resp.OrderID),
    zap.String("channel", req.Channel))
```

## Support

### Fiuu Support
- **Technical Support**: support@fiuu.com
- **Sales/Reseller**: sales@fiuu.com
- **Channel/Partner**: channel@fiuu.com
- **R&D Suggestions**: technical@fiuu.com

### Resources
- **Documentation**: https://github.com/FiuuPayment/Integration-Fiuu_JavaScript_Seamless_Integration
- **Website**: https://fiuu.com
- **GitHub**: https://github.com/FiuuPayment

## Migration from Mock

To migrate from mock payment processing:

1. **Set up Fiuu account** and get credentials
2. **Configure environment variables** with Fiuu settings
3. **Update frontend** to use seamless integration
4. **Test in sandbox** environment
5. **Deploy to production** with live credentials

The service automatically detects Fiuu credentials and switches from mock to real processing.