package reasoner

// Ruleset groups a versioned collection of declarative rules.
type Ruleset struct {
	Version string `yaml:"version"`
	Rules   []Rule `yaml:"rules"`
}

// Rule defines a single declarative rule evaluated by the Rule Engine.
type Rule struct {
	ID       string `yaml:"id"`
	Intent   string `yaml:"intent"`
	Priority int    `yaml:"priority"`

	When ConditionBlock `yaml:"when"`
	Then RuleAction     `yaml:"then"`

	Explanation string `yaml:"explanation"`
}

// ConditionBlock represents logical groupings of conditions.
// - All: every condition must match
// - Any: at least one condition must match
type ConditionBlock struct {
	All []Condition `yaml:"all,omitempty"`
	Any []Condition `yaml:"any,omitempty"`
}

// Condition represents a single signal comparison.
type Condition struct {
	Signal string  `yaml:"signal"`
	Op     string  `yaml:"op"`
	Value  float64 `yaml:"value"`
}

// RuleAction defines the effect of a rule when matched.
type RuleAction struct {
	Status          string  `yaml:"status"`
	ConfidenceBoost float64 `yaml:"confidence_boost"`
}
