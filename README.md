# Protoc Gen Resources

Generates following methods to match k8s runtime.Object interface:

* `GetResourceGroup() string`
* `GetResourceVersion() string`
* `GetResourceKind() string`
* `GetObjectKind() schema.ObjectKind`
* `DeepCopyInto`
* `DeepCopy`
* `DeepCopyObject() runtime.Object`

Supported proto3 types:

- [x] scalars
- [x] optional scalars
- [ ] enums
- [ ] messages
- [ ] 3rd party messages
- [ ] repeated fields
- [ ] maps
- [ ] oneOf