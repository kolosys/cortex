package cortex_test

import (
	"context"
	"testing"

	"github.com/kolosys/cortex"
)

func TestPolicyEvaluate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     cortex.PolicyConfig
		tool    string
		path    string
		useTarget bool
		preset  string
		want    string
		wantSet bool
	}{
		{
			name:    "edit src/** matches src/a.go",
			cfg:     cortex.PolicyConfig{ID: "edit-src", Tool: "edit", Pattern: "src/**", Decision: cortex.PolicyAllow},
			tool:    "edit",
			path:    "src/a.go",
			want:    "allow",
			wantSet: true,
		},
		{
			name:    "edit src/** does not match docs/a.go",
			cfg:     cortex.PolicyConfig{ID: "edit-src", Tool: "edit", Pattern: "src/**", Decision: cortex.PolicyAllow},
			tool:    "edit",
			path:    "docs/a.go",
			wantSet: false,
		},
		{
			name:    "* tool matches any tool",
			cfg:     cortex.PolicyConfig{ID: "any", Tool: "*", Pattern: "src/**", Decision: cortex.PolicyDeny},
			tool:    "read",
			path:    "src/a.go",
			want:    "deny",
			wantSet: true,
		},
		{
			name:    "unmatched leaves decision unset",
			cfg:     cortex.PolicyConfig{ID: "edit-src", Tool: "edit", Pattern: "src/**", Decision: cortex.PolicyAsk},
			tool:    "write",
			path:    "src/a.go",
			wantSet: false,
		},
		{
			name:      "path read from target when path absent",
			cfg:       cortex.PolicyConfig{ID: "edit-src", Tool: "edit", Pattern: "src/**", Decision: cortex.PolicyAllow},
			tool:      "edit",
			path:      "src/a.go",
			useTarget: true,
			want:      "allow",
			wantSet:   true,
		},
		{
			name:    "empty pattern matches any path",
			cfg:     cortex.PolicyConfig{ID: "any-path", Tool: "edit", Decision: cortex.PolicyAsk},
			tool:    "edit",
			path:    "docs/secret.md",
			want:    "ask",
			wantSet: true,
		},
		{
			name:    "* matches one segment only",
			cfg:     cortex.PolicyConfig{ID: "one-seg", Tool: "edit", Pattern: "src/*", Decision: cortex.PolicyAllow},
			tool:    "edit",
			path:    "src/a.go",
			want:    "allow",
			wantSet: true,
		},
		{
			name:    "* does not cross segments",
			cfg:     cortex.PolicyConfig{ID: "one-seg", Tool: "edit", Pattern: "src/*", Decision: cortex.PolicyAllow},
			tool:    "edit",
			path:    "src/pkg/a.go",
			wantSet: false,
		},
		{
			name:    "** matches nested path",
			cfg:     cortex.PolicyConfig{ID: "nested", Tool: "edit", Pattern: "src/**", Decision: cortex.PolicyAllow},
			tool:    "edit",
			path:    "src/pkg/a.go",
			want:    "allow",
			wantSet: true,
		},
		{
			name:    "deny overwrites existing allow",
			cfg:     cortex.PolicyConfig{ID: "deny", Tool: "*", Decision: cortex.PolicyDeny},
			tool:    "edit",
			preset:  "allow",
			want:    "deny",
			wantSet: true,
		},
		{
			name:    "existing deny not overwritten by allow",
			cfg:     cortex.PolicyConfig{ID: "allow", Tool: "*", Decision: cortex.PolicyAllow},
			tool:    "edit",
			preset:  "deny",
			want:    "deny",
			wantSet: true,
		},
		{
			name:    "allow wins over existing ask",
			cfg:     cortex.PolicyConfig{ID: "allow", Tool: "*", Decision: cortex.PolicyAllow},
			tool:    "edit",
			preset:  "ask",
			want:    "allow",
			wantSet: true,
		},
		{
			name:    "unmatched preserves existing decision",
			cfg:     cortex.PolicyConfig{ID: "edit-src", Tool: "edit", Pattern: "src/**", Decision: cortex.PolicyDeny},
			tool:    "edit",
			path:    "docs/a.go",
			preset:  "allow",
			want:    "allow",
			wantSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := cortex.MustPolicy(tt.cfg)
			evalCtx := cortex.NewEvalContext()
			if tt.tool != "" {
				evalCtx.Set(cortex.PolicyKeyTool, tt.tool)
			}
			if tt.path != "" {
				if tt.useTarget {
					evalCtx.Set(cortex.PolicyKeyTarget, tt.path)
				} else {
					evalCtx.Set(cortex.PolicyKeyPath, tt.path)
				}
			}
			if tt.preset != "" {
				evalCtx.Set(cortex.PolicyKeyDecision, tt.preset)
			}

			if err := rule.Evaluate(context.Background(), evalCtx); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, ok := evalCtx.Get(cortex.PolicyKeyDecision)
			if tt.wantSet != ok {
				t.Fatalf("decision set=%v, want %v", ok, tt.wantSet)
			}
			if !tt.wantSet {
				return
			}
			if got != tt.want {
				t.Errorf("decision=%v, want %q", got, tt.want)
			}
		})
	}
}

func TestPolicyValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     cortex.PolicyConfig
		wantErr bool
	}{
		{"missing ID", cortex.PolicyConfig{Tool: "edit", Decision: cortex.PolicyAllow}, true},
		{"missing tool", cortex.PolicyConfig{ID: "id", Decision: cortex.PolicyAllow}, true},
		{"missing decision", cortex.PolicyConfig{ID: "id", Tool: "edit"}, true},
		{"bad decision", cortex.PolicyConfig{ID: "id", Tool: "edit", Decision: "maybe"}, true},
		{"valid star tool", cortex.PolicyConfig{ID: "id", Tool: "*", Decision: cortex.PolicyDeny}, false},
		{"valid with pattern", cortex.PolicyConfig{ID: "id", Tool: "edit", Pattern: "src/**", Decision: cortex.PolicyAsk}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cortex.NewPolicy(tt.cfg)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPolicyMustPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	cortex.MustPolicy(cortex.PolicyConfig{})
}

func TestPolicyAccessors(t *testing.T) {
	rule := cortex.MustPolicy(cortex.PolicyConfig{
		ID:          "edit-src",
		Name:        "Edit src",
		Description: "Allow edits under src",
		Tool:        "edit",
		Pattern:     "src/**",
		Decision:    cortex.PolicyAllow,
		Deps:        []string{"dep1"},
	})
	if rule.ID() != "edit-src" {
		t.Errorf("ID=%q", rule.ID())
	}
	if rule.Name() != "Edit src" {
		t.Errorf("Name=%q", rule.Name())
	}
	if rule.Tool() != "edit" {
		t.Errorf("Tool=%q", rule.Tool())
	}
	if rule.Pattern() != "src/**" {
		t.Errorf("Pattern=%q", rule.Pattern())
	}
	if rule.Decision() != cortex.PolicyAllow {
		t.Errorf("Decision=%q", rule.Decision())
	}
	if len(rule.Dependencies()) != 1 {
		t.Errorf("Deps=%v", rule.Dependencies())
	}
}

func TestPolicyEngineDenyWins(t *testing.T) {
	engine := cortex.New("gate", cortex.DefaultConfig())
	if err := engine.AddRules(
		cortex.MustPolicy(cortex.PolicyConfig{ID: "allow-src", Tool: "edit", Pattern: "src/**", Decision: cortex.PolicyAllow}),
		cortex.MustPolicy(cortex.PolicyConfig{ID: "deny-src", Tool: "edit", Pattern: "src/**", Decision: cortex.PolicyDeny}),
	); err != nil {
		t.Fatal(err)
	}

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set(cortex.PolicyKeyTool, "edit")
	evalCtx.Set(cortex.PolicyKeyPath, "src/a.go")

	if _, err := engine.Evaluate(context.Background(), evalCtx); err != nil {
		t.Fatal(err)
	}
	got, _ := evalCtx.GetString(cortex.PolicyKeyDecision)
	if got != string(cortex.PolicyDeny) {
		t.Errorf("decision=%q, want deny", got)
	}
}
