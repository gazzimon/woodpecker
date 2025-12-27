package reasoner

import "fmt"

// EvaluateRules returns all rules that match the given intent and signal snapshot.
func EvaluateRules(
	intentID string,
	signals map[string]float64,
	rules []Rule,
) ([]Rule, error) {

	var matched []Rule

	for _, rule := range rules {
		if rule.Intent != intentID {
			continue
		}

		if evaluateConditionBlock(rule.When, signals) {
			matched = append(matched, rule)
		}
	}

	return matched, nil
}

func evaluateConditionBlock(block ConditionBlock, signals map[string]float64) bool {
	// ALL conditions (AND)
	for _, cond := range block.All {
		if !evaluateCondition(cond, signals) {
			return false
		}
	}

	// ANY conditions (OR)
	if len(block.Any) > 0 {
		ok := false
		for _, cond := range block.Any {
			if evaluateCondition(cond, signals) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}

	return true
}

func evaluateCondition(cond Condition, signals map[string]float64) bool {
	value, ok := signals[cond.Signal]
	if !ok {
		return false
	}

	switch cond.Op {
	case "gte":
		return value >= cond.Value
	case "gt":
		return value > cond.Value
	case "lte":
		return value <= cond.Value
	case "lt":
		return value < cond.Value
	case "eq":
		return value == cond.Value
	default:
		panic(fmt.Sprintf("unsupported operator: %s", cond.Op))
	}
}
