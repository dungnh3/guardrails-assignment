package rule

type DetectFn func(line string) (bool, error)

type Severity string

const (
	HighSeverity   Severity = "HIGH"
	LowSeverity    Severity = "LOW"
	MediumSeverity Severity = "MEDIUM"
)

type Type string

const (
	SASTRuleType Type = "sast"
)

type (
	Rule struct {
		Type     Type     `json:"type"`
		RuleId   string   `json:"rule_id"`
		Metadata Metadata `json:"metadata"`
		DetectFn DetectFn
	}

	Metadata struct {
		Description string   `json:"description"`
		Severity    Severity `json:"severity"`
	}
)
