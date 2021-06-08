package gitconfigerror

import (
	"fmt"
)

// usage: https://golang.org/pkg/errors/?source=post_page

var (
	ErrSectionOrKeyIsInvalid                  = &SectionOrKeyIsInvalidErr{message: "section or key is invalid"}
	ErrNoSectionOrNameProvided                = &NoSectionOrNameProvidedErr{message: "no section or name provided"}
	ErrConfigFileIsInvalid                    = &ConfigFileIsInvalidErr{message: "config file is invalid"}
	ErrConfigFileCannotBeWritten              = &ConfigFileCannotBeWrittenErr{message: "config file cannot be written"}
	ErrTryingToUnsetAnOptionWhichDoesNotExist = &TryingToUnsetAnOptionWhichDoesNotExistErr{message: "trying to unset an option which does not exist"}
	ErrTryingToUseAnInvalidRegexp             = &TryingToUseAnInvalidRegexpErr{message: "trying to use an invalid regexp"}
)

type SectionOrKeyIsInvalidErr struct {
	message string
}

func (e *SectionOrKeyIsInvalidErr) Error() string {
	return e.message
}

type NoSectionOrNameProvidedErr struct {
	message string
}

func (e *NoSectionOrNameProvidedErr) Error() string {
	return e.message
}

type ConfigFileIsInvalidErr struct {
	message string
}

func (e *ConfigFileIsInvalidErr) Error() string {
	return e.message
}

type ConfigFileCannotBeWrittenErr struct {
	message string
}

func (e *ConfigFileCannotBeWrittenErr) Error() string {
	return e.message
}

type TryingToUnsetAnOptionWhichDoesNotExistErr struct {
	message string
}

func (e *TryingToUnsetAnOptionWhichDoesNotExistErr) Error() string {
	return e.message
}

type TryingToUseAnInvalidRegexpErr struct {
	message string
}

func (e *TryingToUseAnInvalidRegexpErr) Error() string {
	return e.message
}

type UnknownErr struct {
	message string
}

func (e *UnknownErr) Error() string {
	return e.message
}

// what is actually returned from gitconfig command execution
const (
	msgSectionOrKeyInvalid                  = "exit status 1"
	msgNoSectionOrNameProvided              = "exit status 2"
	msgConfigFileInvalid                    = "exit status 3"
	msgConfigFileCannotBeWritten            = "exit status 4"
	msgTryingToUnsetOptionWhichDoesNotExist = "exit status 5"
	msgTryingToUseAnInvalidRegexp           = "exit status 6"
)

func New(err error) error {
	if err == nil {
		return nil
	}

	switch err.Error() {
	case msgSectionOrKeyInvalid:
		return ErrSectionOrKeyIsInvalid
	case msgNoSectionOrNameProvided:
		return ErrNoSectionOrNameProvided
	case msgConfigFileInvalid:
		return ErrConfigFileIsInvalid
	case msgConfigFileCannotBeWritten:
		return ErrConfigFileCannotBeWritten
	case msgTryingToUnsetOptionWhichDoesNotExist:
		return ErrTryingToUnsetAnOptionWhichDoesNotExist
	case msgTryingToUseAnInvalidRegexp:
		return ErrTryingToUseAnInvalidRegexp
	default:
		return &UnknownErr{message: fmt.Sprintf("unknown gitconfig error: %s", err)}
	}
}
