package main

type ConfigParser struct {
}

func NewConfigParser() *ConfigParser {
	return &ConfigParser{}
}

func (c *ConfigParser) getTasks() []Task {
	return []Task{
		{
			Name:        "help",
			Description: "print this output",
			Command:     "help",
		},
		{
			Name:        "init",
			Description: "initialize a new config in current working directory",
			Command:     "init",
		},
		{
			Name:        "version",
			Description: "print version information",
			Command:     "version",
		},
	}
}

func (c *ConfigParser) getDescription() string {
	return "dosh - shell-independent task manager"
}

func (c *ConfigParser) getEpilog() string {
	return "For more information, visit https://github.com/gkmngrgn/dosh"
}
