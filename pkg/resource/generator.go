package resource

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/dgodyna/protoc-gen-resource/pkg/templates"
	"go/format"
	"google.golang.org/protobuf/compiler/protogen"
	"sort"
)

//go:embed templates/package.gotmpl
var packageTmpl string

type generator struct {
	// firstPartyMessages messages which are present in file to generate.
	// for these messages we'll not generate copy inside other messages, just call
	// already generated deepcopy functions.
	firstPartyMessages map[*protogen.Message]interface{}

	// order in which messages will be generated.
	order []*protogen.Message

	// sw contains generated file content.
	sw *templates.SnippetWriter

	// protoPackage holds protobuf package.
	protoPackage string

	// goPackage golds go package.
	goPackage string
}

func Generate(gen *protogen.Plugin, filePath string) error {

	file := gen.FilesByPath[filePath]

	genFile := gen.NewGeneratedFile(
		file.GeneratedFilenamePrefix+".deepcopy.pb.go",
		file.GoImportPath,
	)

	generator := newGenerator(file)

	// if no messages - skip generation
	if len(generator.order) == 0 {
		genFile.Skip()
		return nil
	}

	// generate package and imports
	err := generator.generate()
	if err != nil {
		return err
	}

	sources := []byte(fmt.Sprintf("%v", generator.sw.Out()))

	formattedSources, err := format.Source(sources)
	if err != nil {
		return fmt.Errorf("unable to format generated sources : %w", err)
	}

	_, err = genFile.Write(formattedSources)

	return err
}

// generate all the deepcopy file content.
func (g *generator) generate() error {
	// init package and imports
	g.sw.Do(packageTmpl, templates.Args{"package": g.goPackage})
	// process all the messages

	for _, m := range g.order {
		err := g.genMessage(m)
		if err != nil {
			return err
		}
	}
	return g.sw.Error()
}

// doMessage generate single message
func (g *generator) genMessage(m *protogen.Message) error {
	err := g.genGvk(m)
	if err != nil {
		return fmt.Errorf("unable to generate GVK for message '%s' : %w", m.GoIdent.GoName, err)
	}
	g.deepCopyIntoMessage(m)
	if g.sw.Error() != nil {
		return fmt.Errorf("unable to generate DeepCopyInto function for message '%s' : %w", m.GoIdent.GoName, err)
	}
	g.deepCopy(m)
	if g.sw.Error() != nil {
		return fmt.Errorf("unable to generate DeepCopy function for message '%s' : %w", m.GoIdent.GoName, err)
	}
	g.deepCopyObject(m)
	if g.sw.Error() != nil {
		return fmt.Errorf("unable to generate DeepCopyObject function for message '%s' : %w", m.GoIdent.GoName, err)
	}
	return nil
}

// newGenerator creates a new instance of generator from provided protogen file.
// It'll collect all messages from file and construct order.
func newGenerator(file *protogen.File) *generator {
	// collect all nested messages
	var messages []*protogen.Message
	for _, m := range file.Messages {
		collectMessages(m, &messages)
	}

	// sort to have exactly same order
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].GoIdent.GoName > messages[j].GoIdent.GoName
	})

	firstPartyMessages := make(map[*protogen.Message]interface{}, len(messages))
	for _, m := range messages {
		firstPartyMessages[m] = new(interface{})
	}

	return &generator{
		firstPartyMessages: firstPartyMessages,
		order:              messages,
		sw:                 templates.NewSnippetWriter(bytes.NewBuffer([]byte{}), "{{", "}}", map[string]interface{}{}),
		protoPackage:       *file.Proto.Package,
		goPackage:          string(file.GoPackageName),
	}
}

// collectMessages will recursively collect all proto messages
func collectMessages(m *protogen.Message, all *[]*protogen.Message) {
	for _, subM := range m.Messages {
		collectMessages(subM, all)
	}
	*all = append(*all, m)
}
