package entity

type TopUp struct {
	OperatorID       string `json:"operator_id"`
	Amount           string `json:"amount"`
	UseLocalAmount   string `json:"use_local_amount"`
	CustomIdentifier string `json:"custom_identifier"`
	RecipientEmail   string `json:"recipient_email"`
	RecipientPhone   Phone
}

type Phone struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
}

type TopUpResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transactionId"`
}

