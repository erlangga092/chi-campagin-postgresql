package user

type UserFormatter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{}
	formatter.ID = user.ID
	formatter.Name = user.Name
	formatter.Occupation = user.Occupation
	formatter.Email = user.Email
	formatter.Token = token

	return formatter
}
