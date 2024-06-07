# Azurite

Not available until the next release of testcontainers-go <a href="https://github.com/testcontainers/testcontainers-go"><span class="tc-version">:material-tag: main</span></a>

## Introduction

The Testcontainers module for Azurite.

## Adding this module to your project dependencies

Please run the following command to add the Azurite module to your Go dependencies:

```
go get github.com/testcontainers/testcontainers-go/modules/azurite
```

## Usage example

<!--codeinclude-->
[Creating a Azurite container](../../modules/azurite/examples_test.go) inside_block:runAzuriteContainer
<!--/codeinclude-->

## Module reference

The Azurite module exposes one entrypoint function to create the Azurite container, and this function receives two parameters:

```golang
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*AzuriteContainer, error)
```

- `context.Context`, the Go context.
- `testcontainers.ContainerCustomizer`, a variadic argument for passing options.

### Container Options

When starting the Azurite container, you can pass options in a variadic way to configure it.

#### Image

If you need to set a different Azurite Docker image, you can use `testcontainers.WithImage` with a valid Docker image
for Azurite. E.g. `testcontainers.WithImage("mcr.microsoft.com/azure-storage/azurite:3.23.0")`.

{% include "../features/common_functional_options.md" %}

### Container Methods

The Azurite container exposes the following methods:
