package config

import (
	"testing"
)

func TestNewFeatureFlags(t *testing.T) {
	flags := NewFeatureFlags()

	if flags.IsParsingModuleEnabled() {
		t.Error("Expected parsing module to be disabled by default")
	}

	if flags.IsTelegramModuleEnabled() {
		t.Error("Expected telegram module to be disabled by default")
	}

	if flags.IsEventSourcingEnabled() {
		t.Error("Expected event sourcing to be disabled by default")
	}
}

func TestFeatureFlags_IsParsingModuleEnabled(t *testing.T) {
	flags := &FeatureFlags{UseNewParsingModule: true}

	if !flags.IsParsingModuleEnabled() {
		t.Error("Expected parsing module to be enabled")
	}
}
