
# SDKgen CLI

The CLI provides a simple binary written in go to access the https://sdkgen.app/ code generator. The following example
shows how you can use the binary.

## Build

To get the binary you can either use one of our pre-build binaries which you can download from our
[release page](https://github.com/apioo/sdkgen-cli/releases), or you can also simply build the binary
by yourself with:

> go build

## Install

The `install` command reads the `sdkgen.json` file from the current working directory and generates all defined
dependencies. Take a look at the [schema](https://sdkgen.app/schema) page to see the structure of the
`sdkgen.json` schema file. The install comand generates a `sdkgen.lock` file where all resolved specifications are
placed, so that subsequent install comand calls always use the same specification.

> sdkgen install --client-id="[user]" --client-secret="[password]"

## Update

The `update` command works almost identical to the install command except that it does not read the
`sdkgen.lock` file and always fetches the latest specification.

> sdkgen update --client-id="[user]" --client-secret="[password]"

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

This repository contains also a simple docker file to execute the generator if you dont want to use the binary directly.
To use the docker file you need to set the fitting credentials at the `docker-compose.yml` file, then you can run:

> docker-compose up

This builds the docker image and executes the `sdkgen install` command, it reads the `sdkgen.json` from the `./input`
folder and writes the generated code to the `./output` folder. You can also directly use the
`apiootech/sdkgen` [docker image](https://hub.docker.com/r/apiootech/sdkgen).
