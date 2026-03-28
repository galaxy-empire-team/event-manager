package config

type App struct {
	LogLevel  string `split_words:"true" default:"info"`
	LogFormat string `split_words:"true" default:"json"`
}
