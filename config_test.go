package cortex_test

import (
	"testing"

	"github.com/kolosys/cortex"
)

func TestEvalModeString(t *testing.T) {
	tests := []struct {
		mode     cortex.EvalMode
		expected string
	}{
		{cortex.ModeFailFast, "fail_fast"},
		{cortex.ModeCollectAll, "collect_all"},
		{cortex.ModeContinueOnError, "continue_on_error"},
		{cortex.EvalMode(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.mode.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.mode.String())
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := cortex.DefaultConfig()

	if cfg.Mode != cortex.ModeFailFast {
		t.Errorf("expected ModeFailFast, got %v", cfg.Mode)
	}
	if cfg.Timeout != 0 {
		t.Errorf("expected Timeout=0, got %v", cfg.Timeout)
	}
	if !cfg.ShortCircuit {
		t.Error("expected ShortCircuit=true")
	}
	if !cfg.EnableMetrics {
		t.Error("expected EnableMetrics=true")
	}
	if cfg.MaxRules != 0 {
		t.Errorf("expected MaxRules=0, got %d", cfg.MaxRules)
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := cortex.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Invalid timeout
	cfg.Timeout = -1
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative timeout")
	}

	// Reset and test invalid MaxRules
	cfg = cortex.DefaultConfig()
	cfg.MaxRules = -1
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative max rules")
	}
}
