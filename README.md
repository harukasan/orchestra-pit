orchestra-pit
=============

**STILL UNDER CONSTRUCTION**

**orchestra-pit** is a library to build system orchestration tools.
It also includes orchestration command named **opit**.

- written in Go
- consists just one binary and the recipe file
- no longer dependents on Ruby, Python or any interpreters

## How it works

```
$ vim recipe.json
{
  "resources": [
    {
      "type": "package",
      "name": "sl",
      "state": "installed"
    }
  ]
}

$ opit
Usage: opit COMMAND [ARGUMENTS]

opit is a command line interface of orchestra-pit.

Commands:
  apply    apply the recipe to the host
  test     test whether states of the host satisfies the given recipe file
  version  print version string

Use "opit COMMAND -h" for more information about the commands.

$ opit test
Started at 2015-07-23T10:50:30+09:00
[FAIL] the package "sl" is not installed

$ opit apply
Started at 2015-07-23T10:51:16+09:00
[DONE] name: "sl", state: "installed"

```

## TODO

- supports yaml format
- supports template variables
- supports mrb?
- more tests
- documentation
- auto document generation

## Influenced works

- [serverkit/serverkit](https://github.com/serverkit/serverkit)
- [serverspec/specinfra](https://github.com/serverspec/specinfra)
- [chef/chef](htps://github.com/chef/chef)
