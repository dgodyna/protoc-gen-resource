package resource

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GoType returns corresponding go type of the field
func GoType(f *protogen.Field) string {

	if f.Desc.IsList() {
		return "[]" + getUnderlingTypeName(f)
	}
	return getUnderlingTypeName(f)
}

func getUnderlingTypeName(f *protogen.Field) string {

	switch f.Desc.Kind() {
	// scalars first
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "[]byte"
	// scalars end
	default:
		panic(fmt.Sprintf("type '%+v' is not supported for conversion to go type"))
	}
}
