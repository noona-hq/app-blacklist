package services

type User struct {
	ID           string `json:"id"`
	CompanyID    string `json:"company_id"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
