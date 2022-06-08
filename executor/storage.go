package executor

import "strings"

type RuleStorage struct {
	// Path->Operation->Rule
	rules map[string]map[Operation]*Rule
}

func NewRuleStorage() *RuleStorage {
	return &RuleStorage{rules: make(map[string]map[Operation]*Rule)}
}

func (s *RuleStorage) Add(path string, operation Operation, rule *Rule) {
	_, ok := s.rules[path]
	if !ok {
		s.rules[path] = map[Operation]*Rule{
			operation: rule,
		}
		return
	}

	s.rules[path][operation] = rule
}

// Contains will fetch Rule by Path and Operation, but in case of AnyOperation other records will be ignored.
// AnyOperation has the highest priority in the search
func (s *RuleStorage) Contains(path string, operation Operation) *Rule {
	operationRules, ok := s.rules[path]
	if !ok {
		for rulePath, opRules := range s.rules {
			if strings.Contains(path, rulePath) {
				operationRules = opRules
				break
			}
		}
	}

	if rule, ok := operationRules[AnyOperation]; ok {
		return rule
	}

	return operationRules[operation]
}
