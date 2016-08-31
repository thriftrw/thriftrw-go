package main

import "go.uber.org/thriftrw/plugin"

// A valid plugin that does not do anything.

func main() {
	plugin.Main(&plugin.Plugin{Name: "empty"})
}
