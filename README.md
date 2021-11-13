# `elevate`
[![Go Reference](https://pkg.go.dev/badge/github.com/arindas/elevate.svg)](https://pkg.go.dev/github.com/arindas/elevate)

`elevate` is a HTTP file upload server. Files are uploaded to the server's host file system.

## Usage
```
$ ./elevate --help
Usage of ./elevate:
  -base_dir string
        Base directory for storing files. (default ".")
```

## Web interface

![screenshot](./assets/screenshot.png)

## Feature Set

- Upload multiple files to the server file system, over http, at once.
- Selected files are listed before uploading.
- Target directory in the server's host file system can be specified as shown above.
- Additionally, one can specify to which subdirectory you want to upload in the html form.
- Minimal page loading time. The entire web page is self contained in a single html file.
- Manages to be one of the most barebones file upload server that does exactly what it says.
- Files in a single upload request are handled in parallel.

## Build Instructions

### Pre-requisites:

- go 1.15 or newer installed

### Commands
Simply clone the repository and build with the go tool.

```
git clone https://github.com/arindas/elevate.git
cd elevate
go build
```

The binary produced can be freely distributed with anyone using the same machine architecture.
