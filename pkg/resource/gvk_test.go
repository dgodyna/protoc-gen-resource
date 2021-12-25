package resource

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"reflect"
	"testing"
)

func Test_extractFromComments(t *testing.T) {
	type args struct {
		m *protogen.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *gvk
		want1   bool
		wantErr bool
	}{
		{
			name: "No GVK",
			args: args{
				m: &protogen.Message{
					GoIdent: protogen.GoIdent{
						GoName: "GoName",
					},
					Comments: protogen.CommentSet{
						Leading: `
Test Message
Line 1

Line 3

Line5
`,
					},
				},
			},
		},
		{
			name: "Only group specified",
			args: args{
				m: &protogen.Message{
					GoIdent: protogen.GoIdent{
						GoName: "GoName",
					},
					Comments: protogen.CommentSet{
						Leading: `
Test Message
Line 1
+protoc-gen-resource:group=TEST
Line 31

Line5
`,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Only version specified",
			args: args{
				m: &protogen.Message{
					GoIdent: protogen.GoIdent{
						GoName: "GoName",
					},
					Comments: protogen.CommentSet{
						Leading: `
Test Message
Line 1
+protoc-gen-resource:version=TEST
Line 3

Line5
`,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "group, version, no kind",
			args: args{
				m: &protogen.Message{
					GoIdent: protogen.GoIdent{
						GoName: "GoName",
					},
					Comments: protogen.CommentSet{
						Leading: `
Test Message
Line 1
+protoc-gen-resource:version=TEST_comment
+protoc-gen-resource:group=TEST_comment
Line 3

Line5
`,
					},
				},
			},
			want: &gvk{
				Group:   "TEST_comment",
				Version: "TEST_comment",
				Kind:    "GoName",
			},
			want1: true,
		},
		{
			name: "group, version, kind",
			args: args{
				m: &protogen.Message{
					GoIdent: protogen.GoIdent{
						GoName: "GoName",
					},
					Comments: protogen.CommentSet{
						Leading: `
Test Message
Line 1
+protoc-gen-resource:version=TEST_comment
+protoc-gen-resource:group=TEST_comment
+protoc-gen-resource:kind=TEST_comment
Line 3

Line5
`,
					},
				},
			},
			want: &gvk{
				Group:   "TEST_comment",
				Version: "TEST_comment",
				Kind:    "TEST_comment",
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := extractFromComments(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractFromComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractFromComments() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractFromComments() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_extractFromPackage(t *testing.T) {
	type args struct {
		packageStr string
		m          *protogen.Message
	}
	tests := []struct {
		name  string
		args  args
		want  *gvk
		want1 bool
	}{
		{
			name: "too short",
			args: args{
				packageStr: "mycompany",
			},
		},
		{
			name: "not fit version",
			args: args{
				packageStr: "com.mycompany.apiv2",
			},
		},
		{
			name: "ends by version",
			args: args{
				packageStr: "com.mycompany.api.v1alpha1",
				m: &protogen.Message{GoIdent: protogen.GoIdent{
					GoName: "TestResource",
				}},
			},
			want: &gvk{
				Group:   "api.mycompany.com",
				Version: "v1alpha1",
				Kind:    "TestResource",
			},
			want1: true,
		},
		{
			name: "ends by hub",
			args: args{
				packageStr: "com.mycompany.api.hub",
				m: &protogen.Message{GoIdent: protogen.GoIdent{
					GoName: "TestResource",
				}},
			},
			want: &gvk{
				Group:   "api.mycompany.com",
				Version: "hub",
				Kind:    "TestResource",
			},
			want1: true,
		},
		{
			name: "model",
			args: args{
				packageStr: "com.mycompany.api.hub.model",
				m: &protogen.Message{GoIdent: protogen.GoIdent{
					GoName: "TestResource",
				}},
			},
			want: &gvk{
				Group:   "api.mycompany.com",
				Version: "hub",
				Kind:    "TestResource",
			},
			want1: true,
		},
		{
			name: "services",
			args: args{
				packageStr: "com.mycompany.api.hub.services",
				m: &protogen.Message{GoIdent: protogen.GoIdent{
					GoName: "TestResource",
				}},
			},
			want: &gvk{
				Group:   "api.mycompany.com",
				Version: "hub",
				Kind:    "TestResource",
			},
			want1: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			f := &protogen.File{
				Proto: &descriptorpb.FileDescriptorProto{
					Package: &tt.args.packageStr,
				},
			}
			got, got1 := extractFromPackage(*f.Proto.Package, tt.args.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractFromPackage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractFromPackage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
