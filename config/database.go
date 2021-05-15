package config

// database struct
type database struct {
	Username    string `yaml:"postgress.username"`
	Password    string `yaml:"postgress.password"`
	Host        string `yaml:"postgress.host"`
	Schema      string `yaml:"postgress.schema"`
	Automigrate bool   `yaml:"postgress.automigrate"`
	Logger      bool   `yaml:"postgress.logger"`
	Namespace   string
}
