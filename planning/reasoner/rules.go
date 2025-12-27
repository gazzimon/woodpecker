package reasoner

type Ruleset struct {
	Version string `yaml:"version"`
	Rules   []Rule `yaml:"rules"`
}

type Rule struct {
	ID     string `yaml:"id"`
	Intent string `yaml:"intent"`

	When ConditionBlock `yaml:"when"`
	Then RuleAction     `yaml:"then"`

	Explanation string `yaml:"explanation"`
}

type ConditionBlock struct {
	All []Condition `yaml:"all,omitempty"`
	Any []Condition `yaml:"any,omitempty"`
}

type Condition struct {
	Signal string  `yaml:"signal"`
	Op     string  `yaml:"op"`
	Value  float64 `yaml:"value"`
}

type RuleAction struct {
	Status           string  `yaml:"status"`
	ConfidenceBoost  float64 `yaml:"confidence_boost"`
}
