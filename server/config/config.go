package config

import "github.com/namsral/flag"

type Config struct {
	StaticDir         string
	RedirectPlainHTTP bool
	SandboxMode       string
	SandboxBin        string
	Port              uint
}

func ParseArgs() (config *Config) {
	config = new(Config)
	flag.UintVar(&config.Port, "port", 8080, "")
	flag.StringVar(&config.StaticDir, "static-dir", "client/dist", "")
	flag.BoolVar(&config.RedirectPlainHTTP, "redirect-plain-http", false, "")
	flag.StringVar(&config.SandboxMode, "sandbox", "", "")
	flag.StringVar(&config.SandboxBin, "sandbox-bin", "", "")
	flag.Parse()

	return
}
