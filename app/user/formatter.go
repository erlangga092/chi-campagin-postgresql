package user

import "funding-app/app/key"

type UserFormatter struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Occupation   string `json:"occupation"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func FormatUser(user User, token key.Token) UserFormatter {
	formatter := UserFormatter{}
	formatter.ID = user.ID
	formatter.Name = user.Name
	formatter.Occupation = user.Occupation
	formatter.Email = user.Email
	formatter.AccessToken = token.AccessToken
	formatter.RefreshToken = token.RefreshToken

	return formatter
}
