package provide

import (
	"flag"
	"fmt"
)

var (
	goVersion = ""
	buildTime = ""
	gitHash   = ""

	ConfigPath = "./configs/config.toml"
	version    = false
)

func InitFlag() {
	flag.StringVar(&ConfigPath, "f", "./configs/config.toml", "config file default(./configs/config.toml)")
	flag.BoolVar(&version, "v", false, "show version")
	flag.Parse()

	if version {
		BuildInfo()
		return
	}
}

func BuildInfo() string {
	return fmt.Sprintf("goVersion: %s\nbuildTime: %s\ngitHash: %s\n", goVersion, buildTime, gitHash)
}
