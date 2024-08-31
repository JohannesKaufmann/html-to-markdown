package commonmark

import (
	"fmt"
	"testing"
)

func TestValidateConfig_Empty(t *testing.T) {
	cfg := fillInDefaultConfig(&config{})
	if cfg.HeadingStyle != "atx" {
		t.Error("the config value was not filled with the default value")
	}

	err := validateConfig(&cfg)
	if err != nil {
		t.Errorf("expected no error but got %+v", err)
	}
}
func TestValidateConfig_Success(t *testing.T) {
	cfg := fillInDefaultConfig(&config{
		HeadingStyle: "setext",
	})
	if cfg.HeadingStyle != "setext" {
		t.Error("the config value was overridden")
	}

	err := validateConfig(&cfg)
	if err != nil {
		t.Errorf("expected no error but got %+v", err)
	}
}
func TestValidateConfig_RandomValue(t *testing.T) {
	cfg := fillInDefaultConfig(&config{
		HeadingStyle: "random",
	})

	err := validateConfig(&cfg)
	if err == nil {
		t.Error("expected an error")
	}
	e, ok := err.(*ValidateConfigError)
	if !ok {
		t.Error("expected an error of type ValidateConfigError")
	}
	if e.Key != "HeadingStyle" {
		t.Errorf("expected a different value for 'key' but got %q", e.Key)
	}
	if e.Value != "random" {
		t.Errorf("expected a different value for 'actual' but got %q", e.Value)
	}

	formatted := err.Error()
	if formatted != "invalid value for HeadingStyle:\"random\" must be one of \"atx\" or \"setext\"" {
		t.Errorf("expected a different formatted message but got %q", formatted)
	}
}

func TestValidateConfig_KeyWithValue(t *testing.T) {
	cfg := fillInDefaultConfig(&config{
		StrongDelimiter: "*",
	})

	err := validateConfig(&cfg)
	if err == nil {
		t.Error("expected an error")
	}
	e, ok := err.(*ValidateConfigError)
	if !ok {
		t.Fatal("expected an error of type ValidateConfigError")
	}

	// The default error message for the golang api
	formatted1 := err.Error()
	expected1 := `invalid value for StrongDelimiter:"*" must be exactly 2 characters of "**" or "__"`
	if formatted1 != expected1 {
		t.Errorf("expected a different formatted message but got %q", formatted1)
	}

	// The error message for the cli
	if e.Key == "StrongDelimiter" {
		e.KeyWithValue = fmt.Sprintf("--%s=%q", "strong_delimiter", e.Value)
	}
	formatted2 := err.Error()
	expected2 := `invalid value for --strong_delimiter="*" must be exactly 2 characters of "**" or "__"`
	if formatted2 != expected2 {
		t.Errorf("expected a different formatted message but got %q", formatted2)
	}
}
