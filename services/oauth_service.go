package services

type OAuthService interface {
	CreateUser(user User) error
	GetUserForCompany(companyID string) (User, error)
}
