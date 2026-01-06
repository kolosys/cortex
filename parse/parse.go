package parse

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/kolosys/cortex"
)

// Parser parses config into rules.
type Parser struct {
	formulas map[string]cortex.FormulaFunc
}

// NewParser creates a new parser.
func NewParser() *Parser {
	return &Parser{
		formulas: make(map[string]cortex.FormulaFunc),
	}
}

// RegisterFormula registers a named formula function for use in config.
func (p *Parser) RegisterFormula(name string, fn cortex.FormulaFunc) {
	p.formulas[name] = fn
}

// ParseJSON parses a rule set from JSON.
func (p *Parser) ParseJSON(data []byte) (*RuleSet, error) {
	var rs RuleSet
	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}
	return &rs, nil
}

// ToLookups converts lookup definitions to Lookup instances.
func (p *Parser) ToLookups(rs *RuleSet) ([]cortex.Lookup, error) {
	lookups := make([]cortex.Lookup, 0, len(rs.Lookups))

	for _, def := range rs.Lookups {
		lookup, err := p.buildLookup(def)
		if err != nil {
			return nil, fmt.Errorf("lookup %q: %w", def.Name, err)
		}
		lookups = append(lookups, lookup)
	}

	return lookups, nil
}

func (p *Parser) buildLookup(def LookupDef) (cortex.Lookup, error) {
	switch def.Type {
	case "map":
		if len(def.Items) == 0 {
			return nil, fmt.Errorf("map lookup requires items")
		}
		return cortex.NewMapLookup(def.Name, def.Items), nil

	case "range":
		if len(def.Entries) == 0 {
			return nil, fmt.Errorf("range lookup requires entries")
		}
		ranges := make([]cortex.RangeEntry[float64], len(def.Entries))
		for i, e := range def.Entries {
			max := math.Inf(1)
			if e.Max != nil {
				max = *e.Max
			}
			val, err := toFloat64(e.Value)
			if err != nil {
				return nil, fmt.Errorf("entry %d: %w", i, err)
			}
			ranges[i] = cortex.RangeEntry[float64]{
				Min:   e.Min,
				Max:   max,
				Value: val,
			}
		}
		return cortex.NewRangeLookup(def.Name, ranges), nil

	default:
		return nil, fmt.Errorf("unknown lookup type: %s", def.Type)
	}
}

// ToRules converts rule definitions to Rule instances.
func (p *Parser) ToRules(rs *RuleSet) ([]cortex.Rule, error) {
	rules := make([]cortex.Rule, 0, len(rs.Rules))

	for _, def := range rs.Rules {
		if def.Disabled {
			continue
		}

		rule, err := p.buildRule(def)
		if err != nil {
			return nil, fmt.Errorf("rule %q: %w", def.ID, err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (p *Parser) buildRule(def RuleDefinition) (cortex.Rule, error) {
	switch def.Type {
	case "assignment":
		return p.buildAssignment(def)
	case "formula":
		return p.buildFormula(def)
	case "lookup":
		return p.buildLookupRule(def)
	case "allocation":
		return p.buildAllocation(def)
	case "buildup":
		return p.buildBuildup(def)
	default:
		return nil, fmt.Errorf("unknown rule type: %s", def.Type)
	}
}

func (p *Parser) buildAssignment(def RuleDefinition) (*cortex.AssignmentRule, error) {
	var cfg AssignmentDef
	if err := unmarshalConfig(def.Config, &cfg); err != nil {
		return nil, err
	}

	return cortex.NewAssignment(cortex.AssignmentConfig{
		ID:          def.ID,
		Name:        def.Name,
		Description: def.Description,
		Deps:        def.Deps,
		Target:      cfg.Target,
		Value:       cfg.Value,
	})
}

func (p *Parser) buildFormula(def RuleDefinition) (*cortex.FormulaRule, error) {
	var cfg FormulaDef
	if err := unmarshalConfig(def.Config, &cfg); err != nil {
		return nil, err
	}

	config := cortex.FormulaConfig{
		ID:          def.ID,
		Name:        def.Name,
		Description: def.Description,
		Deps:        def.Deps,
		Target:      cfg.Target,
		Inputs:      cfg.Inputs,
		Expression:  cfg.Expression,
	}

	// Use registered function if specified
	if cfg.Function != "" {
		fn, ok := p.formulas[cfg.Function]
		if !ok {
			return nil, fmt.Errorf("unknown formula function: %s", cfg.Function)
		}
		config.Formula = fn
	}

	return cortex.NewFormula(config)
}

func (p *Parser) buildLookupRule(def RuleDefinition) (*cortex.LookupRule, error) {
	var cfg LookupRuleDef
	if err := unmarshalConfig(def.Config, &cfg); err != nil {
		return nil, err
	}

	return cortex.NewLookup(cortex.LookupConfig{
		ID:          def.ID,
		Name:        def.Name,
		Description: def.Description,
		Deps:        def.Deps,
		Table:       cfg.Table,
		Key:         cfg.Key,
		Target:      cfg.Target,
		Default:     cfg.Default,
		Required:    cfg.Required,
	})
}

func (p *Parser) buildAllocation(def RuleDefinition) (*cortex.AllocationRule, error) {
	var cfg AllocationDef
	if err := unmarshalConfig(def.Config, &cfg); err != nil {
		return nil, err
	}

	strategy, err := cortex.ParseAllocationStrategy(cfg.Strategy)
	if err != nil {
		return nil, err
	}

	targets := make([]cortex.AllocationTarget, len(cfg.Targets))
	for i, t := range cfg.Targets {
		targets[i] = cortex.AllocationTarget{
			Key:    t.Key,
			Amount: t.Amount,
		}
	}

	return cortex.NewAllocation(cortex.AllocationConfig{
		ID:          def.ID,
		Name:        def.Name,
		Description: def.Description,
		Deps:        def.Deps,
		Source:      cfg.Source,
		Strategy:    strategy,
		Targets:     targets,
		Remainder:   cfg.Remainder,
		Precision:   cfg.Precision,
	})
}

func (p *Parser) buildBuildup(def RuleDefinition) (*cortex.BuildupRule, error) {
	var cfg BuildupDef
	if err := unmarshalConfig(def.Config, &cfg); err != nil {
		return nil, err
	}

	op, err := cortex.ParseBuildupOperation(cfg.Operation)
	if err != nil {
		return nil, err
	}

	return cortex.NewBuildup(cortex.BuildupConfig{
		ID:          def.ID,
		Name:        def.Name,
		Description: def.Description,
		Deps:        def.Deps,
		Buildup:     cfg.Buildup,
		Operation:   op,
		Source:      cfg.Source,
		Initial:     cfg.Initial,
		Target:      cfg.Target,
	})
}

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
	case json.Number:
		return n.Float64()
	default:
		return 0, fmt.Errorf("expected number, got %T", v)
	}
}

// ParseAndBuild is a convenience function that parses JSON and builds an Engine.
func ParseAndBuild(name string, data []byte, config *cortex.Config) (*cortex.Engine, error) {
	parser := NewParser()
	return parser.ParseAndBuildEngine(name, data, config)
}

// ParseAndBuildEngine parses JSON and builds an Engine.
func (p *Parser) ParseAndBuildEngine(name string, data []byte, config *cortex.Config) (*cortex.Engine, error) {
	rs, err := p.ParseJSON(data)
	if err != nil {
		return nil, err
	}

	engine := cortex.New(name, config)

	lookups, err := p.ToLookups(rs)
	if err != nil {
		return nil, err
	}
	if err := engine.RegisterLookups(lookups...); err != nil {
		return nil, err
	}

	rules, err := p.ToRules(rs)
	if err != nil {
		return nil, err
	}
	if err := engine.AddRules(rules...); err != nil {
		return nil, err
	}

	return engine, nil
}
