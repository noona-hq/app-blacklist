package logger

type Config struct {
	Level      string `default:"info"`
	Stacktrace bool   `default:"false"`
	Structured bool   `default:"false"`
}
