package rule

import (
	"regexp"
)

var PrefixMatchedG101 = []string{"public_key", "private_key"}

var G101 = Rule{
	Type:   SASTRuleType,
	RuleId: "G101",
	Metadata: Metadata{
		Description: "Look for hard coded credentials",
		Severity:    HighSeverity,
	},
	DetectFn: DetectG101,
}

func DetectG101(line string) (bool, error) {
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
