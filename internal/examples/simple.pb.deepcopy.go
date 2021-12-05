// Code generated by protoc-gen-resource. DO NOT EDIT.

package examples

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (*ABitOfScalars) GetResourceGroup() string {
	return "test.api.nrm.netcracker.com"
}

// API Version, equals to "hub"
func (*ABitOfScalars) GetResourceVersion() string {
	return "hub"
}

// Resource Kind, equals to "ABitOfScalars"
func (*ABitOfScalars) GetResourceKind() string {
	return "ABitOfScalars"
}

// GetObjectKind to satisfy runtime.Object interface
func (x *ABitOfScalars) GetObjectKind() schema.ObjectKind {
	typeMeta := meta.TypeMeta{}
	typeMeta.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "test.api.nrm.netcracker.com",
		Version: "hub",
		Kind:    "ABitOfScalars",
	})
	return &typeMeta
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ABitOfScalars) DeepCopyInto(out *ABitOfScalars) {
	in.DoubleType = out.DoubleType
	in.FloatType = out.FloatType
	in.Int32Type = out.Int32Type
	in.Int64Type = out.Int64Type
	in.Uint32Type = out.Uint32Type
	in.Uint64Type = out.Uint64Type
	in.Sint32Type = out.Sint32Type
	in.Sint64Type = out.Sint64Type
	in.Fixed32Type = out.Fixed32Type
	in.Fixed64Type = out.Fixed64Type
	in.Sfixed32Type = out.Sfixed32Type
	in.Sfixed64Type = out.Sfixed64Type
	in.BoolType = out.BoolType
	in.StringType = out.StringType
	in.BytesType = out.BytesType
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *ABitOfScalars) DeepCopy() *ABitOfScalars {
	if in == nil {
		return nil
	}
	out := new(ABitOfScalars)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, returning deepcopy as runtime.Object interface.
func (in *ABitOfScalars) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
