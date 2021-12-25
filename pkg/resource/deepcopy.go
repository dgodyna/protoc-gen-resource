package resource

import (
	_ "embed"
	"fmt"
	"github.com/dgodyna/protoc-gen-resource/pkg/generator"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"sort"
)

//go:embed templates/deepcopy.gotmpl
var deepCopyTmpl string

//go:embed templates/deepcopy_into.gotmpl
var deepCopyIntoTmpl string

//go:embed templates/deepcopy_object.gotmpl
var deepCopyObjectTmpl string

func deepCopyMessage(message *protogen.Message, sw *generator.SnippetWriter) {
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

	sw.Do(deepCopyTmpl, generator.Args{
		"scalarTypes":       scalarTypes,
		"optionScalarTypes": optionScalarTypes,
		"enumTypes":         enumTypes,
		"messages":          messages,
		"type":              message.GoIdent.GoName,
	})
}

func deepCopyInto(message *protogen.Message, sw *generator.SnippetWriter) {
	sw.Do(deepCopyIntoTmpl, generator.Args{
		"type": message.GoIdent.GoName,
	})
}

func deepCopyObject(message *protogen.Message, sw *generator.SnippetWriter) {
	sw.Do(deepCopyObjectTmpl, generator.Args{
		"type": message.GoIdent.GoName,
	})
}
