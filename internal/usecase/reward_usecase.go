package usecase

import (
	"backend/configs"
	"backend/internal/entity"
	"backend/internal/repository"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type rewardUseCase struct {
	clientID     string
	clientSecret string
	grantType    string
	audience     string
	userRepo     repository.UserRepository
	rewardRepo   repository.RewardRepository
}

func NewRewardUseCase(cfg configs.Reward, userRepo repository.UserRepository, rewardRepo repository.RewardRepository) RewardUseCase {
	return &rewardUseCase{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		grantType:    cfg.GrantType,
		audience:     cfg.Audience,
		userRepo:     userRepo,
		rewardRepo:   rewardRepo,
	}
}

// GetBallance implements TopUpUseCase.
func (t *rewardUseCase) GetBalance(ctx context.Context) (res float64, err error) {
	token, err := t.GetToken(ctx)
	if err != nil {
		return 0, err
	}

	url := "https://topups-sandbox.reloadly.com/accounts/balance"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/com.reloadly.topups-v1+json")
	req.Header.Add("Authorization", "Bearer "+token.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var balance entity.Balance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return 0, err
	}

	return balance.Balance, nil
}

// MakeTopUp implements TopUpUseCase.
func (t *rewardUseCase) GetReward(ctx context.Context, userId uint, operatorID string, amount string, number string) (res string, err error) {
	url := "https://topups-sandbox.reloadly.com/topups"

	coinUser, err := t.userRepo.GetCoinByUserID(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("failed to get coin: %w", err)
	}

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return "", fmt.Errorf("failed to convert amount to int: %w", err)
	}

	rewardPrice, err := t.rewardRepo.GetPriceByAmount(ctx, amountInt)
	if err != nil {
		return "", fmt.Errorf("failed to get reward price: %w", err)
	}

	if coinUser < rewardPrice {
		return "", err
	}

	token, err := t.GetToken(ctx)
	if err != nil {
		return "", err
	}

	phone := *&entity.Phone{
		CountryCode: "ID",
		Number:      number,
	}

	data := *&entity.RewardRequest{
		OperatorID:     operatorID,
		Amount:         amount,
		UseLocalAmount: "true",
		RecipientEmail: "example@example.com",
		RecipientPhone: phone,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	payload := bytes.NewReader(jsonBytes)

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Accept", "application/com.reloadly.topups-v1+json")
	req.Header.Add("Authorization", "Bearer "+token.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	coinUser -= rewardPrice
	
	var topUpStatus entity.RewardStatus
	if err := json.NewDecoder(resp.Body).Decode(&topUpStatus); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if topUpStatus.Status != "SUCCESSFUL" {
		return "", fmt.Errorf("failed to top up: %w", err)
	}

	return "success", nil
}

// GetToken implements TopUpUseCase.
func (t *rewardUseCase) GetToken(ctx context.Context) (res *entity.TokenReward, err error) {
	url := "https://auth.reloadly.com/oauth/token"

	data := *&entity.TokenRewardRequest{
		ClientID:     t.clientID,
		ClientSecret: t.clientSecret,
		GrantType:    t.grantType,
		Audience:     t.audience,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	payload := bytes.NewReader(jsonBytes)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var token entity.TokenReward
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &token, nil

}
