package resource

import (
	"github.com/dgodyna/protoc-gen-resource/pkg/protoc"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/pluginpb"
	"gotest.tools/assert"
	"os"
	"path/filepath"
	"testing"
)

// Note this test check exact generation syntax and will need update after each change in generation.
// It's required to don't break existing generation during refactoring and implementation.
// The test of generation behaviour located at //examples/test.
func TestGenerate(t *testing.T) {
	type args struct {
		descriptorPath string
		fileToGenerate string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantFilePath string
	}{
		{
			name: "Simple Types",
			args: args{
				descriptorPath: filepath.Join("testdata", "descriptors", "simple.descriptor"),
				fileToGenerate: "simple.proto",
			},
			wantFilePath: filepath.Join("testdata", "etalons", "simple.pb.deepcopy.go.etalone"),
		},
		{
			name: "Optionals",
			args: args{
				descriptorPath: filepath.Join("testdata", "descriptors", "optionals.descriptor"),
				fileToGenerate: "optionals.proto",
			},
			wantFilePath: filepath.Join("testdata", "etalons", "optionals.pb.deepcopy.go.etalone"),
		},
		{
			name: "Enums",
			args: args{
				descriptorPath: filepath.Join("testdata", "descriptors", "enums.descriptor"),
				fileToGenerate: "enums.proto",
			},
			wantFilePath: filepath.Join("testdata", "etalons", "enums.pb.deepcopy.go.etalone"),
		},
		{
			name: "Messages",
			args: args{
				descriptorPath: filepath.Join("testdata", "descriptors", "messages.descriptor"),
				fileToGenerate: "messages.proto",
			},
			wantErr:      false,
			wantFilePath: filepath.Join("testdata", "etalons", "messages.pb.deepcopy.go.etalone"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			req, err := protoc.ReadCodeGenerationRequest(tt.args.descriptorPath, tt.args.fileToGenerate)

			assert.NilError(t, err, "unable to create code generation request")
			assert.Assert(t, req != nil, "codegeneration request is nil")

			gen, err := protogen.Options{}.New(req)
			assert.NilError(t, err, "unable to create protogen plugin")
			gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

			if err := Generate(gen, tt.args.fileToGenerate); (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotResponse := gen.Response()
			if len(gotResponse.File) > 0 {
				gotResponse.File[0].Name = nil
			}

			expectedResponse := loadResponse(t, tt.args.fileToGenerate, tt.wantFilePath)

			assert.DeepEqual(t, expectedResponse, gotResponse, protocmp.Transform())
		})
	}
}

func loadResponse(t *testing.T, filesKV ...string) *pluginpb.CodeGeneratorResponse {
	resp := &pluginpb.CodeGeneratorResponse{
		Error:             nil,
		SupportedFeatures: proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)),
		File:              nil,
	}

	if len(filesKV)%2 == 1 {
		t.Fatal("invalid number of files specified")
	}

	// construct response for each file
	for i := 0; i < len(filesKV); i += 2 {
		data, err := os.ReadFile(filesKV[i+1])
		assert.NilError(t, err)

		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			InsertionPoint:    nil,
			Content:           proto.String(string(data)),
			GeneratedCodeInfo: nil,
		})
	}

	return resp
}
