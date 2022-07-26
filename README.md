# go-openapi-transform

Parses JSON or YAML encoded OpenAPI specification files into OpenAPI model objects for the purpose of testing the compatiblility of OpenAPI files and the OpenAPI parser ([neotoolkit/openapi](https://pkg.go.dev/github.com/neotoolkit/openapi)).

## Build

To build the application, run `go build`

## Configuration

Place your OpenAPI specification files in the `./docs` directory.

Configure the application by editing the `config.json` file.

|Property|Type|Description|
|---|---|---|
|ignoreFiles|Array\<String\>|List of filenames to ignore|
|whitelist|Array\<String\>|List of filenames to transform. When present, only these files will be parsed and the `ignoreFiles` configuration will be overridden.|

## Run

Run the application with `go run .`
