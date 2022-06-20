package internal

import (
	"fmt"
	"regexp"
)

func IsEmailValid(e string) bool {

	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

type Credentials struct {
	Email    string `form:"Email" json:"Email"`
	Password string `form:"Password" json:"Password"`

	// for the claims
	UserID   int64  `form:"UserID" json:"UserID"`
	AdminID  int64  `form:"AdminID" json:"AdminID"`
}

// for login function
// first arg USER ID, second arg ADMIN ID
func ValidateCredentials(cred Credentials) (int64, int64, error) {
	if cred.Email == "" {
		return 0,0,fmt.Errorf("Email invalid")
	}
	if cred.Password == "" {
		return 0,0,fmt.Errorf("Password invalid")
	}

	admin := Admin{}
	admin.UseDBWriterPreferred()
	notFound, err := admin.GetByPassword(cred.Password, cred.Email)
	if notFound {

		// admin is null
		user := User{}
		user.UseDBWriterPreferred()
		notFound, err = user.GetByPassword(cred.Password, cred.Email)
		if notFound {
			return 0,0,fmt.Errorf("User is nil")
		}
		if err != nil {
			return 0,0,err
		}
		return user.ID,0, nil
	}
	if err != nil {
		return 0,0,err
	}

	return 0, admin.ID, nil
}
