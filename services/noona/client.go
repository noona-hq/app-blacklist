package noona

import (
	"context"
	"fmt"
	"time"

	"github.com/noona-hq/app-blacklist/utils"
	noona "github.com/noona-hq/noona-sdk-go"
	"github.com/pkg/errors"
)

const (
	BlacklistCustomerGroupName = "Blacklist"
)

type Client struct {
	cfg    Config
	Client *noona.ClientWithResponses
}

func (a Client) GetUser() (*noona.User, error) {
	userResponse, err := a.Client.GetUserWithResponse(context.Background(), &noona.GetUserParams{
		Expand: &noona.Expand{"companies"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error getting user")
	}

	if userResponse.StatusCode() != 200 {
		return nil, errors.New("Error getting user")
	}

	if userResponse.JSON200 == nil {
		return nil, errors.New("Error getting user")
	}

	return userResponse.JSON200, nil
}

func (a Client) SetupWebhook(companyID string) error {
	webhook := noona.Webhook{
		Title:       utils.StringPtr("Blacklist"),
		Description: utils.StringPtr("Watches event creation to enforce blacklist app functioinality."),
		CallbackUrl: utils.StringPtr(a.cfg.AppBaseURL + "/webhook"),
		Company: func() *noona.ExpandableCompany {
			company := noona.ExpandableCompany{}
			company.FromID(noona.ID(companyID))
			return &company
		}(),
		Enabled: utils.BoolPtr(true),
		Headers: &noona.WebhookHeaders{
			{
				Key:    utils.StringPtr("Authorization"),
				Values: &[]string{"Bearer " + a.cfg.AppWebhookToken},
			},
		},
		Events: &noona.WebhookEvents{
			noona.WebhookEventEventCreated,
		},
	}

	webhookResponse, err := a.Client.CreateWebhookWithResponse(context.Background(), &noona.CreateWebhookParams{}, noona.CreateWebhookJSONRequestBody(webhook))
	if err != nil {
		return errors.Wrap(err, "Error creating webhook")
	}

	if webhookResponse.StatusCode() != 200 {
		return errors.New("Error creating webhook")
	}

	return nil
}

func (a Client) SetupBlacklistCustomerGroup(companyID string) error {
	existingGroups, err := a.Client.ListCustomerGroupsWithResponse(context.Background(), companyID, &noona.ListCustomerGroupsParams{})
	if err != nil {
		return errors.Wrap(err, "Error listing customer groups")
	}

	if existingGroups.StatusCode() != 200 {
		return errors.New("Error listing customer groups")
	}

	for _, group := range *existingGroups.JSON200 {
		if group.Title != nil && *group.Title == BlacklistCustomerGroupName {
			// Group already exists
			return nil
		}
	}

	group := noona.CustomerGroup{
		Title:       utils.StringPtr(BlacklistCustomerGroupName),
		Description: utils.StringPtr("Customers in this group are blacklisted and all their online appointment requests will be automatically rejected."),
		Company:     &companyID,
	}

	groupResponse, err := a.Client.CreateCustomerGroupWithResponse(context.Background(), &noona.CreateCustomerGroupParams{}, noona.CreateCustomerGroupJSONRequestBody(group))
	if err != nil {
		return errors.Wrap(err, "Error creating customer group")
	}

	if groupResponse.StatusCode() != 200 {
		return errors.New("Error creating customer group")
	}

	return nil
}

func (a Client) ShouldBlacklistEvent(event noona.Event) (bool, *noona.Customer, error) {
	if event.Unconfirmed == nil || !*event.Unconfirmed {
		return false, nil, nil
	}

	if event.Customer == nil {
		return false, nil, nil
	}

	customerID, err := event.Customer.AsID()
	if err != nil {
		return false, nil, errors.Wrap(err, "Error getting customer id")
	}

	customerResponse, err := a.Client.GetCustomerWithResponse(context.Background(), string(customerID), &noona.GetCustomerParams{
		Expand: &noona.Expand{"groups"},
	})
	if err != nil {
		return false, nil, errors.Wrap(err, "Error getting customer")
	}

	if customerResponse.StatusCode() != 200 {
		return false, nil, errors.New("Error getting customer")
	}

	if customerResponse.JSON200.Groups == nil {
		// Customer has no groups
		return false, nil, nil
	}

	for _, group := range *customerResponse.JSON200.Groups {
		expandedGroup, err := group.AsCustomerGroup()
		if err != nil {
			continue
		}

		if expandedGroup.Title != nil && *expandedGroup.Title == BlacklistCustomerGroupName {
			return true, customerResponse.JSON200, nil
		}
	}

	return false, nil, nil
}

func (a Client) DenyEvent(event noona.Event, customer *noona.Customer) error {
	_, err := a.Client.UpdateEventWithResponse(context.Background(), *event.Id, &noona.UpdateEventParams{}, noona.UpdateEventJSONRequestBody{
		DeclinedAt: utils.TimePtr(time.Now()),
	})
	if err != nil {
		return errors.Wrap(err, "Error denying event")
	}

	if event.Employee == nil || event.Company == nil || customer == nil || (customer.Name == nil && customer.PhoneNumber == nil) {
		return nil
	}

	employeeID, err := event.Employee.AsID()
	if err != nil {
		return nil
	}

	companyID, err := event.Company.AsID()
	if err != nil {
		return nil
	}

	// Send notification to employee

	customerIdentifier := customer.Name
	if customerIdentifier == nil {
		customerIdentifier = customer.PhoneNumber
	}

	notification := noona.NotificationCreate{
		Company:  string(companyID),
		Employee: string(employeeID),
		Title:    "Blacklist triggered",
		Message:  fmt.Sprintf("Appointment request by %s rejected because customer is blacklisted", *customerIdentifier),
	}

	a.Client.CreateNotificationWithResponse(context.Background(), noona.CreateNotificationJSONRequestBody(notification))

	return nil
}
