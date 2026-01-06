// Package parse provides JSON parsing for cortex rule definitions.
package parse

import "encoding/json"

// RuleSet represents a collection of rules from config.
type RuleSet struct {
	Version string           `json:"version"`
	Name    string           `json:"name"`
	Lookups []LookupDef      `json:"lookups,omitempty"`
	Rules   []RuleDefinition `json:"rules"`
}

// LookupDef defines a lookup table in config.
type LookupDef struct {
	Name    string         `json:"name"`
	Type    string         `json:"type"` // "map" or "range"
	Entries []LookupEntry  `json:"entries,omitempty"`
	Items   map[string]any `json:"items,omitempty"` // for map type
}

// LookupEntry defines a single entry in a range lookup.
type LookupEntry struct {
	Min   float64  `json:"min"`
	Max   *float64 `json:"max"` // nil means +infinity
	Value any      `json:"value"`
}

// RuleDefinition is a config-driven rule.
type RuleDefinition struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"` // assignment, formula, allocation, lookup, buildup
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Deps        []string       `json:"deps,omitempty"`
	Disabled    bool           `json:"disabled,omitempty"`
	Config      map[string]any `json:"config"`
}

// AssignmentDef is the config structure for assignment rules.
type AssignmentDef struct {
	Target string `json:"target"`
	Value  any    `json:"value"`
}

// FormulaDef is the config structure for formula rules.
type FormulaDef struct {
	Target     string   `json:"target"`
	Expression string   `json:"expression,omitempty"`
	Function   string   `json:"function,omitempty"` // named registered function
	Inputs     []string `json:"inputs,omitempty"`
}

// LookupRuleDef is the config structure for lookup rules.
type LookupRuleDef struct {
	Table    string `json:"table"`
	Key      string `json:"key"`
	Target   string `json:"target"`
	Default  any    `json:"default,omitempty"`
	Required bool   `json:"required,omitempty"`
}

// AllocationDef is the config structure for allocation rules.
type AllocationDef struct {
	Source    string             `json:"source"`
	Strategy  string             `json:"strategy"`
	Targets   []AllocationTarget `json:"targets"`
	Remainder string             `json:"remainder,omitempty"`
	Precision int                `json:"precision,omitempty"`
}

// AllocationTarget defines an allocation destination.
type AllocationTarget struct {
	Key    string  `json:"key"`
	Amount float64 `json:"amount"`
}

// BuildupDef is the config structure for buildup rules.
type BuildupDef struct {
	Buildup   string  `json:"buildup"`
	Operation string  `json:"operation"`
	Source    string  `json:"source,omitempty"`
	Initial   float64 `json:"initial,omitempty"`
	Target    string  `json:"target,omitempty"`
}

// unmarshalConfig unmarshals a map into a struct.
func unmarshalConfig(cfg map[string]any, target any) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
