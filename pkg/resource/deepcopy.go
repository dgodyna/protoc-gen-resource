package resource

import (
	_ "embed"
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
	for _, f := range message.Fields {
		if f.Desc.Kind().IsValid() {
			if f.Desc.HasOptionalKeyword() {
				// bytes slice is not a pointer
				if f.Desc.Kind() == protoreflect.BytesKind {
					scalarTypes = append(scalarTypes, f.GoName)
					continue
				}
				optionScalarTypes = append(optionScalarTypes, f.GoName)
			} else {
				scalarTypes = append(scalarTypes, f.GoName)
			}
		}
	}

	sort.Slice(scalarTypes, func(i, j int) bool {
		return scalarTypes[i] > scalarTypes[j]
	})
	sort.Slice(optionScalarTypes, func(i, j int) bool {
		return optionScalarTypes[i] > optionScalarTypes[j]
	})

	sw.Do(deepCopyTmpl, generator.Args{
		"scalarTypes":       scalarTypes,
		"optionScalarTypes": optionScalarTypes,
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
