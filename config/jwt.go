package config

// JWT struct
type JWT struct {
	RSecret string `yaml:"jwt.rSecret"`
	Secret  string `yaml:"jwt.secret"`
}
