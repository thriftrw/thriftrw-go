// Copyright (c) 2016 Uber Technologies, Inc.
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

package gen

import (
	"bytes"
	"fmt"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/thriftrw/thriftrw-go/compile"
)

// TODO: This should be moved into a separate program or template

// YARPC generates YARPC-specific code for a service.
//
// Each service gets a ${serviceName}server and ${serviceName}client package.
func YARPC(i thriftPackageImporter, s *compile.ServiceSpec) (map[string]*bytes.Buffer, error) {
	files := make(map[string]*bytes.Buffer)

	thriftPackage, err := i.Package(s.ThriftFile())
	if err != nil {
		return nil, err
	}

	serverPkgName := strings.ToLower(s.Name) + "server"
	serverImportPath := filepath.Join(thriftPackage, "yarpc", serverPkgName)
	serverFile, err := yarpcGenerator{
		NewGenerator(i, serverImportPath, serverPkgName), i,
	}.server(s)
	if err != nil {
		return nil, err
	}
	files[filepath.Join(serverPkgName, "server.go")] = serverFile

	clientPkgName := strings.ToLower(s.Name) + "client"
	clientImportPath := filepath.Join(thriftPackage, "yarpc", clientPkgName)
	clientFile, err := yarpcGenerator{
		NewGenerator(i, clientImportPath, clientPkgName), i,
	}.client(s)
	if err != nil {
		return nil, err
	}
	files[filepath.Join(clientPkgName, "client.go")] = clientFile

	return files, nil
}

type yarpcGenerator struct {
	g Generator
	i thriftPackageImporter
}

func (yg yarpcGenerator) server(s *compile.ServiceSpec) (*bytes.Buffer, error) {
	if err := yg.iface(s, true); err != nil {
		return nil, err
	}

	err := yg.g.DeclareFromTemplate(
		`
		<$thrift := import "github.com/yarpc/yarpc-go/encoding/thrift">
		<$protocol := import "github.com/thriftrw/thriftrw-go/protocol">
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">

		func New(impl Interface) <$thrift>.Service {
			return service{handler{impl}}
		}

		type service struct{h handler}

		func (service) Name() string {
			 return "<.Name>"
		 }

		func (service) Protocol() <$protocol>.Protocol {
			return <$protocol>.Binary
		}

		func (s service) Handlers() map[string]<$thrift>.Handler {
			return map[string]<$thrift>.Handler{
				<range .Functions>
					"<.Name>": <$thrift>.HandlerFunc(s.h.<goCase .Name>),
				<end>
			}
		}

		type handler struct{impl Interface}

		<$service := .>
		<range .Functions>
			<$servicePackage := servicePackage $service>
			<$Args := printf "%s.%sArgs" $servicePackage (goCase .Name)>
			<$Helper := printf "%s.%sHelper" $servicePackage (goCase .Name)>

			<$vars := newNamespace>
			<$req := $vars.NewName "req">
			<$body := $vars.NewName "body">

			func (h handler) <goCase .Name>(
				<$req> *<$thrift>.Request,
				<$body> <$wire>.Value,
			) (<$wire>.Value, *<$thrift>.Response, error) {

				<$args := $vars.NewName "args">
				var <$args> <$Args>
				if err := <$args>.FromWire(<$body>); err != nil {
					return <$wire>.Value{}, nil, err
				}

				<$res := $vars.NewName "res">
				<$succ := $vars.NewName "success">
				<if .ResultSpec.ReturnType>
					<$succ>,
				<end>
				<$res>, err := h.impl.<goCase .Name>(
					<$req>,
					<range .ArgsSpec><$args>.<goCase .Name>,<end>
				)

				<$result := $vars.NewName "result">
				<$result>, err := <$Helper>.WrapResponse(
					<if .ResultSpec.ReturnType>
						<$succ>,
					<end>
					err)
				<$w := $vars.NewName "w">
				var <$w> <$wire>.Value
				if err == nil {
					<$w>, err = <$result>.ToWire()
				}
				return <$w>, <$res>, err
			}
		<end>
		`,
		s)

	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	if err := yg.g.Write(buff, token.NewFileSet()); err != nil {
		return nil, fmt.Errorf(
			"failed to write YARPC server for service %q: %v", s.Name, err)
	}

	return buff, nil
}

func (yg yarpcGenerator) client(s *compile.ServiceSpec) (*bytes.Buffer, error) {
	if err := yg.iface(s, false); err != nil {
		return nil, err
	}

	err := yg.g.DeclareFromTemplate(
		`
		<$transport := import "github.com/yarpc/yarpc-go/transport">
		<$thrift := import "github.com/yarpc/yarpc-go/encoding/thrift">
		<$protocol := import "github.com/thriftrw/thriftrw-go/protocol">

		func New(c <$transport>.Channel) Interface {
			return client{
				c: <$thrift>.New(<$thrift>.Config{
					Service: "<.Name>",
					Channel: c,
					Protocol: <$protocol>.Binary,
				}),
			}
		}

		type client struct{c <$thrift>.Client}

		<$service := .>
		<range .Functions>
			<$wire := import "github.com/thriftrw/thriftrw-go/wire">
			<$servicePackage := servicePackage $service>

			<$Result := printf "%s.%sResult" $servicePackage (goCase .Name)>
			<$Helper := printf "%s.%sHelper" $servicePackage (goCase .Name)>

			<$vars := newNamespace>
			func (c client) <goCase .Name>(
				<$vars.NewName "req"> *<$thrift>.Request,
				<range .ArgsSpec>
					<if .Required>
						<$vars.NewName .Name> <typeReference .Type>,
					<else>
						<$vars.NewName .Name> <typeReferencePtr .Type>,
					<end>
				<end>
			) (
				<if .ResultSpec.ReturnType>
					<$vars.NewName "success"> <typeReference .ResultSpec.ReturnType>,
				<end>
				<$vars.NewName "res"> *<$thrift>.Response,
				err error,
			 ) {
				<$req := $vars.Rotate "req">
				<$args := $vars.NewName "args">
				<$args> := <$servicePackage>.<goCase .Name>Helper.Args(
					<range .ArgsSpec>
						<$vars.Rotate .Name>,
					<end>
				)

				<$res := $vars.Rotate "res">
				<$body := $vars.NewName "body">

				<$w := $vars.NewName "w">
				var <$w> <$wire>.Value
				<$w>, err = <$args>.ToWire()
				if err != nil {
					return
				}

				var <$body> <$wire>.Value
				<$body>, <$res>, err = c.c.Call("<.Name>", <$req>, <$w>)
				if err != nil {
					return
				}

				<$result := $vars.NewName "result">
				var <$result> <$Result>
				if err = <$result>.FromWire(<$body>); err != nil {
					return
				}

				<if .ResultSpec.ReturnType>
					<$succ := $vars.Rotate "success">
					<$succ>, err = <$Helper>.UnwrapResponse(&<$result>)
				<else>
					err = <$Helper>.UnwrapResponse(&<$result>)
				<end>
				return
			}
		<end>
		`, s)

	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	if err := yg.g.Write(buff, token.NewFileSet()); err != nil {
		return nil, fmt.Errorf(
			"failed to write YARPC server for service %q: %v", s.Name, err)
	}

	return buff, nil
}

// iface generates the service interface for the service and the client.
func (yg yarpcGenerator) iface(s *compile.ServiceSpec, isServer bool) error {
	return yg.g.DeclareFromTemplate(
		`
		<$thrift := import "github.com/yarpc/yarpc-go/encoding/thrift">

		type Interface interface {
			<if .Parent>
				<if isServer>
					<import (serverPackage .Parent)>.Interface
				<else>
					<import (clientPackage .Parent)>.Interface
				<end>
			<end>

			<range .Functions>
				<$params := newNamespace>
				<goCase .Name>(
					<$params.NewName "req"> *<$thrift>.Request,
					<range .ArgsSpec>
						<if .Required>
							<$params.NewName .Name> <typeReference .Type>,
						<else>
							<$params.NewName .Name> <typeReferencePtr .Type>,
						<end>
					<end>
				) (
					<if .ResultSpec.ReturnType>
						<typeReference .ResultSpec.ReturnType>,
					<end>
					*<$thrift>.Response,
					 error,
				 )
			<end>
		}
		`,
		s,
		TemplateFunc("isServer", func() bool { return isServer }),
		TemplateFunc("clientPackage", func(service *compile.ServiceSpec) (string, error) {
			return yg.clientPackage(service)
		}),
		TemplateFunc("serverPackage", func(service *compile.ServiceSpec) (string, error) {
			return yg.serverPackage(service)
		}),
	)
}

func (yg yarpcGenerator) basePackage(s *compile.ServiceSpec) (string, error) {
	pkg, err := yg.i.Package(s.ThriftFile())
	if err != nil {
		return "", err
	}
	return filepath.Join(pkg, "yarpc"), nil
}

func (yg yarpcGenerator) serverPackage(s *compile.ServiceSpec) (string, error) {
	pkg, err := yg.basePackage(s)
	if err != nil {
		return "", err
	}
	return filepath.Join(pkg, strings.ToLower(s.Name)+"server"), nil
}

func (yg yarpcGenerator) clientPackage(s *compile.ServiceSpec) (string, error) {
	pkg, err := yg.basePackage(s)
	if err != nil {
		return "", err
	}
	return filepath.Join(pkg, strings.ToLower(s.Name)+"client"), nil
}
