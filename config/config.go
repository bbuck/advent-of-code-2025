package config

import "flag"

type Config struct {
	Solve bool
}

var solve = flag.Bool("solve", false, "Use the real problem input as input")

// Parses system flags and returns a cleanly wrapped config object
func Parse() Config {
	flag.Parse()

	return Config{
		Solve: *solve,
	}
}
