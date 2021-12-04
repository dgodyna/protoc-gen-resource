package resource

import (
	"google.golang.org/protobuf/compiler/protogen"
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
Line 3

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
