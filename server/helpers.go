package server

import (
	noona "github.com/noona-hq/noona-sdk-go"
)

func getCompanyNameFromUser(user *noona.User) (name string) {
	name = "company"

	if user.Companies == nil || len(*user.Companies) == 0 {
		return
	}

	company, err := (*user.Companies)[0].AsCompany()
	if err != nil {
		return
	}

	if company.Name != nil {
		name = *company.Name
	}

	return
}
