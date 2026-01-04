// Package cortex provides a rules engine for business logic evaluation.
//
// Cortex supports five rule types:
//   - Assignment: Set values directly on the context
//   - Formula: Calculate values using expressions or functions
//   - Allocation: Distribute values across multiple targets
//   - Lookup: Retrieve values from lookup tables
//   - Buildup: Accumulate/aggregate values (running totals, sums, etc.)
//
// Example:
//
//	engine := cortex.New("payroll", cortex.DefaultConfig())
//	engine.RegisterLookup(cortex.NewRangeLookup("tax_brackets", brackets))
//	engine.AddRules(
//	    cortex.MustAssignment(cortex.AssignmentConfig{...}),
//	    cortex.MustFormula(cortex.FormulaConfig{...}),
//	)
//	result, err := engine.Evaluate(ctx, cortex.NewEvalContext())
package cortex

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Logger provides a simple logging interface.
type Logger interface {
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, err error, kv ...any)
}

// Metrics provides a simple metrics interface.
type Metrics interface {
	Inc(name string, kv ...any)
	Add(name string, v float64, kv ...any)
	Histogram(name string, v float64, kv ...any)
}

// Tracer provides a simple tracing interface.
type Tracer interface {
	Start(ctx context.Context, name string, kv ...any) (context.Context, func(err error))
}

// Observability holds observability hooks for the engine.
type Observability struct {
	Logger  Logger
	Metrics Metrics
	Tracer  Tracer
}

type nopLogger struct{}

func (nopLogger) Debug(msg string, kv ...any)            {}
func (nopLogger) Info(msg string, kv ...any)             {}
func (nopLogger) Warn(msg string, kv ...any)             {}
func (nopLogger) Error(msg string, err error, kv ...any) {}

type nopMetrics struct{}

func (nopMetrics) Inc(name string, kv ...any)                  {}
func (nopMetrics) Add(name string, v float64, kv ...any)       {}
func (nopMetrics) Histogram(name string, v float64, kv ...any) {}

type nopTracer struct{}

func (nopTracer) Start(ctx context.Context, name string, kv ...any) (context.Context, func(err error)) {
	return ctx, func(err error) {}
}

func newObservability() *Observability {
	return &Observability{
		Logger:  nopLogger{},
		Metrics: nopMetrics{},
		Tracer:  nopTracer{},
	}
}

// Engine evaluates rules in sequence.
type Engine struct {
	name   string
	config *Config
	obs    *Observability
	closed atomic.Bool

	mu      sync.RWMutex
	rules   []Rule
	ruleIDs map[string]struct{}
	lookups map[string]Lookup
}

// New creates a new rules engine.
func New(name string, config *Config) *Engine {
	if config == nil {
		config = DefaultConfig()
	}
	config.applyDefaults()

	return &Engine{
		name:    name,
		config:  config,
		obs:     newObservability(),
		rules:   make([]Rule, 0),
		ruleIDs: make(map[string]struct{}),
		lookups: make(map[string]Lookup),
	}
}

// Name returns the engine name.
func (e *Engine) Name() string {
	return e.name
}

// Config returns the engine configuration.
func (e *Engine) Config() *Config {
	return e.config
}

// WithObservability sets observability hooks.
func (e *Engine) WithObservability(obs *Observability) *Engine {
	if obs != nil {
		if obs.Logger != nil {
			e.obs.Logger = obs.Logger
		}
		if obs.Metrics != nil {
			e.obs.Metrics = obs.Metrics
		}
		if obs.Tracer != nil {
			e.obs.Tracer = obs.Tracer
		}
	}
	return e
}

// AddRule adds a rule to the engine.
func (e *Engine) AddRule(rule Rule) error {
	if e.closed.Load() {
		return ErrEngineClosed
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.ruleIDs[rule.ID()]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateRule, rule.ID())
	}

	if e.config.MaxRules > 0 && len(e.rules) >= e.config.MaxRules {
		return fmt.Errorf("%w: max rules limit reached (%d)", ErrInvalidRule, e.config.MaxRules)
	}

	e.rules = append(e.rules, rule)
	e.ruleIDs[rule.ID()] = struct{}{}
	return nil
}

// AddRules adds multiple rules to the engine.
func (e *Engine) AddRules(rules ...Rule) error {
	for _, rule := range rules {
		if err := e.AddRule(rule); err != nil {
			return err
		}
	}
	return nil
}

// RegisterLookup registers a lookup table.
func (e *Engine) RegisterLookup(lookup Lookup) error {
	if e.closed.Load() {
		return ErrEngineClosed
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.lookups[lookup.Name()]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateLookup, lookup.Name())
	}

	e.lookups[lookup.Name()] = lookup
	return nil
}

// RegisterLookups registers multiple lookup tables.
func (e *Engine) RegisterLookups(lookups ...Lookup) error {
	for _, lookup := range lookups {
		if err := e.RegisterLookup(lookup); err != nil {
			return err
		}
	}
	return nil
}

// Rules returns the number of rules.
func (e *Engine) Rules() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.rules)
}

// Lookups returns the number of registered lookups.
func (e *Engine) Lookups() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.lookups)
}

// Evaluate runs all rules against the provided context.
func (e *Engine) Evaluate(ctx context.Context, evalCtx *EvalContext) (*Result, error) {
	if e.closed.Load() {
		return nil, ErrEngineClosed
	}

	if evalCtx == nil {
		return nil, ErrNilContext
	}

	// Apply timeout if configured
	if e.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.config.Timeout)
		defer cancel()
	}

	// Start trace
	ctx, endTrace := e.obs.Tracer.Start(ctx, "cortex.evaluate", "engine", e.name)
	startTime := time.Now()

	// Copy lookups to eval context
	e.mu.RLock()
	rules := e.rules
	for _, lookup := range e.lookups {
		evalCtx.RegisterLookup(lookup)
	}
	e.mu.RUnlock()

	var errors []RuleError

	for _, rule := range rules {
		// Check context cancellation
		select {
		case <-ctx.Done():
			endTrace(ctx.Err())
			e.obs.Metrics.Inc("cortex.evaluation.timeout", "engine", e.name)
			return nil, fmt.Errorf("%w: %v", ErrTimeout, ctx.Err())
		default:
		}

		// Check if halted
		if e.config.ShortCircuit && evalCtx.IsHalted() {
			break
		}

		// Evaluate rule
		err := e.evaluateRule(ctx, rule, evalCtx)
		if err != nil {
			evalCtx.incErrors()

			var ruleErr *RuleError
			if re, ok := err.(*RuleError); ok {
				ruleErr = re
			} else {
				ruleErr = NewRuleError(rule.ID(), "", "evaluate", err)
			}

			errors = append(errors, *ruleErr)
			e.obs.Logger.Error("rule evaluation failed", err, "rule_id", rule.ID())
			e.obs.Metrics.Inc("cortex.rules.failed", "engine", e.name, "rule_id", rule.ID())

			switch e.config.Mode {
			case ModeFailFast:
				endTrace(err)
				return newResult(evalCtx, errors), err
			case ModeCollectAll, ModeContinueOnError:
				continue
			}
		}

		evalCtx.incRulesEvaluated()
	}

	duration := time.Since(startTime)
	endTrace(nil)

	// Emit metrics
	if e.config.EnableMetrics {
		e.obs.Metrics.Inc("cortex.evaluations", "engine", e.name)
		e.obs.Metrics.Add("cortex.rules.evaluated", float64(evalCtx.RulesEvaluated()), "engine", e.name)
		e.obs.Metrics.Histogram("cortex.evaluation.duration", duration.Seconds(), "engine", e.name)
	}

	result := newResult(evalCtx, errors)
	return result, nil
}

func (e *Engine) evaluateRule(ctx context.Context, rule Rule, evalCtx *EvalContext) error {
	ctx, endTrace := e.obs.Tracer.Start(ctx, "cortex.rule", "rule_id", rule.ID())
	startTime := time.Now()

	err := rule.Evaluate(ctx, evalCtx)

	duration := time.Since(startTime)
	endTrace(err)

	if e.config.EnableMetrics {
		e.obs.Metrics.Histogram("cortex.rule.duration", duration.Seconds(), "rule_id", rule.ID())
	}

	return err
}

// Close closes the engine.
func (e *Engine) Close() error {
	e.closed.Store(true)
	return nil
}

// Clone creates a copy of the engine with the same configuration and lookups,
// but without any rules.
func (e *Engine) Clone(name string) *Engine {
	e.mu.RLock()
	defer e.mu.RUnlock()

	clone := New(name, e.config)
	clone.obs = e.obs

	for k, v := range e.lookups {
		clone.lookups[k] = v
	}

	return clone
}
