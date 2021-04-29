package apollo

import "testing"

func Test_errConfigMissingField(t *testing.T) {
	err := errConfigMissingField("")
	if s := err.Error(); s == "" {
		t.Errorf("error has no text")
	}
}

func Test_errParamsMissingField(t *testing.T) {
	err := errParamsMissingField("")
	if s := err.Error(); s == "" {
		t.Errorf("error has no text")
	}
}
