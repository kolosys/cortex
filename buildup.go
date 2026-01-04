package cortex

import (
	"context"
	"fmt"
	"math"
	"sync"
)

// BuildupOperation defines how values are accumulated.
type BuildupOperation int

const (
	// BuildupSum adds values together.
	BuildupSum BuildupOperation = iota

	// BuildupMin keeps the minimum value.
	BuildupMin

	// BuildupMax keeps the maximum value.
	BuildupMax

	// BuildupAvg computes running average.
	BuildupAvg

	// BuildupCount counts occurrences.
	BuildupCount

	// BuildupProduct multiplies values together.
	BuildupProduct
)

func (op BuildupOperation) String() string {
	switch op {
	case BuildupSum:
		return "sum"
	case BuildupMin:
		return "min"
	case BuildupMax:
		return "max"
	case BuildupAvg:
		return "avg"
	case BuildupCount:
		return "count"
	case BuildupProduct:
		return "product"
	default:
		return "unknown"
	}
}

// ParseBuildupOperation parses a string into a BuildupOperation.
func ParseBuildupOperation(s string) (BuildupOperation, error) {
	switch s {
	case "sum":
		return BuildupSum, nil
	case "min":
		return BuildupMin, nil
	case "max":
		return BuildupMax, nil
	case "avg", "average":
		return BuildupAvg, nil
	case "count":
		return BuildupCount, nil
	case "product":
		return BuildupProduct, nil
	default:
		return 0, fmt.Errorf("%w: unknown buildup operation %q", ErrInvalidRule, s)
	}
}

// Buildup represents a running accumulator.
type Buildup struct {
	Name      string
	Operation BuildupOperation

	mu    sync.Mutex
	value float64
	count int64
}

// Add adds a value to the buildup.
func (b *Buildup) Add(value float64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.count++

	switch b.Operation {
	case BuildupSum:
		b.value += value
	case BuildupMin:
		if b.count == 1 || value < b.value {
			b.value = value
		}
	case BuildupMax:
		if b.count == 1 || value > b.value {
			b.value = value
		}
	case BuildupAvg:
		b.value += value
	case BuildupCount:
		b.value = float64(b.count)
	case BuildupProduct:
		if b.count == 1 {
			b.value = value
		} else {
			b.value *= value
		}
	}
}

// Current returns the current accumulated value.
func (b *Buildup) Current() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.Operation == BuildupAvg && b.count > 0 {
		return b.value / float64(b.count)
	}
	return b.value
}

// Count returns the number of values added.
func (b *Buildup) Count() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.count
}

// Reset resets the buildup to its initial state.
func (b *Buildup) Reset(initial float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.value = initial
	b.count = 0
}

// BuildupRule accumulates values (running totals, aggregations).
type BuildupRule struct {
	baseRule
	buildup   string
	operation BuildupOperation
	source    string // context key containing value to add
	initial   float64
	target    string // optional: write current value to this key after adding
}

// BuildupConfig configures a buildup rule.
type BuildupConfig struct {
	ID          string
	Name        string
	Description string
	Deps        []string

	// Buildup is the buildup accumulator name (created if not exists).
	Buildup string

	// Operation is the accumulation operation.
	Operation BuildupOperation

	// Source is the context key containing the value to add.
	Source string

	// Initial is the initial value for the buildup (if created).
	Initial float64

	// Target is an optional context key to write the current value after adding.
	Target string
}

// NewBuildup creates a new buildup rule.
func NewBuildup(cfg BuildupConfig) (*BuildupRule, error) {
	if cfg.ID == "" {
		return nil, fmt.Errorf("%w: buildup rule requires ID", ErrInvalidRule)
	}
	if cfg.Buildup == "" {
		return nil, fmt.Errorf("%w: buildup rule %q requires buildup name", ErrInvalidRule, cfg.ID)
	}
	if cfg.Source == "" && cfg.Operation != BuildupCount {
		return nil, fmt.Errorf("%w: buildup rule %q requires source (except for count)", ErrInvalidRule, cfg.ID)
	}

	// Set sensible initial values based on operation
	initial := cfg.Initial
	if initial == 0 {
		switch cfg.Operation {
		case BuildupMin:
			initial = math.Inf(1)
		case BuildupMax:
			initial = math.Inf(-1)
		case BuildupProduct:
			initial = 1
		}
	}

	return &BuildupRule{
		baseRule: baseRule{
			id:          cfg.ID,
			name:        cfg.Name,
			description: cfg.Description,
			deps:        cfg.Deps,
		},
		buildup:   cfg.Buildup,
		operation: cfg.Operation,
		source:    cfg.Source,
		initial:   initial,
		target:    cfg.Target,
	}, nil
}

// MustBuildup creates a new buildup rule, panicking on error.
func MustBuildup(cfg BuildupConfig) *BuildupRule {
	r, err := NewBuildup(cfg)
	if err != nil {
		panic(err)
	}
	return r
}

// Evaluate adds to the buildup accumulator.
func (r *BuildupRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error {
	var value float64

	if r.operation == BuildupCount {
		value = 1
	} else {
		v, err := evalCtx.GetFloat64(r.source)
		if err != nil {
			return NewRuleError(r.id, string(RuleTypeBuildup), "evaluate", err)
		}
		value = v
	}

	b := evalCtx.GetOrCreateBuildup(r.buildup, r.operation, r.initial)
	b.Add(value)

	if r.target != "" {
		evalCtx.Set(r.target, b.Current())
	}

	return nil
}

// BuildupName returns the buildup accumulator name.
func (r *BuildupRule) BuildupName() string {
	return r.buildup
}
