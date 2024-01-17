package script

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

const (
	defaultConfigPath = "./scenario.yml"
)

type Scenario struct {
	Templates map[string]Template `yaml:"templates"`
	Steps     []Step              `yaml:"steps,flow"`
}

func New(path string) (*Scenario, error) {
	cp := path
	if path == "" {
		cp = defaultConfigPath
	}

	bytes, err := os.ReadFile(cp)
	if err != nil {
		return nil, err
	}

	c := Scenario{}
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	err = c.setDefaults()
	if err != nil {
		return nil, err
	}

	err = c.validate()
	if err != nil {
		return nil, err
	}

	c.applyTemplate()

	return &c, nil
}

func (c *Scenario) setDefaults() error {
	return defaults.Set(c)
}

func (c *Scenario) validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scenario) applyTemplate() {
	for i, st := range s.Steps {
		if st.Template == nil {
			continue
		}
		s.Steps[i].Text = fmt.Sprintf(s.Templates[*st.Template].Text, st.Text)

		if s.Templates[*st.Template].Actions == nil {
			continue
		}
		s.Steps[i].Actions = append(s.Steps[i].Actions, s.Templates[*st.Template].Actions...)
	}
}
