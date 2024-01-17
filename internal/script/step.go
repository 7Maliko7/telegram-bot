package script

type Step struct {
	Template *string    `yaml:"template"`
	Text     string     `yaml:"text"`
	Actions  [][]string `yaml:"actions"`
	Triggers []string   `yaml:"triggers"`
}
