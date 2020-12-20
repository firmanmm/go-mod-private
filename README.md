# Go Mod Private
A simple helper toolkit to help development with Private Git Server using Golang Module. It allows you to use Git Private Server for your vendor. You can also use `gomp get` just as `go get`. 

*Note* : Currently `gomp` didn't support `./...` operation.

## Installation
Before installing `gomp`, you have to make sure that you have `git` installed on your system.
To install `gomp` command line interface, simply use syntax below :
```
go get github.com/firmanmm/go-mod-private/cmd/gomp
```

## Usage
If you don't know what this tool does, simply add `-h` when using `gomp` cli. 

Example : 
```
$ gomp -h
NAME:
   GoModPrivate - Go Module for Private Git Server Repository

USAGE:
   gomp.exe [global options] command [command options] [arguments...]

VERSION:
   0.0.0

DESCRIPTION:
   Allow you to use Git Private Server as your source dependency while using Go Module

AUTHOR:
   Firman "Rendoru" Maulana 

COMMANDS:
   get             Perform Go Get operation, will switch to git clone and git pull if a matching SSH Credential has been registered
   add_credential  Add credential to be used for [get] command
   sync            Synchronize vendor.gomp and go.mod to mod.gomp
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --gomp value   Read from given gomp file (default: "mod.gomp")
   --help, -h     show help
   --version, -v  print the version
```
### Getting Package
If you want to use `gomp` just like using `go get` you just need to replace `go get` with `gomp get`.
Syntax : 
```
$ gomp get [arguments...]
```
Example :
```
$ gomp get github.com/firmanmm/suberror
```
`gomp` will try to scan your `mod.gomp` file to determine any `SshCredential` that match your requested package. `gomp` will match with the longest matching `SshCredential` if there are multiple credentials that can be used by your request. Upon match, `gomp` will perform `git clone` internally to get your desired package. If there are **no credentials** that match your request, `gomp` will use `go get` to get your request.

#### How it works
- `gomp` try to find any matching credential
    - If nothing match, `gomp` will use `go get`
- `gomp` will perform `git clone` using your `SshCredentials` instead `go get` 
- `gomp` will add the cloned repository to `vendor.gomp`
- `gomp` will then update your `go.mod` to use files from `vendor.gomp`. It will only update your `gomp get`'ed repository
- `gomp` will keep track of your private repository in `mod.gomp`

### Adding Credentials
`gomp` support adding credentials to `mod.gomp` via CLI. This way, you don't have to edit your `mod.gomp` manually.
Syntax :
```
$ gomp add_credential --host=[Required] --user=[Required] --base=[Optional] --pattern=[Optional]
```
- `--host` : Your target host. **Required**
- `--user` : Your target user. **Required**
- `--base` : Base path for searching for package.
- `--pattern` : Will be used to match your requested package with this credential.
Example : 
```
$ gomp add_credential --host=rendoru.com --user=someone --base=/home/someone --pattern="rendoru.com(.*)"
```
### Synchronizing
`gomp` can automatically synchronize `vendor.gomp` and `go.mod` with the `mod.gomp`. This way you don't have to manually `get` them one by one.
Syntax : 
``` 
$ gomp sync
```

## Example
Below are the example of Go Mod Private's output

### Example mod.go
```
module github.com/firmanmm/go-mod-private

go 1.12

require (
	github.com/go-ini/ini v1.46.0 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	github.com/urfave/cli v1.21.0
	gopkg.in/ini.v1 v1.46.0 // indirect
)

//GO_MOD_PRIVATE_START
//This is an auto generated section made by Go Mod Private
//For more information visit https://github.com/firmanmm/go-mod-private
//Please add vendor.gomp to your .gitignore

replace (
	rendoru.com/module/sync-mq => ./vendor.gomp/rendoru.com/module/sync-mq
	rendoru.com/tool/concurrent => ./vendor.gomp/rendoru.com/tool/concurrent

)

//GO_MOD_PRIVATE_END
```

#### Example mod.gomp
```
{
    "SshCredentials": [
        {
            "Matcher": "rendoru.com(.*)",
            "Host": "rendoru.com",
            "Username": "someone",
            "BasePath": "/home/someone"
        }
    ],
    "PrivateRepositories": [
        "rendoru.com/module/sync-mq",
        "rendoru.com/tool/concurrent"
    ]
}
```
