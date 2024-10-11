package config

import "flag"

type Flags struct {
	YAMLPath string
}

func ParseFlags() Flags {
	var yamlPath = flag.String("yaml-path", "configs/dev.yaml", "Path to the .yaml configuration file")

	flag.Parse()
	return Flags{
		YAMLPath: *yamlPath,
	}
}
