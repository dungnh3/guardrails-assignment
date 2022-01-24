package rule

type IChecked interface {
	Detect(line string) (bool, error)
}

type Severity string

const (
	HighSeverity   Severity = "HIGH"
	LowSeverity    Severity = "LOW"
	MediumSeverity Severity = "MEDIUM"
)

type RuleType string

const (
	SASTRuleType RuleType = "sast"
)

type RuleId string

const (
	G101RuleId RuleId = "G101"
)

type (
	Rule struct {
		Type     RuleType `json:"type"`
		RuleId   RuleId   `json:"rule_id"`
		Metadata Metadata `json:"metadata"`
	}

	Metadata struct {
		Description string   `json:"description"`
		Severity    Severity `json:"severity"`
	}
)
