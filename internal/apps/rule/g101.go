package rule

import (
	"regexp"
)

var PrefixMatchedG101 = []string{"public_key", "private_key"}

type G101Rule struct {
	*Rule
}

func G101() *G101Rule {
	return &G101Rule{
		&Rule{
			Type:   SASTRuleType,
			RuleId: G101RuleId,
			Metadata: Metadata{
				Description: "Look for hard coded credentials",
				Severity:    HighSeverity,
			},
		},
	}
}

func (g *G101Rule) Detect(line string) (bool, error) {
	for _, prefix := range PrefixMatchedG101 {
		matched, err := regexp.MatchString(prefix, line)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}
