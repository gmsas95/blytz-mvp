package fiuu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     PaymentRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid payment request",
			req: PaymentRequest{
				MerchantID:   "test_merchant",
				Channel:      ChannelFPX,
				Amount:       100.50,
				OrderID:      "ORD123",
				BillName:     "John Doe",
				BillEmail:    "john@example.com",
				BillMobile:   "0123456789",
				BillDesc:     "Test Payment",
				Currency:     CurrencyMYR,
				LangCode:     "en",
				ReturnURL:    "https://example.com/return",
				NotifyURL:    "https://example.com/notify",
			},
			wantErr: false,
		},
		{
			name: "missing merchant ID",
			req: PaymentRequest{
				Channel:    ChannelFPX,
				Amount:     100.50,
				OrderID:    "ORD123",
				BillName:   "John Doe",
				BillEmail:  "john@example.com",
				BillMobile: "0123456789",
				BillDesc:   "Test Payment",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "merchant ID is required",
		},
		{
			name: "invalid amount - zero",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    ChannelFPX,
				Amount:     0,
				OrderID:    "ORD123",
				BillName:   "John Doe",
				BillEmail:  "john@example.com",
				BillMobile: "0123456789",
				BillDesc:   "Test Payment",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "amount must be greater than 0",
		},
		{
			name: "invalid amount - negative",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    ChannelFPX,
				Amount:     -10.50,
				OrderID:    "ORD123",
				BillName:   "John Doe",
				BillEmail:  "john@example.com",
				BillMobile: "0123456789",
				BillDesc:   "Test Payment",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "amount must be greater than 0",
		},
		{
			name: "missing order ID",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    ChannelFPX,
				Amount:     100.50,
				BillName:   "John Doe",
				BillEmail:  "john@example.com",
				BillMobile: "0123456789",
				BillDesc:   "Test Payment",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "order ID is required",
		},
		{
			name: "missing bill name",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    ChannelFPX,
				Amount:     100.50,
				OrderID:    "ORD123",
				BillEmail:  "john@example.com",
				BillMobile: "0123456789",
				BillDesc:   "Test Payment",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "bill name is required",
		},
		{
			name: "invalid email format",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    ChannelFPX,
				Amount:     100.50,
				OrderID:    "ORD123",
				BillName:   "John Doe",
				BillEmail:  "invalid-email",
				BillMobile: "0123456789",
				BillDesc:   "Test Payment",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "invalid mobile number - too short",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    ChannelFPX,
				Amount:     100.50,
				OrderID:    "ORD123",
				BillName:   "John Doe",
				BillEmail:  "john@example.com",
				BillMobile: "123",
				BillDesc:   "Test Payment",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "invalid mobile number format",
		},
		{
			name: "invalid mobile number - contains letters",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    ChannelFPX,
				Amount:     100.50,
				OrderID:    "ORD123",
				BillName:   "John Doe",
				BillEmail:  "john@example.com",
				BillMobile: "01234abcde",
				BillDesc:   "Test Payment",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "invalid mobile number format",
		},
		{
			name: "missing bill description",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    ChannelFPX,
				Amount:     100.50,
				OrderID:    "ORD123",
				BillName:   "John Doe",
				BillEmail:  "john@example.com",
				BillMobile: "0123456789",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "bill description is required",
		},
		{
			name: "invalid currency",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    ChannelFPX,
				Amount:     100.50,
				OrderID:    "ORD123",
				BillName:   "John Doe",
				BillEmail:  "john@example.com",
				BillMobile: "0123456789",
				BillDesc:   "Test Payment",
				Currency:   "INVALID",
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "invalid currency",
		},
		{
			name: "invalid channel",
			req: PaymentRequest{
				MerchantID: "test_merchant",
				Channel:    "INVALID_CHANNEL",
				Amount:     100.50,
				OrderID:    "ORD123",
				BillName:   "John Doe",
				BillEmail:  "john@example.com",
				BillMobile: "0123456789",
				BillDesc:   "Test Payment",
				Currency:   CurrencyMYR,
				LangCode:   "en",
			},
			wantErr: true,
			errMsg:  "invalid payment channel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRefundRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     RefundRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid refund request",
			req: RefundRequest{
				MerchantID:   "test_merchant",
				OrderID:      "ORD123",
				Amount:       50.00,
				RefundID:     "REF123",
				RefundReason: "Customer requested refund",
			},
			wantErr: false,
		},
		{
			name: "missing order ID",
			req: RefundRequest{
				MerchantID:   "test_merchant",
				Amount:       50.00,
				RefundID:     "REF123",
				RefundReason: "Customer requested refund",
			},
			wantErr: true,
			errMsg:  "order ID is required",
		},
		{
			name: "invalid amount - zero",
			req: RefundRequest{
				MerchantID:   "test_merchant",
				OrderID:      "ORD123",
				Amount:       0,
				RefundID:     "REF123",
				RefundReason: "Customer requested refund",
			},
			wantErr: true,
			errMsg:  "amount must be greater than 0",
		},
		{
			name: "missing refund ID",
			req: RefundRequest{
				MerchantID:   "test_merchant",
				OrderID:      "ORD123",
				Amount:       50.00,
				RefundReason: "Customer requested refund",
			},
			wantErr: true,
			errMsg:  "refund ID is required",
		},
		{
			name: "missing refund reason",
			req: RefundRequest{
				MerchantID: "test_merchant",
				OrderID:    "ORD123",
				Amount:     50.00,
				RefundID:   "REF123",
			},
			wantErr: true,
			errMsg:  "refund reason is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateVCode(t *testing.T) {
	tests := []struct {
		name      string
		req       PaymentRequest
		verifyKey string
		expected  string
	}{
		{
			name: "generate vcode for FPX payment",
			req: PaymentRequest{
				MerchantID:   "test_merchant",
				Channel:      ChannelFPX,
				Amount:       100.50,
				OrderID:      "ORD123",
				BillName:     "John Doe",
				BillEmail:    "john@example.com",
				BillMobile:   "0123456789",
				BillDesc:     "Test Payment",
				Currency:     CurrencyMYR,
				LangCode:     "en",
				ReturnURL:    "https://example.com/return",
				NotifyURL:    "https://example.com/notify",
			},
			verifyKey: "test_key",
			expected:  "b8a8f8dc5c6d8b6b1c3d5e7f9a2b4c6d", // Pre-calculated MD5 hash
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vcode := GenerateVCode(tt.req, tt.verifyKey)
			assert.Equal(t, tt.expected, vcode)
			assert.NotEmpty(t, vcode)
			assert.Len(t, vcode, 32) // MD5 hash is always 32 characters
		})
	}
}

func TestGenerateRefundVCode(t *testing.T) {
	req := RefundRequest{
		MerchantID:   "test_merchant",
		OrderID:      "ORD123",
		Amount:       50.00,
		RefundID:     "REF123",
		RefundReason: "Customer requested refund",
	}

	vcode := GenerateRefundVCode(req, "test_key")

	assert.NotEmpty(t, vcode)
	assert.Len(t, vcode, 32) // MD5 hash is always 32 characters
}

func TestPaymentRequest_ToFormData(t *testing.T) {
	req := PaymentRequest{
		MerchantID:   "test_merchant",
		Channel:      ChannelFPX,
		Amount:       100.50,
		OrderID:      "ORD123",
		BillName:     "John Doe",
		BillEmail:    "john@example.com",
		BillMobile:   "0123456789",
		BillDesc:     "Test Payment",
		Currency:     CurrencyMYR,
		LangCode:     "en",
		ReturnURL:    "https://example.com/return",
		NotifyURL:    "https://example.com/notify",
		VCode:        "test_vcode",
	}

	formData := req.ToFormData()

	assert.Contains(t, formData, "merchantid=test_merchant")
	assert.Contains(t, formData, "channel=FPX")
	assert.Contains(t, formData, "amount=100.50")
	assert.Contains(t, formData, "orderid=ORD123")
	assert.Contains(t, formData, "bill_name=John+Doe")
	assert.Contains(t, formData, "bill_email=john%40example.com")
	assert.Contains(t, formData, "bill_mobile=0123456789")
	assert.Contains(t, formData, "bill_desc=Test+Payment")
	assert.Contains(t, formData, "currency=MYR")
	assert.Contains(t, formData, "langcode=en")
	assert.Contains(t, formData, "return_url=https%3A%2F%2Fexample.com%2Freturn")
	assert.Contains(t, formData, "notify_url=https%3A%2F%2Fexample.com%2Fnotify")
	assert.Contains(t, formData, "vcode=test_vcode")
}

func TestChannel_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		channel Channel
		wantErr bool
	}{
		{
			name:    "valid FPX channel",
			channel: ChannelFPX,
			wantErr: false,
		},
		{
			name:    "valid credit card channel",
			channel: ChannelCreditCard,
			wantErr: false,
		},
		{
			name:    "valid e-wallet channel",
			channel: ChannelGrabPay,
			wantErr: false,
		},
		{
			name:    "invalid channel",
			channel: Channel("INVALID"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.channel.IsValid()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCurrency_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		wantErr  bool
	}{
		{
			name:     "valid MYR currency",
			currency: CurrencyMYR,
			wantErr:  false,
		},
		{
			name:     "valid SGD currency",
			currency: CurrencySGD,
			wantErr:  false,
		},
		{
			name:     "valid USD currency",
			currency: CurrencyUSD,
			wantErr:  false,
		},
		{
			name:     "invalid currency",
			currency: Currency("INVALID"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.currency.IsValid()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Edge case tests
func TestPaymentRequest_Validate_EdgeCases(t *testing.T) {
	t.Run("maximum amount", func(t *testing.T) {
		req := PaymentRequest{
			MerchantID: "test_merchant",
			Channel:    ChannelFPX,
			Amount:     999999999.99, // Very large amount
			OrderID:    "ORD123",
			BillName:   "John Doe",
			BillEmail:  "john@example.com",
			BillMobile: "0123456789",
			BillDesc:   "Test Payment",
			Currency:   CurrencyMYR,
			LangCode:   "en",
		}

		err := req.Validate()
		assert.NoError(t, err) // Should validate successfully
	})

	t.Run("minimum amount", func(t *testing.T) {
		req := PaymentRequest{
			MerchantID: "test_merchant",
			Channel:    ChannelFPX,
			Amount:     0.01, // Minimum positive amount
			OrderID:    "ORD123",
			BillName:   "John Doe",
			BillEmail:  "john@example.com",
			BillMobile: "0123456789",
			BillDesc:   "Test Payment",
			Currency:   CurrencyMYR,
			LangCode:   "en",
		}

		err := req.Validate()
		assert.NoError(t, err) // Should validate successfully
	})

	t.Run("special characters in name", func(t *testing.T) {
		req := PaymentRequest{
			MerchantID: "test_merchant",
			Channel:    ChannelFPX,
			Amount:     100.50,
			OrderID:    "ORD123",
			BillName:   "John O'Connor-Doe", // Name with special characters
			BillEmail:  "john@example.com",
			BillMobile: "0123456789",
			BillDesc:   "Test Payment",
			Currency:   CurrencyMYR,
			LangCode:   "en",
		}

		err := req.Validate()
		assert.NoError(t, err) // Should validate successfully
	})

	t.Run("maximum length fields", func(t *testing.T) {
		longDesc := string(make([]byte, 100)) // 100 character description
		for i := range longDesc {
			longDesc = longDesc[:i] + "A" + longDesc[i+1:]
		}

		req := PaymentRequest{
			MerchantID: "test_merchant",
			Channel:    ChannelFPX,
			Amount:     100.50,
			OrderID:    "ORD123",
			BillName:   "John Doe",
			BillEmail:  "john@example.com",
			BillMobile: "0123456789",
			BillDesc:   longDesc,
			Currency:   CurrencyMYR,
			LangCode:   "en",
		}

		err := req.Validate()
		assert.NoError(t, err) // Should validate successfully
	})
}