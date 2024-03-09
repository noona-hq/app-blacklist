package core

import (
	"time"

	"github.com/noona-hq/blacklist/logger"
	"github.com/noona-hq/blacklist/services/noona"
	"github.com/noona-hq/blacklist/store"
	"github.com/noona-hq/blacklist/store/entity"
	noonasdk "github.com/noona-hq/noona-sdk-go"
	"github.com/pkg/errors"
)

type Service struct {
	logger          logger.Logger
	noona           noona.Service
	store           store.Store
	anonymousClient noona.AnonymousClient
}

func New(logger logger.Logger, noona noona.Service, store store.Store) (Service, error) {
	anonymousClient, err := noona.AnonymousClient()
	if err != nil {
		return Service{}, errors.Wrap(err, "Error creating anonymous Noona client")
	}

	return Service{logger, noona, store, anonymousClient}, nil
}

func (s Service) OnboardUser(code string) (*noonasdk.User, error) {
	s.logger.Infow("Onboarding user to app")

	token, err := s.anonymousClient.CodeTokenExchange(code)
	if err != nil {
		return nil, errors.Wrap(err, "error exchanging code for token")
	}

	client, err := s.noona.Client(*token)
	if err != nil {
		return nil, errors.Wrap(err, "error getting auth noona client")
	}

	noonaUser, err := client.GetUser()
	if err != nil {
		return nil, errors.Wrap(err, "error getting user")
	}

	user, err := s.noonaUserAsUser(noonaUser, token)
	if err != nil {
		return nil, errors.Wrap(err, "error converting noona user to user")
	}

	if err := s.scaffoldNoonaResourcesForApp(client, user.CompanyID); err != nil {
		return nil, errors.Wrap(err, "error scaffolding noona resources")
	}

	if err := s.store.CreateUser(user); err != nil {
		return nil, errors.Wrap(err, "error creating user")
	}

	s.logger.Infow("User onboarded to app", "email", user.Email, "company_id", user.CompanyID)

	return noonaUser, nil
}

func (s Service) scaffoldNoonaResourcesForApp(client noona.Client, companyID string) error {
	if err := client.SetupWebhook(companyID); err != nil {
		return errors.Wrap(err, "error setting up webhook")
	}

	if err := client.SetupBlacklistCustomerGroup(companyID); err != nil {
		return errors.Wrap(err, "error setting up blacklist customer group")
	}

	return nil
}

// ProcessWebhookCallback processes a webhook callback from Noona
// Returning an error will cause the webhook to be retried
// Returning nil will acknowledge the webhook
func (s Service) ProcessWebhookCallback(callback noonasdk.CallbackData) error {
	event, err := callback.Data.AsEvent()
	if err != nil {
		return errors.Wrap(err, "Error getting event from callback data")
	}

	s.logger.Infow("Webhook callback received", "type", callback.Type, "event_id", *event.Id)

	if event.Unconfirmed == nil || !*event.Unconfirmed {
		s.logger.Infow("Event is confirmed, will not enforce blacklist", "event_id", *event.Id)
		return nil
	}

	companyID, err := event.Company.AsID()
	if err != nil {
		s.logger.Errorw("Error getting company id from event", "event_id", *event.Id, "error", err)
		return nil
	}

	user, err := s.store.GetUserForCompany(string(companyID))
	if err != nil {
		s.logger.Errorw("Error getting user for company", "event_id", *event.Id, "company_id", string(companyID), "error", err)
		return nil
	}

	oAuthToken, err := s.getOAuthTokenFromUser(user)
	if err != nil {
		s.logger.Errorw("Error getting OAuth token from user", "event_id", *event.Id, "error", err)
		return nil
	}

	authClient, err := s.noona.Client(oAuthToken)
	if err != nil {
		s.logger.Errorw("Error getting noona client", "event_id", *event.Id, "error", err)
		return nil
	}

	shouldBlackListEvent, customer, err := authClient.ShouldBlacklistEvent(event)
	if err != nil {
		s.logger.Errorw("Error checking if event should be blacklisted", "event_id", *event.Id, "error", err)
		return nil
	}

	if shouldBlackListEvent {
		if err := authClient.DenyEvent(event, customer); err != nil {
			s.logger.Errorw("Error blacklisting event", "event_id", *event.Id, "error", err)
			return nil
		}

		s.logger.Infow("Event blacklisted", "event_id", *event.Id)
		return nil
	}

	s.logger.Infow("Event not blacklisted", "event_id", *event.Id)

	return nil
}

func (s Service) noonaUserAsUser(user *noonasdk.User, token *noonasdk.OAuthToken) (entity.User, error) {
	if user == nil || token == nil {
		return entity.User{}, errors.New("user or token is nil")
	}

	if user.Companies == nil || len(*user.Companies) == 0 {
		return entity.User{}, errors.New("user has no associated companies")
	}

	company, err := (*user.Companies)[0].AsCompany()
	if err != nil {
		return entity.User{}, errors.Wrap(err, "error getting company")
	}

	return entity.User{
		Email:     *user.Email,
		CompanyID: *company.Id,
		Token: entity.Token{
			AccessToken:          *token.AccessToken,
			AccessTokenExpiresAt: *token.ExpiresAt,
			RefreshToken:         *token.RefreshToken,
		},
	}, nil
}

func (s Service) getOAuthTokenFromUser(user entity.User) (noonasdk.OAuthToken, error) {
	oAuthToken := noonasdk.OAuthToken{
		RefreshToken: &user.Token.RefreshToken,
		AccessToken:  &user.Token.AccessToken,
		ExpiresAt:    &user.Token.AccessTokenExpiresAt,
	}

	if oAuthToken.ExpiresAt.Before(time.Now().Add(time.Minute * 5)) {
		token, err := s.anonymousClient.RefreshTokenExchange(user.Token.RefreshToken)
		if err != nil {
			return noonasdk.OAuthToken{}, errors.Wrap(err, "error refreshing token")
		}

		oAuthToken = noonasdk.OAuthToken{
			RefreshToken: token.RefreshToken,
			AccessToken:  token.AccessToken,
			ExpiresAt:    token.ExpiresAt,
		}

		if _, err := s.store.UpdateUser(user.ID, entity.User{Token: entity.Token{
			AccessToken:          *token.AccessToken,
			AccessTokenExpiresAt: *token.ExpiresAt,
		}}); err != nil {
			s.logger.Errorw("Error updating user", "error", err)
		}
	}

	return oAuthToken, nil
}
