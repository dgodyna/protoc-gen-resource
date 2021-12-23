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

func TestGenerate(t *testing.T) {
	type args struct {
		fileToGenerate string
		filePath       []string
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
				fileToGenerate: "simple.proto_",
				filePath:       []string{filepath.Join("testdata", "protos")},
			},
			wantFilePath: filepath.Join("testdata", "etalons", "simple.pb.deepcopy.go.etalone"),
		},
		{
			name: "Optionals",
			args: args{
				fileToGenerate: "optionals.proto_",
				filePath:       []string{filepath.Join("testdata", "protos")},
			},
			wantFilePath: filepath.Join("testdata", "etalons", "optionals.pb.deepcopy.go.etalone"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := protoc.Parser{
				ProtoPaths:        tt.args.filePath,
				IncludeImports:    true,
				IncludeSourceInfo: true,
			}.CodeGenerationRequest(tt.args.fileToGenerate)

			assert.NilError(t, err, "unable to create code generation request")

			gen, err := protogen.Options{}.New(req)
			assert.NilError(t, err, "unable to create protogen plugin")
			gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

			if err := Generate(gen, tt.args.fileToGenerate); (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}

			resp := gen.Response()
			if len(resp.File) > 0 {
				resp.File[0].Name = nil
			}

			response := loadResponse(t, tt.args.fileToGenerate, tt.wantFilePath)

			assert.DeepEqual(t, response, resp, protocmp.Transform())
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
