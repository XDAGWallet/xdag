package cli

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
)

var validatePwd = func(input string) error {
	if len(input) < 1 {
		return errors.New("password must have at least 1 characters")
	}
	return nil
}

func ShowPassword() string {
	prompt := promptui.Prompt{
		Label:    "Password",
		Validate: validatePwd,
		Mask:     '*',
	}

	result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input password %v\n", err)
		result, err = prompt.Run()
	}

	return result
}

func reshowPassword(pwd string) bool {
	prompt := promptui.Prompt{
		Label:    "Confirm password",
		Validate: validatePwd,
		Mask:     '*',
	}

	result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input password failed %v\n", err)
		result, err = prompt.Run()
	}

	if pwd != result {
		fmt.Println("Password not match")
		return false
	}
	WalletAccount.Password = pwd
	return true
}
