package resource

import (
	_ "embed"
	"github.com/dgodyna/protoc-gen-resource/pkg/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

//go:embed templates/deepcopy.gotmpl
var deepCopyTmpl string

//go:embed templates/deepcopy_into.gotmpl
var deepCopyIntoTmpl string

//go:embed templates/deepcopy_object.gotmpl
var deepCopyObjectTmpl string

func deepCopyMessage(message *protogen.Message, sw *generator.SnippetWriter) {
	var scalarTypes []string
	for _, f := range message.Fields {
		if f.Desc.Kind().IsValid() {
			scalarTypes = append(scalarTypes, f.GoName)
		}
	}

	sw.Do(deepCopyTmpl, generator.Args{
		"scalarTypes": scalarTypes,
		"type":        message.GoIdent.GoName,
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
