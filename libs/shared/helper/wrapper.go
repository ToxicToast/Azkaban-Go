// Package helper functions for wrapping primitive types in protobuf wrappers
package helper

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// WrapInt64 wraps a int pointer in a protobuf Int64Value.
func WrapInt64(v *int64) *wrapperspb.Int64Value {
	if v == nil {
		return nil
	}
	return wrapperspb.Int64(*v)
}
