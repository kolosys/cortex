package cortex

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// EvalContext holds the state during rule evaluation.
// It provides thread-safe access to values, buildups, and lookups.
type EvalContext struct {
	ID string

	mu       sync.RWMutex
	values   map[string]any
	buildups map[string]*Buildup
	lookups  map[string]Lookup
	metadata map[string]string

	halted   bool
	haltedBy string

	rulesEvaluated atomic.Int64
	errCount       atomic.Int64
	startTime      time.Time
}

// NewEvalContext creates a new evaluation context.
func NewEvalContext() *EvalContext {
	return &EvalContext{
		ID:        generateID(),
		values:    make(map[string]any),
		buildups:  make(map[string]*Buildup),
		lookups:   make(map[string]Lookup),
		metadata:  make(map[string]string),
		startTime: time.Now(),
	}
}

// Get retrieves a value from the context.
func (e *EvalContext) Get(key string) (any, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	v, ok := e.values[key]
	return v, ok
}

// GetTyped retrieves a typed value from the context.
func GetTyped[T any](e *EvalContext, key string) (T, bool) {
	v, ok := e.Get(key)
	if !ok {
		var zero T
		return zero, false
	}
	typed, ok := v.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return typed, true
}

// Set stores a value in the context.
func (e *EvalContext) Set(key string, value any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.values[key] = value
}

// SetTyped stores a typed value in the context.
func SetTyped[T any](e *EvalContext, key string, value T) {
	e.Set(key, value)
}

// Delete removes a value from the context.
func (e *EvalContext) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.values, key)
}

// Has checks if a key exists in the context.
func (e *EvalContext) Has(key string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, ok := e.values[key]
	return ok
}

// Keys returns all keys in the context.
func (e *EvalContext) Keys() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	keys := make([]string, 0, len(e.values))
	for k := range e.values {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a copy of all values in the context.
func (e *EvalContext) Values() map[string]any {
	e.mu.RLock()
	defer e.mu.RUnlock()
	cp := make(map[string]any, len(e.values))
	for k, v := range e.values {
		cp[k] = v
	}
	return cp
}

// GetFloat64 retrieves a float64 value from the context.
func (e *EvalContext) GetFloat64(key string) (float64, error) {
	v, ok := e.Get(key)
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrValueNotFound, key)
	}
	return toFloat64(v)
}

// GetInt retrieves an int value from the context.
func (e *EvalContext) GetInt(key string) (int, error) {
	v, ok := e.Get(key)
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrValueNotFound, key)
	}
	return toInt(v)
}

// GetString retrieves a string value from the context.
func (e *EvalContext) GetString(key string) (string, error) {
	v, ok := e.Get(key)
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrValueNotFound, key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%w: expected string, got %T", ErrTypeMismatch, v)
	}
	return s, nil
}

// GetBool retrieves a bool value from the context.
func (e *EvalContext) GetBool(key string) (bool, error) {
	v, ok := e.Get(key)
	if !ok {
		return false, fmt.Errorf("%w: %s", ErrValueNotFound, key)
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("%w: expected bool, got %T", ErrTypeMismatch, v)
	}
	return b, nil
}

// RegisterLookup registers a lookup table in the context.
func (e *EvalContext) RegisterLookup(lookup Lookup) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.lookups[lookup.Name()] = lookup
}

// Lookup performs a lookup in the specified table.
func (e *EvalContext) Lookup(tableName string, key any) (any, bool, error) {
	e.mu.RLock()
	lookup, ok := e.lookups[tableName]
	e.mu.RUnlock()
	if !ok {
		return nil, false, fmt.Errorf("%w: %s", ErrLookupNotFound, tableName)
	}
	v, found := lookup.Get(key)
	return v, found, nil
}

// GetBuildup returns a buildup accumulator for the given key.
func (e *EvalContext) GetBuildup(key string) (*Buildup, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	b, ok := e.buildups[key]
	return b, ok
}

// GetOrCreateBuildup returns an existing buildup or creates a new one.
func (e *EvalContext) GetOrCreateBuildup(key string, op BuildupOperation, initial float64) *Buildup {
	e.mu.Lock()
	defer e.mu.Unlock()
	if b, ok := e.buildups[key]; ok {
		return b
	}
	b := &Buildup{
		Name:      key,
		Operation: op,
		value:     initial,
	}
	e.buildups[key] = b
	return b
}

// SetMetadata sets a metadata key-value pair.
func (e *EvalContext) SetMetadata(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metadata[key] = value
}

// GetMetadata retrieves a metadata value.
func (e *EvalContext) GetMetadata(key string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	v, ok := e.metadata[key]
	return v, ok
}

// Halt stops evaluation with the given rule ID.
func (e *EvalContext) Halt(ruleID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.halted = true
	e.haltedBy = ruleID
}

// IsHalted returns true if evaluation has been halted.
func (e *EvalContext) IsHalted() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.halted
}

// HaltedBy returns the rule ID that halted evaluation.
func (e *EvalContext) HaltedBy() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.haltedBy
}

// incRulesEvaluated increments the rules evaluated counter.
func (e *EvalContext) incRulesEvaluated() {
	e.rulesEvaluated.Add(1)
}

// incErrors increments the error counter.
func (e *EvalContext) incErrors() {
	e.errCount.Add(1)
}

// RulesEvaluated returns the number of rules evaluated.
func (e *EvalContext) RulesEvaluated() int64 {
	return e.rulesEvaluated.Load()
}

// ErrorCount returns the number of errors encountered.
func (e *EvalContext) ErrorCount() int64 {
	return e.errCount.Load()
}

// Duration returns the time since context creation.
func (e *EvalContext) Duration() time.Duration {
	return time.Since(e.startTime)
}

// Clone creates a shallow copy of the context.
func (e *EvalContext) Clone() *EvalContext {
	e.mu.RLock()
	defer e.mu.RUnlock()

	clone := &EvalContext{
		ID:        generateID(),
		values:    make(map[string]any, len(e.values)),
		buildups:  make(map[string]*Buildup, len(e.buildups)),
		lookups:   e.lookups, // share lookups
		metadata:  make(map[string]string, len(e.metadata)),
		startTime: time.Now(),
	}

	for k, v := range e.values {
		clone.values[k] = v
	}
	for k, v := range e.metadata {
		clone.metadata[k] = v
	}

	return clone
}

// toFloat64 converts various numeric types to float64.
func toFloat64(v any) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case float32:
		return float64(n), nil
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case int16:
		return float64(n), nil
	case int8:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case uint32:
		return float64(n), nil
	case uint16:
		return float64(n), nil
	case uint8:
		return float64(n), nil
	default:
		return 0, fmt.Errorf("%w: expected numeric, got %T", ErrTypeMismatch, v)
	}
}

// toInt converts various numeric types to int.
func toInt(v any) (int, error) {
	switch n := v.(type) {
	case int:
		return n, nil
	case int64:
		return int(n), nil
	case int32:
		return int(n), nil
	case int16:
		return int(n), nil
	case int8:
		return int(n), nil
	case uint:
		return int(n), nil
	case uint64:
		return int(n), nil
	case uint32:
		return int(n), nil
	case uint16:
		return int(n), nil
	case uint8:
		return int(n), nil
	case float64:
		return int(n), nil
	case float32:
		return int(n), nil
	default:
		return 0, fmt.Errorf("%w: expected numeric, got %T", ErrTypeMismatch, v)
	}
}

var idCounter atomic.Uint64

func generateID() string {
	return fmt.Sprintf("eval-%d-%d", time.Now().UnixNano(), idCounter.Add(1))
}
