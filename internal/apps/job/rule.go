package job

import "github.com/dungnh3/guardrails-assignment/internal/model"

type ICheckedRule interface {
	Validate(line string) ([]model.Finding, error)
}
