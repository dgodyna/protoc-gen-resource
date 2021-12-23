package protoc

import (
	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/assert"
	"os"
	"os/exec"
	"testing"
)

func TestParser_protocCommand(t *testing.T) {
	t.Skip("must be run only on envs with protoc installed")
	protocPath, err := exec.LookPath("protoc")
	assert.NilError(t, err)

	type fields struct {
		ProtoPaths        []string
		DescriptorsSet    []string
		DescriptorsSetOut string
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
			name: "Protoc from path",
			want: exec.Command(protocPath),
		},
		{
			name: "Configured protoc",
			fields: fields{
				Protoc: "/test/protoc",
			},
			customizedPath: "/t",
			want:           exec.Command("/test/protoc"),
		},
		{
			name:           "Non-existing protoc",
			customizedPath: "/test",
			wantErr:        true,
		},
		{
			name: "ProtoPaths",
			fields: fields{
				ProtoPaths: []string{"dir1/subdir", "dir2", "dir3"},
				Protoc:     "/test/protoc",
			},
			want: exec.Command("/test/protoc", "-Idir1/subdir", "-Idir2", "-Idir3"),
		},
		{
			name: "DescriptorsSet",
			fields: fields{
				DescriptorsSet: []string{"dir1/subdir", "dir2", "dir3"},
				Protoc:         "/test/protoc",
			},
			want: exec.Command("/test/protoc", "--descriptor_set_in=dir1/subdir,dir2,dir3"),
		},
		{
			name: "IncludeImports",
			fields: fields{
				IncludeImports: true,
				Protoc:         "/test/protoc",
			},
			want: exec.Command("/test/protoc", "--include_imports"),
		},
		{
			name: "IncludeSourceInfo",
			fields: fields{
				IncludeSourceInfo: true,
				Protoc:            "/test/protoc",
			},
			want: exec.Command("/test/protoc", "--include_source_info"),
		},
		{
			name: "descriptor_set_out",
			fields: fields{
				DescriptorsSetOut: "test_out",
				Protoc:            "/test/protoc",
			},
			want: exec.Command("/test/protoc", "--descriptor_set_out=test_out"),
		},
		{
			name: "Params",
			fields: fields{
				Protoc: "/test/protoc",
			},
			args: args{filesToGenerate: []string{"file1.proto", "file2.proto", "file3.proto"}},
			want: exec.Command("/test/protoc", "file1.proto", "file2.proto", "file3.proto"),
		},
		{
			name: "A bit of everything",
			fields: fields{
				ProtoPaths:        []string{"dir1/subdir", "dir2", "dir3"},
				IncludeImports:    true,
				IncludeSourceInfo: true,
				Protoc:            "/test/protoc",
			},
			args: args{filesToGenerate: []string{"file1.proto", "file2.proto", "file3.proto"}},
			want: exec.Command("/test/protoc", "-Idir1/subdir", "-Idir2", "-Idir3", "--include_imports", "--include_source_info", "file1.proto", "file2.proto", "file3.proto"),
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.customizedPath != "" {
				path := os.Getenv("PATH")
				err := os.Setenv("PATH", tt.customizedPath)
				assert.NilError(t, err)
				defer func() {
					err := os.Setenv("PATH", path)
					assert.NilError(t, err)
				}()
			}

			p := Parser{
				ProtoPaths:        tt.fields.ProtoPaths,
				DescriptorsSet:    tt.fields.DescriptorsSet,
				IncludeImports:    tt.fields.IncludeImports,
				IncludeSourceInfo: tt.fields.IncludeSourceInfo,
				Protoc:            tt.fields.Protoc,
				DescriptorsSetOut: tt.fields.DescriptorsSetOut,
			}
			got, err := p.protocCommand(tt.args.filesToGenerate...)
			if (err != nil) != tt.wantErr {
				t.Errorf("protocCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.DeepEqual(t, tt.want, got, cmpopts.IgnoreUnexported(exec.Cmd{}))
		})
	}
}
