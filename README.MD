# usersms

## Introduction

User profile management micro-service

This micro-service uses Cockroach/PostgresSQL
for storage and either:
1. the `micro api` for load balancing via the proxy handler

    or
2. Google App Engine

In this Readme you will come accross the variables `<name>` and `<version>`,
both can be found in [config/consts.go](pkg/config/consts.go) where `<version>`
is the Major version code as per http://semver.org

e.g. 
```
	Name        = "usersms"
	VersionFull = "0.1.2"
```
yields

`<name>` => `usersms`

`<version>` => `0`


## Installation/Deployment

Refer to the respective readme files for instructions:

1. [Using micro](MICRO.MD)
2. [Deploy to Google AppEngine](cmd/gcloud/README.MD)

## Manual build

### Pre-requisites

1. [Go](https://golang.org) 1.8+ is required for this.
1. APIDoc installed and added to `PATH` variable: http://apidocjs.com

### Build

Run the build command from inside the root folder:
```
make build
```

`make build` can also be supplied with `goos=[GOOS]` `goarch=[GOARCH]`
`goarm=[GOARM]` values as per the documentation at https://golang.org/cmd/go (untested)

build uses the go toolchain to build binaries. It also generates API docs and
configuration templates.
