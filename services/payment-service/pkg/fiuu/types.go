package fiuu

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Channel represents Fiuu payment channel
type Channel string

const (
	ChannelFPX       Channel = "fpx"
	ChannelCredit    Channel = "credit"
	ChannelPayPal    Channel = "paypal"
	ChannelGrabPay   Channel = "GrabPay"
	ChannelTNG       Channel = "TNG-EWALLET"
	ChannelShopeePay Channel = "ShopeePay"
	ChannelBoost     Channel = "BOOST"
	ChannelApplePay  Channel = "apple_pay"
	ChannelGooglePay Channel = "google_pay"
	ChannelAlipay    Channel = "alipay"
	ChannelFPXB2B    Channel = "FPX_B2B"
	ChannelDuitNowQR Channel = "RPP_DuitNowQR"
	ChannelAtome     Channel = "Atome"
	ChannelRely      Channel = "Rely-PW"
)

// Currency represents supported currencies
type Currency string

const (
	CurrencyMYR Currency = "MYR"
	CurrencySGD Currency = "SGD"
	CurrencyUSD Currency = "USD"
	CurrencyTHB Currency = "THB"
	CurrencyPHP Currency = "PHP"
	CurrencyIDR Currency = "IDR"
)

// PaymentRequest represents Fiuu payment request
type PaymentRequest struct {
	MerchantID            string   `json:"mpsmerchantid"`
	Channel               Channel  `json:"mpschannel"`
	Amount                float64  `json:"mpsamount"`
	OrderID               string   `json:"mpsorderid"`
	BillName              string   `json:"mpsbill_name"`
	BillEmail             string   `json:"mpsbill_email"`
	BillMobile            string   `json:"mpsbill_mobile"`
	BillDesc              string   `json:"mpsbill_desc"`
	Country               string   `json:"mpscountry,omitempty"`
	BillAddress           string   `json:"mpsbill_add,omitempty"`
	BillZip               string   `json:"mpsbill_zip,omitempty"`
	Currency              Currency `json:"mpscurrency"`
	LangCode              string   `json:"mpslangcode,omitempty"`
	ReturnURL             string   `json:"mpsreturnurl,omitempty"`
	NotifyURL             string   `json:"mpsnotifyurl,omitempty"`
	CallbackURL           string   `json:"mpscallbackurl,omitempty"`
	Timer                 int      `json:"mpstimer,omitempty"`
	TimerBox              string   `json:"mpstimerbox,omitempty"`
	CancelURL             string   `json:"mpscancelurl,omitempty"`
	TokenStatus           int      `json:"mpstokenstatus,omitempty"`
	TCCType               string   `json:"mpstcctype,omitempty"`
	GuestCheckout         int      `json:"mpsguestcheckout,omitempty"`
	Non3DS                int      `json:"non_3ds,omitempty"`
	Extra                 string   `json:"mpsextra,omitempty"`
	SplitInfo             string   `json:"mpsplitinfo,omitempty"`
	InstallMonth          int      `json:"mpsinstallmonth,omitempty"`
	IsDDA                 int      `json:"mpsis_DDA,omitempty"`
	BuyerIDNumber         string   `json:"mpsbuyerid_number,omitempty"`
	BuyerIDType           int      `json:"mpsbuyerid_type,omitempty"`
	CryptoCurrency        string   `json:"mpscyrptocurrency,omitempty"`
	PaymentExpirationTime string   `json:"mpspayment_expiration_time,omitempty"`
	VCode                 string   `json:"vcode"`
}

// PaymentResponse represents Fiuu payment response
type PaymentResponse struct {
	TransactionID    string `json:"tranID"`
	OrderID          string `json:"order_id"`
	Amount           string `json:"amount"`
	Currency         string `json:"currency"`
	PaymentStatus    string `json:"payment_status"`
	PaymentChannel   string `json:"payment_channel"`
	ChannelCode      string `json:"channel_code"`
	PaymentRefID     string `json:"payment_ref_id"`
	PayDate          string `json:"pay_date"`
	PayTime          string `json:"pay_time"`
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_desc"`
	Signature        string `json:"signature"`
}

// RefundRequest represents Fiuu refund request
type RefundRequest struct {
	MerchantID   string  `json:"merchantid"`
	OrderID      string  `json:"orderid"`
	Amount       float64 `json:"amount"`
	RefundID     string  `json:"refundid"`
	RefundReason string  `json:"refundreason"`
	VCode        string  `json:"vcode"`
}

// RefundResponse represents Fiuu refund response
type RefundResponse struct {
	RefundID         string `json:"refundid"`
	OrderID          string `json:"orderid"`
	Amount           string `json:"amount"`
	RefundStatus     string `json:"refundstatus"`
	RefundDate       string `json:"refunddate"`
	RefundTime       string `json:"refundtime"`
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_desc"`
	Signature        string `json:"signature"`
}

// ChannelInfo provides information about payment channels
type ChannelInfo struct {
	Code        Channel  `json:"code"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Currency    Currency `json:"currency"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
}

// GetAvailableChannels returns all available Fiuu payment channels
func GetAvailableChannels() []ChannelInfo {
	return []ChannelInfo{
		{
			Code:        ChannelFPX,
			Name:        "FPX Online Banking",
			Type:        "online_banking",
			Currency:    CurrencyMYR,
			Description: "Pay with Malaysian online banking",
			Enabled:     true,
		},
		{
			Code:        ChannelCredit,
			Name:        "Credit/Debit Card",
			Type:        "card",
			Currency:    CurrencyMYR,
			Description: "Pay with credit or debit card",
			Enabled:     true,
		},
		{
			Code:        ChannelGrabPay,
			Name:        "GrabPay",
			Type:        "ewallet",
			Currency:    CurrencyMYR,
			Description: "Pay with GrabPay e-wallet",
			Enabled:     true,
		},
		{
			Code:        ChannelTNG,
			Name:        "Touch 'n Go",
			Type:        "ewallet",
			Currency:    CurrencyMYR,
			Description: "Pay with Touch 'n Go e-wallet",
			Enabled:     true,
		},
		{
			Code:        ChannelShopeePay,
			Name:        "ShopeePay",
			Type:        "ewallet",
			Currency:    CurrencyMYR,
			Description: "Pay with ShopeePay e-wallet",
			Enabled:     true,
		},
		{
			Code:        ChannelBoost,
			Name:        "Boost",
			Type:        "ewallet",
			Currency:    CurrencyMYR,
			Description: "Pay with Boost e-wallet",
			Enabled:     true,
		},
		{
			Code:        ChannelFPXB2B,
			Name:        "FPX B2B",
			Type:        "online_banking",
			Currency:    CurrencyMYR,
			Description: "Business banking transfer",
			Enabled:     true,
		},
		{
			Code:        ChannelDuitNowQR,
			Name:        "DuitNow QR",
			Type:        "qr",
			Currency:    CurrencyMYR,
			Description: "Pay with DuitNow QR code",
			Enabled:     true,
		},
		{
			Code:        ChannelAtome,
			Name:        "Atome",
			Type:        "bnpl",
			Currency:    CurrencyMYR,
			Description: "Buy Now Pay Later with Atome",
			Enabled:     true,
		},
		{
			Code:        ChannelRely,
			Name:        "Rely",
			Type:        "bnpl",
			Currency:    CurrencyMYR,
			Description: "Buy Now Pay Later with Rely",
			Enabled:     true,
		},
		{
			Code:        ChannelPayPal,
			Name:        "PayPal",
			Type:        "wallet",
			Currency:    CurrencyUSD,
			Description: "Pay with PayPal",
			Enabled:     true,
		},
		{
			Code:        ChannelApplePay,
			Name:        "Apple Pay",
			Type:        "wallet",
			Currency:    CurrencyMYR,
			Description: "Pay with Apple Pay",
			Enabled:     false, // Requires additional setup
		},
		{
			Code:        ChannelGooglePay,
			Name:        "Google Pay",
			Type:        "wallet",
			Currency:    CurrencyMYR,
			Description: "Pay with Google Pay",
			Enabled:     false, // Requires additional setup
		},
	}
}

// GenerateVCode generates the verification code (vcode) for Fiuu API
func GenerateVCode(req PaymentRequest, verifyKey string) string {
	// Create a map of all parameters for vcode generation
	params := map[string]string{
		"mpsmerchantid":  req.MerchantID,
		"mpschannel":     string(req.Channel),
		"mpsamount":      fmt.Sprintf("%.2f", req.Amount),
		"mpsorderid":     req.OrderID,
		"mpsbill_name":   req.BillName,
		"mpsbill_email":  req.BillEmail,
		"mpsbill_mobile": req.BillMobile,
		"mpsbill_desc":   req.BillDesc,
		"mpscurrency":    string(req.Currency),
	}

	// Add optional parameters if they exist
	if req.Country != "" {
		params["mpscountry"] = req.Country
	}
	if req.BillAddress != "" {
		params["mpsbill_add"] = req.BillAddress
	}
	if req.BillZip != "" {
		params["mpsbill_zip"] = req.BillZip
	}
	if req.LangCode != "" {
		params["mpslangcode"] = req.LangCode
	}
	if req.ReturnURL != "" {
		params["mpsreturnurl"] = req.ReturnURL
	}
	if req.NotifyURL != "" {
		params["mpsnotifyurl"] = req.NotifyURL
	}
	if req.CallbackURL != "" {
		params["mpscallbackurl"] = req.CallbackURL
	}
	if req.Timer > 0 {
		params["mpstimer"] = fmt.Sprintf("%d", req.Timer)
	}
	if req.TimerBox != "" {
		params["mpstimerbox"] = req.TimerBox
	}
	if req.CancelURL != "" {
		params["mpscancelurl"] = req.CancelURL
	}
	if req.TokenStatus > 0 {
		params["mpstokenstatus"] = fmt.Sprintf("%d", req.TokenStatus)
	}
	if req.TCCType != "" {
		params["mpstcctype"] = req.TCCType
	}
	if req.GuestCheckout > 0 {
		params["mpsguestcheckout"] = fmt.Sprintf("%d", req.GuestCheckout)
	}
	if req.Non3DS > 0 {
		params["non_3ds"] = fmt.Sprintf("%d", req.Non3DS)
	}
	if req.Extra != "" {
		params["mpsextra"] = req.Extra
	}
	if req.SplitInfo != "" {
		params["mpsplitinfo"] = req.SplitInfo
	}
	if req.InstallMonth > 0 {
		params["mpsinstallmonth"] = fmt.Sprintf("%d", req.InstallMonth)
	}
	if req.IsDDA > 0 {
		params["mpsis_DDA"] = fmt.Sprintf("%d", req.IsDDA)
	}
	if req.BuyerIDNumber != "" {
		params["mpsbuyerid_number"] = req.BuyerIDNumber
	}
	if req.BuyerIDType > 0 {
		params["mpsbuyerid_type"] = fmt.Sprintf("%d", req.BuyerIDType)
	}
	if req.CryptoCurrency != "" {
		params["mpscyrptocurrency"] = req.CryptoCurrency
	}
	if req.PaymentExpirationTime != "" {
		params["mpspayment_expiration_time"] = req.PaymentExpirationTime
	}

	// Sort keys for consistent hash generation
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the string to hash
	var parts []string
	for _, key := range keys {
		if params[key] != "" {
			parts = append(parts, key+params[key])
		}
	}

	// Append verify key
	parts = append(parts, verifyKey)

	// Join and create MD5 hash
	dataToHash := strings.Join(parts, "")
	hash := md5.Sum([]byte(dataToHash))

	return fmt.Sprintf("%x", hash)
}

// GenerateRefundVCode generates the verification code for refund requests
func GenerateRefundVCode(req RefundRequest, verifyKey string) string {
	params := map[string]string{
		"merchantid":   req.MerchantID,
		"orderid":      req.OrderID,
		"amount":       fmt.Sprintf("%.2f", req.Amount),
		"refundid":     req.RefundID,
		"refundreason": req.RefundReason,
	}

	// Sort keys for consistent hash generation
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the string to hash
	var parts []string
	for _, key := range keys {
		parts = append(parts, key+params[key])
	}

	// Append verify key
	parts = append(parts, verifyKey)

	// Join and create MD5 hash
	dataToHash := strings.Join(parts, "")
	hash := md5.Sum([]byte(dataToHash))

	return fmt.Sprintf("%x", hash)
}

// ToFormData converts PaymentRequest to URL-encoded form data
func (req PaymentRequest) ToFormData() string {
	values := url.Values{}

	values.Set("mpsmerchantid", req.MerchantID)
	values.Set("mpschannel", string(req.Channel))
	values.Set("mpsamount", fmt.Sprintf("%.2f", req.Amount))
	values.Set("mpsorderid", req.OrderID)
	values.Set("mpsbill_name", req.BillName)
	values.Set("mpsbill_email", req.BillEmail)
	values.Set("mpsbill_mobile", req.BillMobile)
	values.Set("mpsbill_desc", req.BillDesc)
	values.Set("mpscurrency", string(req.Currency))
	values.Set("vcode", req.VCode)

	// Add optional parameters
	if req.Country != "" {
		values.Set("mpscountry", req.Country)
	}
	if req.BillAddress != "" {
		values.Set("mpsbill_add", req.BillAddress)
	}
	if req.BillZip != "" {
		values.Set("mpsbill_zip", req.BillZip)
	}
	if req.LangCode != "" {
		values.Set("mpslangcode", req.LangCode)
	}
	if req.ReturnURL != "" {
		values.Set("mpsreturnurl", req.ReturnURL)
	}
	if req.NotifyURL != "" {
		values.Set("mpsnotifyurl", req.NotifyURL)
	}
	if req.CallbackURL != "" {
		values.Set("mpscallbackurl", req.CallbackURL)
	}
	if req.Timer > 0 {
		values.Set("mpstimer", fmt.Sprintf("%d", req.Timer))
	}
	if req.TimerBox != "" {
		values.Set("mpstimerbox", req.TimerBox)
	}
	if req.CancelURL != "" {
		values.Set("mpscancelurl", req.CancelURL)
	}
	if req.TokenStatus > 0 {
		values.Set("mpstokenstatus", fmt.Sprintf("%d", req.TokenStatus))
	}
	if req.TCCType != "" {
		values.Set("mpstcctype", req.TCCType)
	}
	if req.GuestCheckout > 0 {
		values.Set("mpsguestcheckout", fmt.Sprintf("%d", req.GuestCheckout))
	}
	if req.Non3DS > 0 {
		values.Set("non_3ds", fmt.Sprintf("%d", req.Non3DS))
	}
	if req.Extra != "" {
		values.Set("mpsextra", req.Extra)
	}
	if req.SplitInfo != "" {
		values.Set("mpsplitinfo", req.SplitInfo)
	}
	if req.InstallMonth > 0 {
		values.Set("mpsinstallmonth", fmt.Sprintf("%d", req.InstallMonth))
	}
	if req.IsDDA > 0 {
		values.Set("mpsis_DDA", fmt.Sprintf("%d", req.IsDDA))
	}
	if req.BuyerIDNumber != "" {
		values.Set("mpsbuyerid_number", req.BuyerIDNumber)
	}
	if req.BuyerIDType > 0 {
		values.Set("mpsbuyerid_type", fmt.Sprintf("%d", req.BuyerIDType))
	}
	if req.CryptoCurrency != "" {
		values.Set("mpscyrptocurrency", req.CryptoCurrency)
	}
	if req.PaymentExpirationTime != "" {
		values.Set("mpspayment_expiration_time", req.PaymentExpirationTime)
	}

	return values.Encode()
}

// Validate checks if the payment request is valid
func (req PaymentRequest) Validate() error {
	if req.MerchantID == "" {
		return fmt.Errorf("merchant ID is required")
	}
	if req.Channel == "" {
		return fmt.Errorf("payment channel is required")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if req.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}
	if req.BillName == "" {
		return fmt.Errorf("bill name is required")
	}
	if req.BillEmail == "" {
		return fmt.Errorf("bill email is required")
	}
	if req.BillMobile == "" {
		return fmt.Errorf("bill mobile is required")
	}
	if req.BillDesc == "" {
		return fmt.Errorf("bill description is required")
	}
	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	return nil
}

// Validate checks if the refund request is valid
func (req RefundRequest) Validate() error {
	if req.MerchantID == "" {
		return fmt.Errorf("merchant ID is required")
	}
	if req.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if req.RefundID == "" {
		return fmt.Errorf("refund ID is required")
	}
	if req.RefundReason == "" {
		return fmt.Errorf("refund reason is required")
	}
	return nil
}
