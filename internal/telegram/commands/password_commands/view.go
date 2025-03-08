package password_commands

import (
	"fmt"

	"github.com/vandi37/password-manager/internal/postgresql/module"
)

func ToString(passwords []module.Password, msg string) (string, bool) {
	if len(passwords) <= 0 {
		return "No passwords found", false
	}

	var s = "Passwords" + msg + "\n"

	for i, p := range passwords {
		s += fmt.Sprintf("\n%d. `%s`:`%s`", i+1, p.Company, p.Username)
	}

	return s + "\n\nEnter index of password to do some actions with it", true
}
