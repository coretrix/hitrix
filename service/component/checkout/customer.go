package checkout

type Instrument struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Fingerprint   string `json:"fingerprint"`
	ExpiryMonth   int    `json:"expiry_month"`
	ExpiryYear    int    `json:"expiry_year"`
	Name          string `json:"name"`
	Scheme        string `json:"scheme"`
	Last4         string `json:"last4"`
	Bin           string `json:"bin"`
	CardType      string `json:"card_type"`
	CardCategory  string `json:"card_category"`
	Issuer        string `json:"issuer"`
	IssuerCountry string `json:"issuer_country"`
	ProductID     string `json:"product_id"`
	ProductType   string `json:"product_type"`
	AccountHolder struct {
		BillingAddress struct {
			AddressLine1 string `json:"address_line1"`
			AddressLine2 string `json:"address_line2"`
			City         string `json:"city"`
			State        string `json:"state"`
			Zip          string `json:"zip"`
			Country      string `json:"country"`
		} `json:"billing_address"`
		Phone struct {
			CountryCode string `json:"country_code"`
			Number      string `json:"number"`
		} `json:"phone"`
	} `json:"account_holder"`
}

type CustomerPhone struct {
	Number      string `json:"number"`
	CountryCode string `json:"country_code"`
}

type CustomerResponse struct {
	ID       string        `json:"id"`
	Email    string        `json:"email"`
	Default  string        `json:"default"`
	Name     string        `json:"name"`
	Phone    CustomerPhone `json:"phone,omitempty"`
	Metadata struct {
		CouponCode string `json:"coupon_code"`
		PartnerID  int    `json:"partner_id"`
	} `json:"metadata"`
	Instruments []Instrument `json:"instruments"`
}

type SaveCustomerRequest struct {
	Email    string            `json:"email"`
	Name     string            `json:"name"`
	Phone    *CustomerPhone    `json:"phone,omitempty"`
	Metadata map[string]string `json:"metadata"`
}
