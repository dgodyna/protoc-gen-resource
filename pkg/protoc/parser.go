package protoc

import (
	"fmt"
	"os/exec"
	"strings"
)

// Parser parses proto source into descriptors.
type Parser struct {
	// ProtoPaths Specifies the directories in which to search for imports. Directories will be searched in specified order.
	// If not given, the current working directory is used.
	// If not found in any of the these directories, the DescriptorsSet descriptors will be checked for required proto file.
	// For each element -I<element> protoc argument will be added.
	ProtoPaths []string

	// DescriptorsSet specifies a list of FILES each containing a FileDescriptorSet (a protocol buffer defined in descriptor.proto).
	// The FileDescriptor for each of the PROTO_FILES provided will be loaded from these FileDescriptorSets.
	// If a FileDescriptor appears multiple times, the first occurrence will be used.
	// If is not empty --descriptor_set_in=<FILES> protoc argument will be added, where <FILES> is comma separated list of files.
	DescriptorsSet []string

	// IncludeImports also includes all dependencies of the input files in the set, so that the set is self-contained.
	// If set to true --include_imports protoc argument will be added.
	IncludeImports bool

	// IncludeSourceInfo will create descriptors that include information about the original location of each decl in the source file
	// as well as surrounding comments.
	// If set to true --include_source_info protoc argument will be added.
	IncludeSourceInfo bool //--include_source_info

	// Path to protoc compiler. If not set up - will be searched in PATH.
	Protoc string
}

func (p *Parser) CodeGenerationRequest(filesToGenerate ...string) {

}

// protocCommand create protoc command for provided files
func (p Parser) protocCommand(filesToGenerate ...string) (*exec.Cmd, error) {

	// lookup protoc in PATH if it was not specified
	if p.Protoc == "" {
		protocPath, err := exec.LookPath("protoc")
		if err != nil {
			return nil, fmt.Errorf("cannot find protoc executable: %w", err)
		}
		p.Protoc = protocPath
	}

	var args []string
	for _, i := range p.ProtoPaths {
		args = append(args, fmt.Sprintf("-I%s", i))
	}

	if len(p.DescriptorsSet) != 0 {
		args = append(args, fmt.Sprintf("--descriptor_set_in=%s", strings.Join(p.ProtoPaths, ",")))
	}

	if p.IncludeImports {
		args = append(args, "--include_imports")
	}

	if p.IncludeSourceInfo {
		args = append(args, "--include_source_info")
	}

	args = append(args, filesToGenerate...)

	c := exec.Command(p.Protoc, args...)
	return c, nil
}
