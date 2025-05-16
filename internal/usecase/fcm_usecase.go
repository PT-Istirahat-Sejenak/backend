package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"golang.org/x/oauth2/google"
)

type fcmUseCase struct {
	userRepo repository.UserRepository
}

func NewFcmUseCase(userRepo repository.UserRepository) FCMUseCase {
	return &fcmUseCase{
		userRepo: userRepo,
	}
}

// GetAccessToken implements FCMUseCase.
func (f *fcmUseCase) GetAccessToken(ctx context.Context) (*entity.Fcm, error) {
	serviceAccountPath := "internal/infrastructure/broadcast/donora-f67f2-5c889d5acd0a.json"
	data, err := os.ReadFile(serviceAccountPath)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return nil, err
	}

	token, err := conf.TokenSource(ctx).Token()
	if err != nil {
		return nil, err
	}

	// ambil project_id dari file JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	projectID, ok := parsed["project_id"].(string)
	if !ok {
		return nil, fmt.Errorf("project_id not found or invalid")
	}

	res := &entity.Fcm{
		AccessToken: token.AccessToken,
		ProjectID:   projectID,
	}

	return res, nil
}

// SendFCMV1 implements FCMUseCase.
func (f *fcmUseCase) SendFCMV1(ctx context.Context, userID uint, bloodType, title, body string) error {
	data, err := f.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	user, err := f.userRepo.FindById(ctx, userID)
	if err != nil {
		return err
	}

	if user.FCMToken == "" {
		err = errors.New("user fcm token not found")
		return err
	}

	if user.Role != "pencari" {
		err = errors.New("user role must be pencari")
		return err
	}

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", data.ProjectID)

	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"topic": "blood_" + bloodType,
			"notification": map[string]string{
				"title": title,
				"body":  body,
			},
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+data.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to send: %s", string(bodyBytes))
	}

	return nil
}

// func (f *fcmUseCase) SubscribeByBlood(ctx context.Context, fcmToken, bloodType string) error {
// 	opt := option.WithCredentialsFile("/home/fahrul/gdsc/gsc/backend/internal/infrastructure/broadcast/donora-f67f2-5c889d5acd0a.json")

// 	app, err := firebase.NewApp(ctx, nil, opt)
// 	if err != nil {
// 		return err
// 	}

// 	err = SubscribeUserToBloodTopic(ctx, app, fcmToken, bloodType)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (f *fcmUseCase) SubscribeUserToBloodTopic(ctx context.Context, app *firebase.App, fcmToken string, bloodType string) error {
	client, err := app.Messaging(ctx)
	if err != nil {
		return err
	}

	topic := "blood_" + bloodType
	_, err = client.SubscribeToTopic(ctx, []string{fcmToken}, topic)
	if err != nil {
		return err
	}

	return nil
}
