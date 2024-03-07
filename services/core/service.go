package core

import (
	"github.com/noona-hq/blacklist/logger"
	"github.com/noona-hq/blacklist/services/noona"
	"github.com/noona-hq/blacklist/services/store"
	"github.com/noona-hq/blacklist/services/store/entity"
	noonasdk "github.com/noona-hq/noona-sdk-go"
	"github.com/pkg/errors"
)

type Service struct {
	logger logger.Logger
	noona  noona.Service
	store  store.Store
}

func New(logger logger.Logger, noona noona.Service, store store.Store) Service {
	return Service{logger, noona, store}
}

func (s Service) OnboardUserToBlacklist(code string) (*noonasdk.User, error) {
	s.logger.Infow("Onboarding user to blacklist")

	noAuthClient, err := s.noona.NoAuthNoonaClient()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting no auth noona client")
	}

	token, err := noAuthClient.CodeTokenExchange(code)
	if err != nil {
		return nil, errors.Wrap(err, "Error exchanging code for token")
	}

	authClient, err := s.noona.AuthNoonaClient(*token)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting auth noona client")
	}

	user, err := authClient.GetUser()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting user")
	}

	blacklistUser, err := s.noonaUserAsBlacklistUser(user, token)
	if err != nil {
		return nil, errors.Wrap(err, "Error converting noona user to blacklist user")
	}

	if err := authClient.SetupWebhook(blacklistUser.CompanyID); err != nil {
		return nil, errors.Wrap(err, "Error setting up webhook")
	}

	if err := authClient.SetupBlacklistCustomerGroup(blacklistUser.CompanyID); err != nil {
		return nil, errors.Wrap(err, "Error setting up blacklist customer group")
	}

	if err := s.store.CreateBlacklistUser(blacklistUser); err != nil {
		return nil, errors.Wrap(err, "Error creating blacklist user")
	}

	s.logger.Infow("User onboarded to blacklist", "email", blacklistUser.Email, "company_id", blacklistUser.CompanyID)

	return user, nil
}

// ProcessWebhookCallback processes a webhook callback from Noona and enforces the blacklist
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

	user, err := s.store.GetBlacklistUserForCompany(string(companyID))
	if err != nil {
		s.logger.Errorw("Error getting blacklist user for company", "event_id", *event.Id, "company_id", string(companyID), "error", err)
		return nil
	}

	authClient, err := s.noona.AuthNoonaClientFromBlacklistUser(user)
	if err != nil {
		s.logger.Errorw("Error getting auth noona client from blacklist user", "event_id", *event.Id, "error", err)
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

func (s Service) noonaUserAsBlacklistUser(user *noonasdk.User, token *noonasdk.OAuthToken) (entity.User, error) {
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
