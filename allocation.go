package cortex

import (
	"context"
	"fmt"
	"math"
)

// AllocationStrategy defines how values are distributed.
type AllocationStrategy int

const (
	// StrategyPercentage distributes by percentage (must sum to 100).
	StrategyPercentage AllocationStrategy = iota

	// StrategyFixed distributes fixed amounts.
	StrategyFixed

	// StrategyWeighted distributes by relative weights.
	StrategyWeighted

	// StrategyEqual distributes equally among all targets.
	StrategyEqual

	// StrategyRatio distributes by ratio (e.g., 2:3:5).
	StrategyRatio
)

func (s AllocationStrategy) String() string {
	switch s {
	case StrategyPercentage:
		return "percentage"
	case StrategyFixed:
		return "fixed"
	case StrategyWeighted:
		return "weighted"
	case StrategyEqual:
		return "equal"
	case StrategyRatio:
		return "ratio"
	default:
		return "unknown"
	}
}

// ParseAllocationStrategy parses a string into an AllocationStrategy.
func ParseAllocationStrategy(s string) (AllocationStrategy, error) {
	switch s {
	case "percentage":
		return StrategyPercentage, nil
	case "fixed":
		return StrategyFixed, nil
	case "weighted":
		return StrategyWeighted, nil
	case "equal":
		return StrategyEqual, nil
	case "ratio":
		return StrategyRatio, nil
	default:
		return 0, fmt.Errorf("%w: unknown allocation strategy %q", ErrInvalidRule, s)
	}
}

// AllocationTarget specifies a single allocation destination.
type AllocationTarget struct {
	Key    string  // context key to set
	Amount float64 // percentage, fixed amount, weight, or ratio (based on strategy)
}

// AllocationRule distributes a value across multiple targets.
type AllocationRule struct {
	baseRule
	source    string
	strategy  AllocationStrategy
	targets   []AllocationTarget
	remainder string // optional: key for rounding remainder
	precision int    // decimal precision
}

// AllocationConfig configures an allocation rule.
type AllocationConfig struct {
	ID          string
	Name        string
	Description string
	Deps        []string

	// Source is the context key containing the value to allocate.
	Source string

	// Strategy is the distribution method.
	Strategy AllocationStrategy

	// Targets are the allocation destinations.
	Targets []AllocationTarget

	// Remainder is an optional context key for the rounding remainder.
	Remainder string

	// Precision is the decimal precision (default 2).
	Precision int
}

// NewAllocation creates a new allocation rule.
func NewAllocation(cfg AllocationConfig) (*AllocationRule, error) {
	if cfg.ID == "" {
		return nil, fmt.Errorf("%w: allocation rule requires ID", ErrInvalidRule)
	}
	if cfg.Source == "" {
		return nil, fmt.Errorf("%w: allocation rule %q requires source", ErrInvalidRule, cfg.ID)
	}
	if len(cfg.Targets) == 0 {
		return nil, fmt.Errorf("%w: allocation rule %q requires at least one target", ErrInvalidRule, cfg.ID)
	}

	// Validate based on strategy
	switch cfg.Strategy {
	case StrategyPercentage:
		var sum float64
		for _, t := range cfg.Targets {
			sum += t.Amount
		}
		if math.Abs(sum-100) > 0.0001 {
			return nil, fmt.Errorf("%w: allocation rule %q percentages sum to %.2f, not 100", ErrAllocationSum, cfg.ID, sum)
		}
	case StrategyEqual:
		// No amounts needed
	default:
		// Weighted, Ratio, Fixed: amounts should be positive
		for _, t := range cfg.Targets {
			if t.Amount < 0 {
				return nil, fmt.Errorf("%w: allocation rule %q has negative amount for %q", ErrInvalidRule, cfg.ID, t.Key)
			}
		}
	}

	precision := cfg.Precision
	if precision <= 0 {
		precision = 2
	}

	return &AllocationRule{
		baseRule: baseRule{
			id:          cfg.ID,
			name:        cfg.Name,
			description: cfg.Description,
			deps:        cfg.Deps,
		},
		source:    cfg.Source,
		strategy:  cfg.Strategy,
		targets:   cfg.Targets,
		remainder: cfg.Remainder,
		precision: precision,
	}, nil
}

// MustAllocation creates a new allocation rule, panicking on error.
func MustAllocation(cfg AllocationConfig) *AllocationRule {
	r, err := NewAllocation(cfg)
	if err != nil {
		panic(err)
	}
	return r
}

// Evaluate distributes the source value across targets.
func (r *AllocationRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error {
	source, err := evalCtx.GetFloat64(r.source)
	if err != nil {
		return NewRuleError(r.id, string(RuleTypeAllocation), "evaluate", err)
	}

	allocations, remainder := r.calculate(source)

	for i, t := range r.targets {
		evalCtx.Set(t.Key, allocations[i])
	}

	if r.remainder != "" && remainder != 0 {
		evalCtx.Set(r.remainder, remainder)
	}

	return nil
}

func (r *AllocationRule) calculate(source float64) ([]float64, float64) {
	n := len(r.targets)
	allocations := make([]float64, n)

	switch r.strategy {
	case StrategyPercentage:
		var total float64
		for i, t := range r.targets {
			allocations[i] = r.round(source * t.Amount / 100)
			total += allocations[i]
		}
		return allocations, source - total

	case StrategyFixed:
		var total float64
		for i, t := range r.targets {
			allocations[i] = r.round(t.Amount)
			total += allocations[i]
		}
		return allocations, source - total

	case StrategyWeighted:
		var totalWeight float64
		for _, t := range r.targets {
			totalWeight += t.Amount
		}
		if totalWeight == 0 {
			return allocations, source
		}
		var total float64
		for i, t := range r.targets {
			allocations[i] = r.round(source * t.Amount / totalWeight)
			total += allocations[i]
		}
		return allocations, source - total

	case StrategyEqual:
		each := r.round(source / float64(n))
		var total float64
		for i := range r.targets {
			allocations[i] = each
			total += each
		}
		return allocations, source - total

	case StrategyRatio:
		var totalRatio float64
		for _, t := range r.targets {
			totalRatio += t.Amount
		}
		if totalRatio == 0 {
			return allocations, source
		}
		var total float64
		for i, t := range r.targets {
			allocations[i] = r.round(source * t.Amount / totalRatio)
			total += allocations[i]
		}
		return allocations, source - total
	}

	return allocations, 0
}

func (r *AllocationRule) round(v float64) float64 {
	multiplier := math.Pow(10, float64(r.precision))
	return math.Round(v*multiplier) / multiplier
}

// Source returns the source key.
func (r *AllocationRule) Source() string {
	return r.source
}

// Strategy returns the allocation strategy.
func (r *AllocationRule) Strategy() AllocationStrategy {
	return r.strategy
}

// Targets returns the allocation targets.
func (r *AllocationRule) Targets() []AllocationTarget {
	return r.targets
}
