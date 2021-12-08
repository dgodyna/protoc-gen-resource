package main

import (
	"fmt"
	"github.com/dgodyna/protoc-gen-resource/pkg/protoc"
	"github.com/dgodyna/protoc-gen-resource/pkg/resource"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io"
	"os"
)

func main() {
	panic("hhhhh")
	resp := generate()
	err := writeResponse(resp)
	if err != nil {
		panic(err)
	}
}

func generate() *pluginpb.CodeGeneratorResponse {
	req, err := parseProtocRequest()
	if err != nil {
		return &pluginpb.CodeGeneratorResponse{
			Error: proto.String(fmt.Errorf("unable to parse protoc request : %w", err).Error()),
		}
	}

	return protoc.ApplyPluginFunction(resource.Generate, req)
}

// writeResponse marshall response and write it to stdout
func writeResponse(resp *pluginpb.CodeGeneratorResponse) error {

	out, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("unable to marshall codegeneration response: %w", err)
	}

	_, err = os.Stdout.Write(out)
	if err != nil {
		return fmt.Errorf("unable to write codegeneration response to stdout : %w", err)
	}

	return nil
}

// parseProtocRequest parse generation request from stdin and unmarshall in to generator request.
func parseProtocRequest() (*pluginpb.CodeGeneratorRequest, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("unable to read codegeneration request : %w", err)
	}

	req := &pluginpb.CodeGeneratorRequest{}
	err = proto.Unmarshal(data, req)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall codegeneration request : %w")
	}

	return req, nil
}
