package reflect

import goreflect "reflect"

// IsNilInterfaceOrPointer returns only true if the given value is a nil
// interface of a nil pointer.
func IsNilInterfaceOrPointer(v interface{}) bool {
	return v == nil ||
		(goreflect.ValueOf(v).Kind() == goreflect.Ptr && goreflect.ValueOf(v).IsNil())
}
