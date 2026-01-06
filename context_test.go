package cortex_test

import (
	"sync"
	"testing"

	"github.com/kolosys/cortex"
)

func TestEvalContextBasic(t *testing.T) {
	ctx := cortex.NewEvalContext()

	if ctx.ID == "" {
		t.Error("expected non-empty ID")
	}

	ctx.Set("key", "value")
	v, ok := ctx.Get("key")
	if !ok || v != "value" {
		t.Errorf("expected value='value', got %v", v)
	}

	if !ctx.Has("key") {
		t.Error("expected Has('key') to be true")
	}

	ctx.Delete("key")
	if ctx.Has("key") {
		t.Error("expected Has('key') to be false after delete")
	}
}

func TestEvalContextTyped(t *testing.T) {
	ctx := cortex.NewEvalContext()

	cortex.SetTyped(ctx, "str", "hello")
	cortex.SetTyped(ctx, "num", 42.5)
	cortex.SetTyped(ctx, "bool", true)

	s, ok := cortex.GetTyped[string](ctx, "str")
	if !ok || s != "hello" {
		t.Errorf("expected 'hello', got %v", s)
	}

	n, ok := cortex.GetTyped[float64](ctx, "num")
	if !ok || n != 42.5 {
		t.Errorf("expected 42.5, got %v", n)
	}

	b, ok := cortex.GetTyped[bool](ctx, "bool")
	if !ok || !b {
		t.Errorf("expected true, got %v", b)
	}

	// Wrong type
	_, ok = cortex.GetTyped[int](ctx, "str")
	if ok {
		t.Error("expected false for wrong type")
	}
}

func TestEvalContextGetFloat64(t *testing.T) {
	ctx := cortex.NewEvalContext()

	tests := []struct {
		name     string
		value    any
		expected float64
	}{
		{"float64", 42.5, 42.5},
		{"float32", float32(42.5), 42.5},
		{"int", 42, 42.0},
		{"int64", int64(42), 42.0},
		{"int32", int32(42), 42.0},
		{"uint", uint(42), 42.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx.Set("value", tt.value)
			f, err := ctx.GetFloat64("value")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, f)
			}
		})
	}
}

func TestEvalContextGetString(t *testing.T) {
	ctx := cortex.NewEvalContext()

	ctx.Set("str", "hello")
	s, err := ctx.GetString("str")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != "hello" {
		t.Errorf("expected 'hello', got %q", s)
	}

	// Not found
	_, err = ctx.GetString("missing")
	if err == nil {
		t.Error("expected error for missing key")
	}

	// Wrong type
	ctx.Set("num", 42)
	_, err = ctx.GetString("num")
	if err == nil {
		t.Error("expected error for wrong type")
	}
}

func TestEvalContextGetBool(t *testing.T) {
	ctx := cortex.NewEvalContext()

	ctx.Set("flag", true)
	b, err := ctx.GetBool("flag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !b {
		t.Error("expected true")
	}

	ctx.Set("flag2", false)
	b, err = ctx.GetBool("flag2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b {
		t.Error("expected false")
	}
}

func TestEvalContextKeys(t *testing.T) {
	ctx := cortex.NewEvalContext()

	ctx.Set("a", 1)
	ctx.Set("b", 2)
	ctx.Set("c", 3)

	keys := ctx.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}

	keySet := make(map[string]bool)
	for _, k := range keys {
		keySet[k] = true
	}

	for _, expected := range []string{"a", "b", "c"} {
		if !keySet[expected] {
			t.Errorf("expected key %q", expected)
		}
	}
}

func TestEvalContextValues(t *testing.T) {
	ctx := cortex.NewEvalContext()

	ctx.Set("x", 1)
	ctx.Set("y", 2)

	values := ctx.Values()
	if len(values) != 2 {
		t.Errorf("expected 2 values, got %d", len(values))
	}

	// Modify the returned map shouldn't affect original
	values["z"] = 3
	if ctx.Has("z") {
		t.Error("modifying returned map should not affect context")
	}
}

func TestEvalContextHalt(t *testing.T) {
	ctx := cortex.NewEvalContext()

	if ctx.IsHalted() {
		t.Error("expected not halted initially")
	}

	ctx.Halt("test-rule")

	if !ctx.IsHalted() {
		t.Error("expected halted after Halt()")
	}
	if ctx.HaltedBy() != "test-rule" {
		t.Errorf("expected haltedBy='test-rule', got %q", ctx.HaltedBy())
	}
}

func TestEvalContextMetadata(t *testing.T) {
	ctx := cortex.NewEvalContext()

	ctx.SetMetadata("key", "value")
	v, ok := ctx.GetMetadata("key")
	if !ok || v != "value" {
		t.Errorf("expected 'value', got %q", v)
	}

	_, ok = ctx.GetMetadata("missing")
	if ok {
		t.Error("expected false for missing metadata")
	}
}

func TestEvalContextClone(t *testing.T) {
	ctx := cortex.NewEvalContext()
	ctx.Set("x", 42)
	ctx.SetMetadata("env", "test")

	clone := ctx.Clone()

	// Clone should have same values
	x, _ := clone.GetFloat64("x")
	if x != 42 {
		t.Errorf("expected x=42, got %f", x)
	}

	// Clone should have different ID
	if clone.ID == ctx.ID {
		t.Error("clone should have different ID")
	}

	// Modifying clone shouldn't affect original
	clone.Set("y", 100)
	if ctx.Has("y") {
		t.Error("modifying clone should not affect original")
	}
}

func TestEvalContextConcurrency(t *testing.T) {
	ctx := cortex.NewEvalContext()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx.Set("key", i)
		}(i)
	}

	// Concurrent reads
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx.Get("key")
		}()
	}

	wg.Wait()
}
