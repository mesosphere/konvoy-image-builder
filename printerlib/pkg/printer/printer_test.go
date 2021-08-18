package printer

import (
	"bytes"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

type vError struct {
	Action string
	Change string
	Reason string
	Err    *field.Error
}

func (e vError) GetAction() string {
	return e.Action
}

func (e vError) GetChange() string {
	return e.Change
}

func (e vError) GetReason() string {
	return e.Reason
}

func (e vError) GetErr() *field.Error {
	return e.Err
}

func TestPrintValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		errors  ValidationErrors
		wantOut string
	}{
		{
			name:    "empty errors prints nothing",
			errors:  nil,
			wantOut: "",
		},
		{
			name:    "null error prints new line",
			errors:  ValidationErrors{vError{}},
			wantOut: "\n",
		},
		{
			name: "error prints all non empty filed",
			errors: ValidationErrors{
				vError{
					Action: "what need to be done",
				},
				vError{
					Action: "what need to be done",
					Change: "what need to be changed",
					Reason: "why",
					Err:    nil,
				},
				vError{
					Action: "what need to be done",
					Change: "what need to be changed",
					Reason: "why",
					Err:    field.Invalid(field.NewPath("path"), "value", "details"),
				},
			},
			wantOut: "Action required: what need to be done\n\n" +
				"Action required: what need to be done\n" +
				"Change: what need to be changed\n" +
				"Reason: why\n\n" +
				"Action required: what need to be done\n" +
				"Change: what need to be changed\n" +
				"Reason: why\n" +
				"Error:  path: Invalid value: \"value\": details\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintValidationErrors(out, tt.errors)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("PrintValidationErrors() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
