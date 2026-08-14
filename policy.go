package cortex

import (
	"context"
	"fmt"
)

// PolicyDecision is the outcome of a matching policy rule.
type PolicyDecision string

const (
	PolicyDeny  PolicyDecision = "deny"
	PolicyAllow PolicyDecision = "allow"
	PolicyAsk   PolicyDecision = "ask"
)

// EvalContext keys read and written by PolicyRule.
const (
	PolicyKeyTool     = "tool"
	PolicyKeyPath     = "path"
	PolicyKeyTarget   = "target"
	PolicyKeyDecision = "decision"
)

// PolicyRule matches a tool name and optional path glob, then sets a decision.
type PolicyRule struct {
	baseRule
	tool     string
	pattern  string
	decision PolicyDecision
}

// PolicyConfig configures a policy rule.
type PolicyConfig struct {
	ID          string
	Name        string
	Description string
	Deps        []string

	// Tool is the tool name to match. "*" matches any tool.
	Tool string

	// Pattern is an optional glob matched against path (or target) in the
	// eval context. * matches one path segment; ** matches across segments.
	// Empty pattern matches any path.
	Pattern string

	// Decision is written to the context on match: deny, allow, or ask.
	// When several rules match, precedence is deny > allow > ask.
	Decision PolicyDecision
}

// NewPolicy creates a new policy rule.
func NewPolicy(cfg PolicyConfig) (*PolicyRule, error) {
	if cfg.ID == "" {
		return nil, fmt.Errorf("%w: policy rule requires ID", ErrInvalidRule)
	}
	if cfg.Tool == "" {
		return nil, fmt.Errorf("%w: policy rule %q requires tool", ErrInvalidRule, cfg.ID)
	}
	switch cfg.Decision {
	case PolicyDeny, PolicyAllow, PolicyAsk:
	default:
		return nil, fmt.Errorf("%w: policy rule %q decision must be deny, allow, or ask", ErrInvalidRule, cfg.ID)
	}

	return &PolicyRule{
		baseRule: baseRule{
			id:          cfg.ID,
			name:        cfg.Name,
			description: cfg.Description,
			deps:        cfg.Deps,
		},
		tool:     cfg.Tool,
		pattern:  cfg.Pattern,
		decision: cfg.Decision,
	}, nil
}

// MustPolicy creates a new policy rule, panicking on error.
func MustPolicy(cfg PolicyConfig) *PolicyRule {
	r, err := NewPolicy(cfg)
	if err != nil {
		panic(err)
	}
	return r
}

// Evaluate matches tool and path against the rule; on match it sets decision.
func (r *PolicyRule) Evaluate(_ context.Context, evalCtx *EvalContext) error {
	if !r.match(evalCtx) {
		return nil
	}
	cur := PolicyDecision(ctxString(evalCtx, PolicyKeyDecision))
	evalCtx.Set(PolicyKeyDecision, string(mergePolicyDecision(cur, r.decision)))
	return nil
}

func (r *PolicyRule) match(evalCtx *EvalContext) bool {
	tool := ctxString(evalCtx, PolicyKeyTool)
	if r.tool != "*" && r.tool != tool {
		return false
	}
	if r.pattern == "" {
		return true
	}
	path := ctxString(evalCtx, PolicyKeyPath)
	if path == "" {
		path = ctxString(evalCtx, PolicyKeyTarget)
	}
	return matchGlob(r.pattern, path)
}

// mergePolicyDecision implements deny > allow > ask. An existing deny is never
// replaced; deny always wins over allow or ask.
func mergePolicyDecision(current, incoming PolicyDecision) PolicyDecision {
	if current == PolicyDeny || incoming == PolicyDeny {
		return PolicyDeny
	}
	if current == PolicyAllow || incoming == PolicyAllow {
		return PolicyAllow
	}
	if incoming != "" {
		return incoming
	}
	return current
}

func ctxString(evalCtx *EvalContext, key string) string {
	s, err := evalCtx.GetString(key)
	if err != nil {
		return ""
	}
	return s
}

// Tool returns the tool matcher.
func (r *PolicyRule) Tool() string { return r.tool }

// Pattern returns the path glob, or empty if any path matches.
func (r *PolicyRule) Pattern() string { return r.pattern }

// Decision returns the decision written on match.
func (r *PolicyRule) Decision() PolicyDecision { return r.decision }
