package printer

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type ValidationErrors []ValidationError

type ValidationError interface {
	// Action represents the action that was taken to cause the error
	GetAction() string
	// Change represents the change that must happen to remove the error
	GetChange() string
	// Reason represents the reason for the error happening
	GetReason() string
	// Err is the actual error message that triggered this error
	GetErr() *field.Error
}
