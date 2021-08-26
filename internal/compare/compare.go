package compare

import (
	"fmt"

	"go.uber.org/multierr"
	"go.uber.org/thriftrw/compile"
)

const (
	addReqError = "adding a required field %s to %s is not backwards compatible"
	changOptToReqError = "changing an optional field %s in %s to required is not backwards compatible"
	deleteServiceError = "deleting service %s is not backwards compatible"
	removeMethodError = "removing method %s in service %s is not backwards compatible"
)

// StructSpecs compares two structs defined in a Thrift file.
func StructSpecs(from, to *compile.StructSpec) error {
	fields := map[int16]*compile.FieldSpec{}
	// Assume that these two should be compared.
	for _, f := range from.Fields {
		// Capture state of all fields here.
		fields[f.ID] = f
	}

	var errors []error
	for _, toField := range to.Fields {
		if fromField, ok := fields[toField.ID]; ok {
			fromRequired := fromField.Required
			toRequired := toField.Required
			if !fromRequired && toRequired {
				errors = append(errors, fmt.Errorf(changOptToReqError, toField.ThriftName(), to.ThriftName()))
			}
		} else {
			if toField.Required {
				errors = append(errors, fmt.Errorf(addReqError, toField.ThriftName(), to.ThriftName()))
			}
		}
	}

	return multierr.Combine(errors...)
}

func Services(toModule, fromModule *compile.Module) error {
	var errors []error
	for n, fromService := range fromModule.Services {
		toServ, ok := toModule.Services[n]
		if !ok {
			// Service was deleted, which is not backwards compatible.
			errors = append(errors, fmt.Errorf(deleteServiceError, n))
			continue
		}
		for f, _ := range fromService.Functions {
			if _, ok :=  toServ.Functions[f]; !ok {

				errors = append(errors, fmt.Errorf(removeMethodError, f, n))
			}
		}
	}

	return multierr.Combine(errors...)
}
