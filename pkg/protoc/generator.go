package protoc

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// GenerateFileFunc is function that takes generator as input and generate each of the files
type GenerateFileFunc func(gen *protogen.Plugin, file string) error

// ApplyPluginFunction applies `f` to all FileToGenerate inside request
//
// Returns combined plugin response
func ApplyPluginFunction(f GenerateFileFunc, req *pluginpb.CodeGeneratorRequest) *pluginpb.CodeGeneratorResponse {
	gen, err := protogen.Options{}.New(req)
	if err != nil {
		return &pluginpb.CodeGeneratorResponse{
			Error: proto.String(err.Error()),
		}
	}
	// Add support for proto3 optional fields
	// https://github.com/protocolbuffers/protobuf/blob/master/docs/implementing_proto3_presence.md#signaling-that-your-code-generator-supports-proto3-optional
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	for _, file := range gen.Request.FileToGenerate {
		err := f(gen, file)
		if err != nil {
			gen.Error(err)
			break
		}
	}

	return gen.Response()
}
