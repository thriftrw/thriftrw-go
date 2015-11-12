package ast

import "fmt"

// Program represents the full syntax tree for a single .thrift file.
type Program struct {

	// Headers

	Includes   []*Include
	Namespaces []*Namespace

	// Definitions

	Constants []*Constant
	Typedefs  []*Typedef
	Enums     []*Enum
	Structs   []*Struct
	Services  []*Service
}

// TODO is AddHeader and AddDefinition even needed

// AddHeader adds a new Header to the AST for the Thrift file.
func (p *Program) AddHeader(header Header) {
	switch hdr := header.(type) {
	case *Include:
		p.Includes = append(p.Includes, hdr)
	case *Namespace:
		p.Namespaces = append(p.Namespaces, hdr)
	default:
		panic(fmt.Sprintf("parser: unknown header type %T: %v", hdr, hdr))
	}
}

// AddDefinition adds a new Definition to the AST for the Thrift file.
func (p *Program) AddDefinition(definition Definition) error {
	switch def := definition.(type) {
	case *Constant:
		p.Constants = append(p.Constants, def)
	case *Typedef:
		p.Typedefs = append(p.Typedefs, def)
	case *Enum:
		p.Enums = append(p.Enums, def)
	case *Struct:
		p.Structs = append(p.Structs, def)
	case *Service:
		p.Services = append(p.Services, def)
	default:
		return fmt.Errorf("parser: unknown definition type %T: %v", def, def)
	}
	return nil
}
