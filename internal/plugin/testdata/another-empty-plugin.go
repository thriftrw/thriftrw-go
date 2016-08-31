package main

import "github.com/thriftrw/thriftrw-go/plugin"

// A valid plugin that does not do anything.

func main() {
	plugin.Main(&plugin.Plugin{Name: "another-empty"})
}
