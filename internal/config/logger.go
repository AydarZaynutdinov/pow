package config

type LoggerConfig struct {
	Level string `yaml:"Level" validate:"oneof=debug info warn error dpanic panic fatal"`
	Time  string `yaml:"Time" validate:"oneof=rfc3339nano rfc3339 iso8601 millis nanos"`
}
