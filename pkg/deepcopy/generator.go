package deepcopy

import "google.golang.org/protobuf/compiler/protogen"

func Generate(gen *protogen.Plugin, filePath string) error {

	file := gen.FilesByPath[filePath]

	genFile := gen.NewGeneratedFile(
		file.GeneratedFilenamePrefix+".pb.deepcopy.go",
		file.GoImportPath,
	)

	_, err := genFile.Write([]byte("package protos\n"))

	return err
	// genFile.Skip()
}
