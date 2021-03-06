package resource

import (
	_ "embed"
	"fmt"
	"github.com/dgodyna/protoc-gen-resource/pkg/templates"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

//go:embed templates/deepcopy.gotmpl
var deepCopyTmpl string

//go:embed templates/deepcopy_object.gotmpl
var deepCopyObjectTmpl string

// deepCopyIntoMessage generates DeepCopyInto function. Fields will be processed in exactly same order as they present in message.
func (g *generator) deepCopyIntoMessage(message *protogen.Message) {
	g.sw.Do("// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.\n", nil)
	g.sw.Do("func (in *{{.GoIdent.GoName}}) DeepCopyInto(out *{{.GoIdent.GoName}}) {\n", message)
	for _, field := range message.Fields {
		g.doField(field)
	}
	g.sw.Do("return\n", nil)
	g.sw.Do("}\n\n", nil)
}

// doField process single message field.
func (g *generator) doField(field *protogen.Field) {

	// list just a flag on description - so check it first
	if field.Desc.IsList() {
		g.doList(field)
		return
	}

	switch field.Desc.Kind() {
	case protoreflect.BoolKind, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind, protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind, protoreflect.FloatKind, protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind,
		protoreflect.DoubleKind, protoreflect.StringKind:
		g.doScalar(field)
	case protoreflect.EnumKind:
		g.doEnum(field)
	case protoreflect.BytesKind:
		g.doBytes(field)
	case protoreflect.MessageKind:
		g.doMessage(field)
	default:
		panic(fmt.Errorf("kind '%s' not supported yet", field.Desc.Kind()))
	}

}

func (g *generator) doList(field *protogen.Field) {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind, protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind, protoreflect.FloatKind, protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind,
		protoreflect.DoubleKind, protoreflect.StringKind, protoreflect.BytesKind:
		g.doScalarList(field)
	case protoreflect.EnumKind:
		g.doEnumList(field)
	case protoreflect.MessageKind:
		g.doMessageList(field)
	default:
		panic(fmt.Errorf("kind '%s' not supported yet", field.Desc.Kind()))
	}
}

// doBytes process bytes fields. bytes are just scalar types, so simple assignment.
func (g *generator) doBytes(field *protogen.Field) {
	g.sw.Do("out.{{ .field.GoName }} = in.{{ .field.GoName }}\n", templates.Args{"field": field})
}

// doMessage process messages. Here we have following
// 1) if message is located in the same proto file - it's 1rst party message, and we can just call DeepCopyInto for it
// 2) else - we are not sure were runtime.Object functions implemented there - so cast to runtime.Object
// 2.1) If case successful - just call conversion func
// 2.2) If it not implements runtime.Object (e.g. 3rd party google messages) - generate DeepCopy for it right here.
func (g *generator) doMessage(field *protogen.Field) {

	// TODO process 1rst party messages
	// TODO process 3rd party messages
	g.sw.Do(`if in.{{ .field.GoName }} != nil {
        _, ok := interface{}(in.{{ .field.GoName }}).(runtime.Object)
        if ok {
          out.{{ .field.GoName }} = in.{{ .field.GoName }}.DeepCopy()
        } else {
          panic(fmt.Errorf("message field '{{.message.GoIdent.GoName}}{{ .field.GoName }}' does not implement runtime.Object"))
        }
	}`, templates.Args{"field": field, "message": field.Parent})
	g.sw.Do("\n", nil)
}

// doEnum process enums fields. Enums are just scalar types, so simple assignment.
func (g *generator) doEnum(field *protogen.Field) {
	g.sw.Do("out.{{ .field.GoName }} = in.{{ .field.GoName }}\n", templates.Args{"field": field})
}

// doScalar process scalars types. We support protobuf 3 optionals - also processing optional keyword.
func (g *generator) doScalar(field *protogen.Field) {
	if field.Desc.HasOptionalKeyword() {
		g.sw.Do("{{ .field.GoName }} := *in.{{ .field.GoName }}\n", templates.Args{"field": field})
		g.sw.Do("out.{{ .field.GoName }} = &{{ .field.GoName }}\n", templates.Args{"field": field})
	} else {
		g.sw.Do("out.{{ .field.GoName }} = in.{{ .field.GoName }}\n", templates.Args{"field": field})
	}
}

// doScalarList process repeatable scalars. Repeatable scalars are not optional, so no need to check it.
func (g *generator) doScalarList(field *protogen.Field) {
	g.sw.Do(`
if in.{{ .field.GoName }} != nil {
	in, out := &in.{{ .field.GoName }}, &out.{{ .field.GoName }}
	*out = make({{ .field | GoType }}, len(*in))
	copy(*out, *in)
}
`, templates.Args{"field": field})
}

func (g *generator) doMessageList(field *protogen.Field) {

	g.sw.Do(`
	inn, outt := &in.{{ .field.GoName }}, &out.{{ .field.GoName }}
	*outt = make([]*{{ .field.Message.GoIdent.GoName }}, len(*inn))
	for i := range *inn {
		if (*inn)[i] != nil {
			in, out := &(*inn)[i], &(*outt)[i]
			_, ok := interface{}(*in).(runtime.Object)
			if ok {
            	*out = new({{ .field.Message.GoIdent.GoName }})
				(*in).DeepCopyInto(*out)
			} else {
				panic(fmt.Errorf("message field '{{.message.GoIdent.GoName}}{{ .field.GoName }}' does not implement runtime.Object"))
			}
		}            
	}
`, templates.Args{"field": field, "message": field.Parent})
}

// doScalarList process repeatable scalars. Repeatable scalars are not optional, so no need to check it.
func (g *generator) doEnumList(field *protogen.Field) {
	g.sw.Do(`
if in.{{ .field.GoName }} != nil {
	in, out := &in.{{ .field.GoName }}, &out.{{ .field.GoName }}
	*out = make([]{{ .field.Enum.GoIdent.GoName }}, len(*in))
	copy(*out, *in)
}
`, templates.Args{"field": field})
}

// deepCopy generate DeepCopy function for provided message.
func (g *generator) deepCopy(message *protogen.Message) {
	g.sw.Do(deepCopyTmpl, templates.Args{
		"type": message.GoIdent.GoName,
	})
}

// deepCopyObject generate DeepCopyObject function for provided message.
func (g *generator) deepCopyObject(message *protogen.Message) {
	g.sw.Do(deepCopyObjectTmpl, templates.Args{
		"type": message.GoIdent.GoName,
	})
}
