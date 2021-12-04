package resource

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/dgodyna/protoc-gen-resource/pkg/generator"
	"go/format"
	"google.golang.org/protobuf/compiler/protogen"
	"sort"
)

//go:embed templates/package.gotmpl
var packageTmpl string

func Generate(gen *protogen.Plugin, filePath string) error {

	file := gen.FilesByPath[filePath]

	var messages []*protogen.Message
	for _, m := range file.Messages {
		collectMessages(m, &messages)
	}

	genFile := gen.NewGeneratedFile(
		file.GeneratedFilenamePrefix+".pb.deepcopy.go",
		file.GoImportPath,
	)

	if len(messages) == 0 {
		genFile.Skip()
		return nil
	}

	// sort them to have always same order
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].GoIdent.GoName > messages[j].GoIdent.GoName
	})

	sw := generator.NewSnippetWriter(bytes.NewBuffer([]byte{}), "{{", "}}", map[string]interface{}{})

	sw.Do(packageTmpl, generator.Args{"package": file.GoPackageName})
	for _, m := range messages {
		deepCopyMessage(m, sw)
		deepCopyInto(m, sw)
		deepCopyObject(m, sw)
	}

	if sw.Error() != nil {
		return sw.Error()
	}

	sources := []byte(fmt.Sprintf("%v", sw.Out()))

	formattedSources, err := format.Source(sources)
	if err != nil {
		return fmt.Errorf("unable to format generated sources : %w", err)
	}

	_, err = genFile.Write(formattedSources)

	return err
}

// collectMessages will recursively collect all proto messages
func collectMessages(m *protogen.Message, all *[]*protogen.Message) {
	for _, subM := range m.Messages {
		collectMessages(subM, all)
	}
	*all = append(*all, m)
}