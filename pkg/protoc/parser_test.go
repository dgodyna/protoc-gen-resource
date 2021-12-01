package protoc

import (
	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
	"os/exec"
	"testing"
)

func TestParser_protocCommand(t *testing.T) {
	type fields struct {
		ProtoPaths        []string
		DescriptorsSet    []string
		IncludeImports    bool
		IncludeSourceInfo bool
		Protoc            string
	}
	type args struct {
		filesToGenerate []string
	}
	tests := []struct {
		name           string
		fields         fields
		customizedPath string
		args           args
		want           *exec.Cmd
		wantErr        bool
	}{
		{
			name:    "Protoc from path",
			fields:  fields{},
			args:    args{},
			want:    exec.Command("protoc"),
			wantErr: false,
		},
		{
			name:           "Configured protoc",
			fields:         fields{},
			customizedPath: "",
			args:           args{},
			want:           nil,
			wantErr:        false,
		},
		{
			name:           "ProtoPaths",
			fields:         fields{},
			customizedPath: "",
			args:           args{},
			want:           nil,
			wantErr:        false,
		},
		{
			name:           "DescriptorsSet",
			fields:         fields{},
			customizedPath: "",
			args:           args{},
			want:           nil,
			wantErr:        false,
		},
		{
			name:           "IncludeImports",
			fields:         fields{},
			customizedPath: "",
			args:           args{},
			want:           nil,
			wantErr:        false,
		},
		{
			name:           "IncludeSourceInfo",
			fields:         fields{},
			customizedPath: "",
			args:           args{},
			want:           nil,
			wantErr:        false,
		},
		{
			name:           "A bit of everything",
			fields:         fields{},
			customizedPath: "",
			args:           args{},
			want:           nil,
			wantErr:        false,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := Parser{
				ProtoPaths:        tt.fields.ProtoPaths,
				DescriptorsSet:    tt.fields.DescriptorsSet,
				IncludeImports:    tt.fields.IncludeImports,
				IncludeSourceInfo: tt.fields.IncludeSourceInfo,
				Protoc:            tt.fields.Protoc,
			}
			got, err := p.protocCommand(tt.args.filesToGenerate...)
			if (err != nil) != tt.wantErr {
				t.Errorf("protocCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Assert(t, cmp.Equal(tt.want, got))
		})
	}
}
