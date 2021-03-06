package ui

import (
	"fmt"
	"strings"

	"github.com/CircleCI-Public/circleci-cli/logger"
	"github.com/manifoldco/promptui"
)

// UserInterface is created to allow us to pass a mock user interface for testing.
// Since we can't properly run integration tests on code that calls PromptUI.
// This is because the first call to PromptUI reads from stdin correctly,
// but subsequent calls return EOF.
type UserInterface interface {
	ReadSecretStringFromUser(log *logger.Logger, message string) (string, error)
	ReadStringFromUser(log *logger.Logger, message string, defaultValue string) string
	AskUserToConfirm(log *logger.Logger, message string) bool
}

// InteractiveUI implements the UserInterface used by the real program, not in tests.
type InteractiveUI struct {
}

// ReadSecretStringFromUser can be used to read a value from the user by masking their input.
// It's useful for token input in our case.
func (InteractiveUI) ReadSecretStringFromUser(_ *logger.Logger, message string) (string, error) {
	prompt := promptui.Prompt{
		Label: message,
		Mask:  '*',
	}

	token, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return token, nil
}

// ReadStringFromUser can be used to read any value from the user or the defaultValue when provided.
func (InteractiveUI) ReadStringFromUser(_ *logger.Logger, message string, defaultValue string) string {
	prompt := promptui.Prompt{
		Label: message,
	}

	if defaultValue != "" {
		prompt.Default = defaultValue
	}

	token, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	return token
}

// AskUserToConfirm will prompt the user to confirm with the provided message.
func (InteractiveUI) AskUserToConfirm(_ *logger.Logger, message string) bool {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	return err == nil && strings.ToLower(result) == "y"
}

// TestingUI implements the UserInterface for our testing purposes.
type TestingUI struct {
	Input   string
	Confirm bool
}

// ReadSecretStringFromUser implements the TestingUI interface for asking a user's input.
// It works by simply printing the message to standard output and passing the input through.
func (ui TestingUI) ReadSecretStringFromUser(log *logger.Logger, message string) (string, error) {
	log.Info(message)
	return ui.Input, nil
}

// ReadStringFromUser implements the TestingUI interface for asking a user's input.
// It works by simply printing the message to standard output and passing the input through.
func (ui TestingUI) ReadStringFromUser(log *logger.Logger, message string, defaultValue string) string {
	log.Info(message)
	return ui.Input
}

// AskUserToConfirm works by printing the provided message to standard out and returning a Confirm dialogue up the chain.
func (ui TestingUI) AskUserToConfirm(log *logger.Logger, message string) bool {
	log.Info(message)
	return ui.Confirm
}

// ShouldAskForToken wraps an AskUserToConfirm dialogue only if their token is empty or blank.
func ShouldAskForToken(token string, log *logger.Logger, ui UserInterface) bool {
	if token == "" {
		return true
	}

	return ui.AskUserToConfirm(log, "A CircleCI token is already set. Do you want to change it")
}

// ShouldAskForEndpoint wraps an AskUserToConfirm dialogue only if their endpoint has changed from the default value.
func ShouldAskForEndpoint(endpoint string, log *logger.Logger, ui UserInterface, defaultValue string) bool {
	if endpoint == defaultValue {
		return true
	}

	return ui.AskUserToConfirm(log, fmt.Sprintf("Do you want to reset the endpoint? (default: %s)", defaultValue))
}
