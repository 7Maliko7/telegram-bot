package script

type Template struct {
	Text     string     `yaml:"text"`
	Actions  [][]string `yaml:"actions"`
	Triggers []string   `yaml:"triggers"`
}
