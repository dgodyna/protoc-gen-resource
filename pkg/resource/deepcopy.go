package resource

import (
	_ "embed"
	"fmt"
	"github.com/dgodyna/protoc-gen-resource/pkg/templates"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"sort"
)

//go:embed templates/deepcopy_into.gotmpl
var deepCopyIntoTmpl string

//go:embed templates/deepcopy.gotmpl
var deepCopyTmpl string

//go:embed templates/deepcopy_object.gotmpl
var deepCopyObjectTmpl string

// deepCopyInto generate DeepCopyInto function for provided message.
func (g *generator) deepCopyInto(message *protogen.Message) {
	var scalarTypes []string
	var optionScalarTypes []string
	var enumTypes []string
	var messages []string
	for _, f := range message.Fields {
		switch f.Desc.Kind() {
		case protoreflect.EnumKind:
			enumTypes = append(enumTypes, f.GoName)
		case protoreflect.BytesKind:
			scalarTypes = append(scalarTypes, f.GoName)
		case protoreflect.MessageKind:
			messages = append(messages, f.GoName)
		case protoreflect.BoolKind, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind, protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind, protoreflect.FloatKind, protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind,
			protoreflect.DoubleKind, protoreflect.StringKind:
			if f.Desc.HasOptionalKeyword() {
				optionScalarTypes = append(optionScalarTypes, f.GoName)
			} else {
				scalarTypes = append(scalarTypes, f.GoName)
			}
		default:
			panic(fmt.Errorf("kind '%s' not supported yet", f.Desc.Kind()))
		}
	}

	sort.Slice(scalarTypes, func(i, j int) bool {
		return scalarTypes[i] > scalarTypes[j]
	})
	sort.Slice(optionScalarTypes, func(i, j int) bool {
		return optionScalarTypes[i] > optionScalarTypes[j]
	})
	sort.Slice(enumTypes, func(i, j int) bool {
		return enumTypes[i] > enumTypes[j]
	})
	sort.Slice(messages, func(i, j int) bool {
		return messages[i] > messages[j]
	})

	g.sw.Do(deepCopyIntoTmpl, templates.Args{
		"scalarTypes":       scalarTypes,
		"optionScalarTypes": optionScalarTypes,
		"enumTypes":         enumTypes,
		"messages":          messages,
		"type":              message.GoIdent.GoName,
	})
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
