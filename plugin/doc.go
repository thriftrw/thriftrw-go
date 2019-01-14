// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package plugin provides the API for writing ThriftRW plugins.
//
// Plugins are standalone programs with names in the format
// thriftrw-plugin-$name where $name is the name of the plugin.
//
// 	// thriftrw-plugin-myfancyplugin/main.go
// 	package main
//
// 	import "go.uber.org/thriftrw/plugin"
//
// 	func main() {
// 		plugin.Main(&plugin.Plugin{
// 			Name: "myfancyplugin",
// 			// ...
// 		})
// 	}
//
// Note that the name in the executable MUST match the name in the Plugin
// struct. Plugins need to be installed and available on the $PATH before they
// can be used. Additionaly, plugins talk to the ThriftRW process using stdout
// and stdin. This means that plugin implementations MUST NOT write to stdout
// or read from stdin. Plugins MAY write to stderr and receive input via
// command line arguments.
//
// To use a plugin, pass its name in with the -p/--plugin option.
//
// 	thriftrw --plugin=myfancyplugin foo.thrift
//
// Arguments may be sent to plugins by including them in the name. These are
// passed through to the plugin process as command line arguments. For
// example,
//
// 	thriftrw --plugin='myfancyplugin --useContext'
//
// Will pass `--useContext` to `thriftrw-plugin-myfancyplugin`.
package plugin
