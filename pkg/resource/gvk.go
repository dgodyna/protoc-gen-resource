package resource

import (
	_ "embed"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

//go:embed templates/gvk.gotmpl
var gvkTmpl string

type gvk struct {
	Group   string
	Version string
	Kind    string
}

// getGVK get group version & kind of resource from proto message.
func getGVK(m *protogen.Message) (*gvk, error) {

	//m.Comments
	return nil, nil
}

// extractFromComments will extract group version kind information from protobuf message comments.
// Group must be specified by following comment: +protoc-gen-resource:group=GROUP
// Version must be specified by following comment: +protoc-gen-resource:version=VERSION
// King may be specified by following comment: +protoc-gen-resource:kind=KIND
// If kind is not specified - message name will be used.
// If any of group or version specified without another one - error will be returned.
func extractFromComments(m *protogen.Message) (*gvk, bool, error) {

	groupFound := false
	versionFound := false
	kindFound := false

	if strings.Contains(string(m.Comments.Leading), "+protoc-gen-resource:group=") {
		groupFound = true
	}
	if strings.Contains(string(m.Comments.Leading), "+protoc-gen-resource:version=") {
		versionFound = true
	}
	if strings.Contains(string(m.Comments.Leading), "+protoc-gen-resource:kind=") {
		kindFound = true
	}

	if (groupFound != versionFound) && (groupFound || versionFound) {
		return nil, false, fmt.Errorf("invalid configuration for GVK, both comments '+protoc-gen-resource:group=GROUP' "+
			"and '+protoc-gen-resource:version=VERSION' must be provided for message '%s'", m.GoIdent.GoName)
	}

	if !groupFound {
		return nil, false, nil
	}

	commentLine := string(m.Comments.Leading)
	// comments are full of new lines and tabs - replace it with space
	commentLine = strings.ReplaceAll(commentLine, "\n", " ")
	commentLine = strings.ReplaceAll(commentLine, "\t", " ")
	// add space to the end to have ability to split right part of comments by spaces
	commentLine = commentLine + " "

	// now let's extract group, version
	group := strings.Split(strings.Split(commentLine, "+protoc-gen-resource:group=")[1], " ")[0]
	version := strings.Split(strings.Split(commentLine, "+protoc-gen-resource:version=")[1], " ")[0]

	kind := m.GoIdent.GoName
	if kindFound {
		kind = strings.Split(strings.Split(commentLine, "+protoc-gen-resource:kind=")[1], " ")[0]
	}

	return &gvk{
		Group:   group,
		Version: version,
		Kind:    kind,
	}, true, nil
}

// extractFromPackage will try to extract group, version, kind from protobuf package.
// Package must match following patterns <GROUP>.<VERSION> or <GROUP>.<VERSION>.[services|model
// Where version must be either 'hub' - for internal version or follow v.* pattern.
// Group will be reversed:
// com.mycompany.product1.api -> api.product1.mycompany.com
// If package is not following patterns - return false.
func extractFromPackage(f *protogen.File, m *protogen.Message) (*gvk, bool) {

	p := *f.Proto.Package

	p = strings.TrimSuffix(p, ".model")
	p = strings.TrimSuffix(p, ".services")

	packageParts := strings.Split(p, ".")

	if len(packageParts) < 2 {
		// does not follow proposed format
		return nil, false
	}

	version := packageParts[len(packageParts)-1]
	// checks version
	if version != "hub" && !strings.HasPrefix(version, "v") {
		return nil, false
	}

	// reverse package parts
	for i, j := 0, len(packageParts)-1; i < j; i, j = i+1, j-1 {
		packageParts[i], packageParts[j] = packageParts[j], packageParts[i]
	}

	return &gvk{
		Group:   strings.Join(packageParts[1:], "."),
		Version: version,
		Kind:    m.GoIdent.GoName,
	}, true
}
