<p align="center">
<img src="resources/icons/app.png" style="width:150px" width="150" alt="Warthog" />
<h2 align="center">Warthog</h2>
<h4 align="center">Cross platform gRPC client</h4>
</p>
<p align="center">
<a href="https://goreportcard.com/report/github.com/forest33/warthog"><img src="https://goreportcard.com/badge/github.com/forest33/warthog"/></a>
<img src="https://img.shields.io/github/go-mod/go-version/forest33/warthog?style=flat-square"/>
</p>


## Features

- Automatic parsing of proto definitions to render services and input messages
- `.proto` file discovery
- Selection of multiple services and methods
- Configuration of TLS, including disabling TLS (plain text)
- Input generation for all scalar types
- Input generation for nested messages
- Input generation for enums, including nested
- Input generation for repeated fields
- Input generation for oneof and map fields
- Support for adding RPC metadata
- Execute unary requests
- Execute server streaming requests
- Execute client streaming requests
- Execute bidirectional requests
- Cancel requests in-flight
- View response messages
- View RPC Header and Trailer
- MacOS build
- Linux build
- Windows build
- Reflection API to determine RPC schema
- Support for Google Well Known Types
- Create multiple workspaces and workspace switching

## Download

Visit the [Releases](https://github.com/Forest33/warthog/releases) page for the latest downloads.

## Install

### Linux

[Download](https://github.com/Forest33/warthog/releases) and unarchive `Warthog*-linux-x86-64.tar.gz` and run.

### MacOS

[Download](https://github.com/Forest33/warthog/releases) and open `Warthog*-darwin-x86-64.dmg`, drag `Warthog` to the `Applications` folder and run from `Applications`.

### Windows

[Download](https://github.com/Forest33/warthog/releases) and run `Warthog*-windows-x86-64.exe`.