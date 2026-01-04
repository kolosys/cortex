package cortex

import "time"

// EvalMode determines how errors are handled during evaluation.
type EvalMode int

const (
	// ModeFailFast stops evaluation on first error.
	ModeFailFast EvalMode = iota

	// ModeCollectAll evaluates all rules and collects all errors.
	ModeCollectAll

	// ModeContinueOnError logs errors but continues evaluation.
	ModeContinueOnError
)

func (m EvalMode) String() string {
	switch m {
	case ModeFailFast:
		return "fail_fast"
	case ModeCollectAll:
		return "collect_all"
	case ModeContinueOnError:
		return "continue_on_error"
	default:
		return "unknown"
	}
}

// Config configures the engine behavior.
type Config struct {
	// Mode determines error handling behavior.
	Mode EvalMode

	// Timeout for entire evaluation (0 = no timeout).
	Timeout time.Duration

	// ShortCircuit allows rules to halt evaluation early.
	ShortCircuit bool

	// EnableMetrics enables detailed metrics collection.
	EnableMetrics bool

	// MaxRules limits the number of rules (0 = unlimited).
	MaxRules int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Mode:          ModeFailFast,
		Timeout:       0,
		ShortCircuit:  true,
		EnableMetrics: true,
		MaxRules:      0,
	}
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Timeout < 0 {
		return ErrInvalidRule
	}
	if c.MaxRules < 0 {
		return ErrInvalidRule
	}
	return nil
}

func (c *Config) applyDefaults() {
	if c.Mode < ModeFailFast || c.Mode > ModeContinueOnError {
		c.Mode = ModeFailFast
	}
}
