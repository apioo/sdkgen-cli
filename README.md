
# SDKgen CLI

The CLI provides a simple binary written in go to access the https://sdkgen.app/ code generator. The following example
shows how you can use the binary.

## Build

To get the binary you can either use one of our pre-build binaries which you can download from our
[release page](https://github.com/apioo/sdkgen-cli/releases), or you can also simply build the binary
by yourself with:

> go build

## Generate

The generate command expects as first type a type which SDK should be generated i.e. `client-typescript`.
The second argument is a [TypeAPI](https://typeapi.org/) specification which describes your REST API.
The last argument is a target directory where all generated files are placed.

> sdkgen generate [type] [schema_file] [output_dir] --client-id="[user]" --client-secret="[password]"

* __type__  
  The type describes what kind of SDK should be generated i.e. TypeScript or Java. You can get a list of all available types with sdkgen types.
* __schema_file__  
  The input [TypeAPI](https://typeapi.org/) specification which describes your REST API.
* __output_dir__  
  The target output dir where all generated files are placed, this directory should be empty.
* __--client-id__  
  The client id is either your username or an app key which you can create at our backend.
* __--client-secret__  
  This client secret is either your password or an app secret which you can create at our backend.

## Types

Lists all available types which can be used at the generate command.

> sdkgen types --client-id="[user]" --client-secret="[password]"

* __--client-id__  
  The client id is either your username or an app key which you can create at our backend.
* __--client-secret__  
  This client secret is either your password or an app secret which you can create at our backend.

## Docker

This repository contains also a simple docker file to execute the generator if you dont want to use the binary
directly. To use the docker file you need to set the fitting credentials at the `.env` file, then you can run:

> docker-compose up

This builds the docker image and reads the `./output/typeapi.json` specification and writes the generated code also to
the `./output` dir.
