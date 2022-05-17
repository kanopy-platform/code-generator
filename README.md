# Kanopy code-generator

An opinionated Kubernetes Runtime and Golang type generator in the builder pattern.  Why the builder pattern? It was chosen to offer a consistent, clean, and readable code that can be extensible by the developer.

The Kanopy code-generator is built using the Kubernetes [gengo](https://github.com/kubernetes/gengo) packages.

## kanopy-codegen
### Usage
```
Kanopy Builder code generator

Usage:
  kanopy-codegen [flags]

Flags:
      --bounding-dirs strings     specify directories to bound the generation
      --build-tag string          A Go build tag to use to identify files generated by this command. Should be unique. (default "ignore_autogenerated")
  -e, --go-header-file string     File containing boilerplate header text. The string YEAR will be replaced with the current 4-digit year. (default "/Users/david.katz/go/src/k8s.io/gengo/boilerplate/boilerplate.go.txt")
  -h, --help                      help for kanopy-codegen
  -i, --input-dirs strings        Comma-separated list of import paths to get input types from.
      --log-level string          Configure log level (default "info")
  -o, --output-base string        Output base; defaults to $GOPATH/src/ or ./ if $GOPATH is not set. (default "/Users/david.katz/go/src")
  -O, --output-file-base string   Base name (without .go suffix) for output files. (default "zz_generated_builders")
  -p, --output-package string     Base package path.
      --trim-path-prefix string   If set, trim the specified prefix from --output-package when generating files.
      --verify-only               If true, only verify existing output, do not write anything.
```
### Execute:
- `go install ./cmd/kanopy-code-gen`
- `kanopy-codegen -o ./<path for output> --input-dirs ./<path to package>`

## cmd/

The main entry point into the application

## internal/cli

Defines the cli interface and flags using [cobra](https://github.com/spf13/cobra)

## pkg/generators

Defines the generator implementation for constructing code for Kubernetes Runtime types

### pkg/generators/builder

Implements the [gengo](https://github.com/kubernetes/gengo) [generator.Generator](https://github.com/kubernetes/gengo/blob/master/generator/generator.go#L90) interface.

### pkg/generators/snippets

Defines individual template snippets used by the builder.

### pkg/generators/tags

Common functions to parse comment tags supported by this generator.

## Supported Comment Tags

### Type Enabled

```golang
// +kanopy:builder=true
type AType struct...
```

### Type Disabled

```golang
// +kanopy:builder=false
type AType struct...
```

## Package Global Settings

As referenced in gengo, global package settings can be provided by adding a `doc.go` to the package. For example:

`doc.go`:
```golang
package mytypes

// +kanopy:builder=package
```

## Definition of Terms

| terms | definition |
| ----- | ---------- |
| embedded | A type that is inheriting the attributes and members of another type |
| member | A direct attribute of a type. For example. `AType.Name = "hello"`  Name is a member of AType |
| root / wrapper type | The top level type defined in the source code to be generated |
| parent type | The immediate parent type of a member |

```golang
type AType struct {  // root / wrapper type
   b.AnotherType // embedded type
}

type AnotherType struct { // parent type
    Name string // member type
}
```

Relationships are always relative to the type as the tree is traversed.  In the example above the `AType` definition is also the `parent` of `b.AnotherType` and `b.AnotherType` is also a _member_ of `AType`.

Code is be generated using the following convention:

- `With<MemberName>` for direct assignments.  e.g. `WithName(string)`
- `Append<MemberName>` for slices e.g. `AppendStrings(...string)`

## Generator States

All kubernetes runtime types have the following in common. They embed [ObjectMeta](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta) and [TypeMeta](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#TypeMeta).

- Given a type is tagged and enabled Then perform code generation.
- Given a type with ObjectMeta
  - generate a Constructor that accepts the name of the resources
  - generate members of ObjectMeta not tagged as `// Read-only`

- Given a type with TypeMeta 
  - generate DeepCopy and DeepCopyInto wrappers of the parent type.

- Given a type with Builtin / Primitive members
  - generate setter functions for each member not tagged as `// Read-only`.
