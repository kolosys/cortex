package cortex

import (
	"context"
	"fmt"
	"math"
)

// Lookup represents a lookup table.
type Lookup interface {
	// Name returns the lookup table name.
	Name() string

	// Get retrieves a value by key.
	Get(key any) (any, bool)
}

// MapLookup provides a simple map-based lookup implementation.
type MapLookup[K comparable, V any] struct {
	name  string
	items map[K]V
}

// NewMapLookup creates a new map-based lookup table.
func NewMapLookup[K comparable, V any](name string, items map[K]V) *MapLookup[K, V] {
	cp := make(map[K]V, len(items))
	for k, v := range items {
		cp[k] = v
	}
	return &MapLookup[K, V]{
		name:  name,
		items: cp,
	}
}

func (l *MapLookup[K, V]) Name() string { return l.name }

func (l *MapLookup[K, V]) Get(key any) (any, bool) {
	k, ok := key.(K)
	if !ok {
		return nil, false
	}
	v, found := l.items[k]
	return v, found
}

// RangeEntry represents a single range in a range lookup.
type RangeEntry[V any] struct {
	Min   float64 // inclusive
	Max   float64 // exclusive (use math.Inf(1) for unbounded)
	Value V
}

// RangeLookup provides range-based lookup (e.g., tax brackets).
type RangeLookup[V any] struct {
	name   string
	ranges []RangeEntry[V]
}

// NewRangeLookup creates a new range-based lookup table.
func NewRangeLookup[V any](name string, ranges []RangeEntry[V]) *RangeLookup[V] {
	cp := make([]RangeEntry[V], len(ranges))
	copy(cp, ranges)
	return &RangeLookup[V]{
		name:   name,
		ranges: cp,
	}
}

func (l *RangeLookup[V]) Name() string { return l.name }

func (l *RangeLookup[V]) Get(key any) (any, bool) {
	k, err := toFloat64(key)
	if err != nil {
		return nil, false
	}
	for _, r := range l.ranges {
		if k >= r.Min && k < r.Max {
			return r.Value, true
		}
	}
	return nil, false
}

// LookupRule retrieves a value from a lookup table.
type LookupRule struct {
	baseRule
	table      string
	keySource  string // context key to use as lookup key
	target     string
	defaultVal any
	required   bool
}

// LookupConfig configures a lookup rule.
type LookupConfig struct {
	ID          string
	Name        string
	Description string
	Deps        []string

	// Table is the lookup table name (must be registered).
	Table string

	// Key is the context key to use as lookup key.
	Key string

	// Target is the context key to store the result.
	Target string

	// Default is the value if not found (ignored if Required is true).
	Default any

	// Required causes an error if the lookup key is not found.
	Required bool
}

// NewLookup creates a new lookup rule.
func NewLookup(cfg LookupConfig) (*LookupRule, error) {
	if cfg.ID == "" {
		return nil, fmt.Errorf("%w: lookup rule requires ID", ErrInvalidRule)
	}
	if cfg.Table == "" {
		return nil, fmt.Errorf("%w: lookup rule %q requires table", ErrInvalidRule, cfg.ID)
	}
	if cfg.Key == "" {
		return nil, fmt.Errorf("%w: lookup rule %q requires key", ErrInvalidRule, cfg.ID)
	}
	if cfg.Target == "" {
		return nil, fmt.Errorf("%w: lookup rule %q requires target", ErrInvalidRule, cfg.ID)
	}

	return &LookupRule{
		baseRule: baseRule{
			id:          cfg.ID,
			name:        cfg.Name,
			description: cfg.Description,
			deps:        cfg.Deps,
		},
		table:      cfg.Table,
		keySource:  cfg.Key,
		target:     cfg.Target,
		defaultVal: cfg.Default,
		required:   cfg.Required,
	}, nil
}

// MustLookup creates a new lookup rule, panicking on error.
func MustLookup(cfg LookupConfig) *LookupRule {
	r, err := NewLookup(cfg)
	if err != nil {
		panic(err)
	}
	return r
}

// Evaluate performs the lookup and sets the result.
func (r *LookupRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error {
	key, ok := evalCtx.Get(r.keySource)
	if !ok {
		return NewRuleError(r.id, string(RuleTypeLookup), "evaluate",
			fmt.Errorf("%w: %s", ErrValueNotFound, r.keySource))
	}

	value, found, err := evalCtx.Lookup(r.table, key)
	if err != nil {
		return NewRuleError(r.id, string(RuleTypeLookup), "evaluate", err)
	}

	if !found {
		if r.required {
			return NewRuleError(r.id, string(RuleTypeLookup), "evaluate",
				fmt.Errorf("%w: %v in table %s", ErrKeyNotFound, key, r.table))
		}
		value = r.defaultVal
	}

	evalCtx.Set(r.target, value)
	return nil
}

// Table returns the lookup table name.
func (r *LookupRule) Table() string {
	return r.table
}

// TaxBracket is a convenience type for common tax bracket lookups.
type TaxBracket struct {
	Min  float64
	Max  float64
	Rate float64
}

// NewTaxBracketLookup creates a range lookup for tax brackets.
func NewTaxBracketLookup(name string, brackets []TaxBracket) *RangeLookup[float64] {
	ranges := make([]RangeEntry[float64], len(brackets))
	for i, b := range brackets {
		max := b.Max
		if max == 0 {
			max = math.Inf(1)
		}
		ranges[i] = RangeEntry[float64]{
			Min:   b.Min,
			Max:   max,
			Value: b.Rate,
		}
	}
	return NewRangeLookup(name, ranges)
}
