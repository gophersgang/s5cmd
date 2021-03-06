package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// AcceptableError interface defines an error which is OK-to-have, for things like "cp -n" etc. It should not be treated as an error (regarding the exit code etc)
type AcceptableError interface {
	error
	Acceptable() bool
}

// AcceptableErrorType embeds the stdlib error interface so that we can have more receivers on it
type AcceptableErrorType struct {
	error
}

// NewAcceptableError creates a new AcceptableError
func NewAcceptableError(s string) AcceptableErrorType {
	return AcceptableErrorType{errors.New(s)}
}

// Acceptable is always true for errors of AcceptableError type
func (e AcceptableErrorType) Acceptable() bool {
	return true
}

// IsRetryableError returns if an error (probably awserr) is retryable, along with an error code
func IsRetryableError(err error) (string, bool) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			//fmt.Println("awsErr", awsErr.Code(), awsErr.Message(), awsErr.OrigErr())

			errCode := awsErr.Code()
			switch errCode {
			case "SlowDown", "SerializationError":
				return errCode, true
			}

			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				//fmt.Println("reqErr", reqErr.StatusCode(), reqErr.RequestID())
				errCode = reqErr.Code()
				switch errCode {
				case "InternalError", "SerializationError":
					return errCode, true
				}
				status := reqErr.StatusCode()
				switch status {
				case 400, 500:
					return fmt.Sprintf("HTTP%d", status), true
				}
			}
		}
	}
	return "", false
}

// CleanupError converts multiline error messages generated by aws-sdk-go into a single line
func CleanupError(err error) (s string) {
	s = strings.Replace(err.Error(), "\n", " ", -1)
	s = strings.Replace(s, "\t", " ", -1)
	s = strings.Replace(s, "  ", " ", -1)
	s = strings.Replace(s, "  ", " ", -1)
	s = strings.TrimSpace(s)
	return
}

// IsAcceptableError determines if the error is an AcceptableError, and if so, returns the error as such
func IsAcceptableError(err error) AcceptableError {
	e, ok := err.(AcceptableError)
	if !ok {
		return nil
	}
	return e
}
