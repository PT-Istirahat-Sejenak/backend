package entity

type RewardRequest struct {
	OperatorID     string `json:"operatorId"`
	Amount         string `json:"amount"`
	UseLocalAmount string `json:"useLocalAmount"`
	RecipientEmail string `json:"recipientEmail"`
	RecipientPhone Phone
}

type Phone struct {
	CountryCode string `json:"countryCode"`
	Number      string `json:"number"`
}

type RewardResponse struct {
	Status string `json:"status"`
}

type Balance struct {
	Balance float64 `json:"balance"`
}

type TokenReward struct {
	Token     string `json:"access_token"`
	Scope     string `json:"scope"`
	Exp       int64  `json:"expires_in"`
	TokenType string `json:"token_type"`
}

type TokenRewardRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Audience     string `json:"audience"`
}

type RewardStatus struct {
	Status string `json:"status"`
}

type PriceReward struct {
	Amount int
	Price  int
}
