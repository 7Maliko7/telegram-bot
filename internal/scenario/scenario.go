package scenario

import (
	"github.com/7Maliko7/telegram-bot/internal/script"
)

type Scenario struct {
	list      []script.Step
	templates map[string]script.Template
}

func New(sc *script.Scenario) *Scenario {
	return &Scenario{
		list:      sc.Steps,
		templates: sc.Templates,
	}
}

func (s *Scenario) Step(id int) *script.Step {
	return &s.list[id]
}

func (s *Scenario) StepCount() int {
	return len(s.list)
}
