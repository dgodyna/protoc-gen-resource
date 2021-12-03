package protoc

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"io/ioutil"
	"os"
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
	IncludeSourceInfo bool

	// Path to protoc compiler. If not set up - will be searched in PATH.
	Protoc string

	// DescriptorsSetOut holds filepath where to store generated proto descriptor. descriptor_set_out
	DescriptorsSetOut string
}

// TODO test that
// CodeGenerationRequest will return pluginpb.CodeGeneratorRequest for provided files.
func (p Parser) CodeGenerationRequest(filesToGenerate ...string) (*pluginpb.CodeGeneratorRequest, error) {

	if p.DescriptorsSetOut != "" {
		if _, err := os.Stat(p.DescriptorsSetOut); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// create file
				f, err := os.Create(p.DescriptorsSetOut)
				if err != nil {
					return nil, fmt.Errorf("unable to create discriptors out file '%s' : %w", p.DescriptorsSetOut, err)
				}
				cerr := f.Close()
				if cerr != nil {
					return nil, fmt.Errorf("unable to close descriptors out file '%s' : %w", p.DescriptorsSetOut, cerr)
				}
			} else {
				return nil, fmt.Errorf("unable to stat descriptors out file '%s' : %w", p.DescriptorsSetOut, err)
			}
		}
	} else {
		f, err := ioutil.TempFile("", "*.pb")
		if err != nil {
			// store descriptor into the temp file
			return nil, fmt.Errorf("cannot create descriptors out temp file: %w", err)
		}
		p.DescriptorsSetOut = f.Name()
		err = f.Close()
		if err != nil {
			return nil, fmt.Errorf("error during close descriptors out temp file '%s' : %w", p.DescriptorsSetOut, err)
		}
	}

	cmd, err := p.protocCommand(filesToGenerate...)
	if err != nil {
		return nil, fmt.Errorf("unable to create prtoc command : %w", err)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cannot start protoc: %w (%s)", err, string(out))
	}

	// Let's read got inside a file
	data, err := ioutil.ReadFile(p.DescriptorsSetOut)
	if err != nil {
		return nil, fmt.Errorf("cannot read protoc output file '%s' : %w", p.DescriptorsSetOut, err)
	}

	// Conversion to *protogen.File
	res := &descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(data, res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal protoc output file: %w", err)
	}

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: filesToGenerate,
		// nil case we didn't specify any plugins here
		Parameter: nil,
		ProtoFile: res.File,
	}

	return req, nil
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
		args = append(args, fmt.Sprintf("--descriptor_set_in=%s", strings.Join(p.DescriptorsSet, ",")))
	}

	if p.IncludeImports {
		args = append(args, "--include_imports")
	}

	if p.IncludeSourceInfo {
		args = append(args, "--include_source_info")
	}

	if p.DescriptorsSetOut != "" {
		args = append(args, fmt.Sprintf("--descriptor_set_out=%s", p.DescriptorsSetOut))

	}

	args = append(args, filesToGenerate...)

	c := exec.Command(p.Protoc, args...)

	return c, nil
}
