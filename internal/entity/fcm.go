package entity

type Fcm struct {
	AccessToken string `json:"access_token"`
	ProjectID   string `json:"project_id"`
}

type RequestFcm struct {
	UserID    uint   `json:"user_id"`
	BloodType string `json:"blood_type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
}
