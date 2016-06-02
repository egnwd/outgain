package config

import "github.com/namsral/flag"

type Config struct {
	StaticDir         string
	RedirectPlainHTTP bool
	SandboxMode       string
	RunnerBin         string
	Port              uint
	DefaultAI         string
	BotAI             string
}

func ParseArgs() (config *Config) {
	config = new(Config)
	flag.UintVar(&config.Port, "port", 8080, "")
	flag.StringVar(&config.StaticDir, "static-dir", "client/dist", "")
	flag.BoolVar(&config.RedirectPlainHTTP, "redirect-plain-http", false, "")
	flag.StringVar(&config.RunnerBin, "runner-bin", "./runner/target/release/runner", "")
	flag.StringVar(&config.DefaultAI, "default-ai", "./default_ai.rb", "")
	flag.StringVar(&config.BotAI, "bot-ai", "./bot_ai.rb", "")
	flag.StringVar(&config.SandboxMode, "sandbox", "", "")
	flag.Parse()

	return
}
