# degit — straightforward project scaffolding

## Introduction

degit is a project scaffolding tool that offers two main features:
- Downloads Git repositories without their `.git` folder and history
- Provides a powerful templating system using Go's `text/template` package

This is a Go implementation of the original [degit](https://github.com/Rich-Harris/degit) tool, enhanced with template processing capabilities for more flexible project scaffolding.

## Usage

### Basics

The simplest use of degit is to download the main branch of a repository from any remote git repository to the current working directory:
```sh
degit https://github.com/user/repo
degit git@gitlab.com:user/repo

# A shortcut for github
degit user/repo
```

#### Create a new folder for the project
```sh
degit user/repo my-new-project
```

#### Specify a tag, branch (commit hash has not support yet)
```sh
degit user/repo#dev         # branch
degit user/repo#v1.2.3      # tag
```

#### Private repository
```sh
degit -i path/to/ssh/private-key user/repo          # auth with key
degit -l identifier -p secret user/repo             # auth with username/password/access token
```

### Template

degit also can be use as the template engine commandline interface to render template file(s).

```sh
# rendering template and replace template file with result
degit template path/of/template

# rendering template but put result into a new file
degit template path/of/template path/of/result

# rendering with given variables
degit template -D foo=bar -D abc=xyz path/of/template
```

#### Options

```sh
degit template -s key=value -s key2=value2 path/of/template
```

| Key                | Example                                    | Description   |
|--------------------|--------------------------------------------|---------------|
| extions            | `extensions=template` `extensions=temp,in` | Used to define the path of result but only remove the extra extension, it also use for 'glob' command line option; separated by comma. |
| delimiter          | `delimiter=<%,%>` `delimiter={%,%}` | Set the action delimiters of template to specific strings, left and right sides separated by comma. |
| nonstrict          | `nonstrict=true` `nonstrict=TRUE` | Use the "zero value" for variables that not given. The value is 'true'; case insensitive |
| removesource       | `removesource=true`               | Remove the "source file" if the file path of result is not same with the path of "source file" |

#### Rendering with glob mode

```sh
# Use '-g' option to enable glob mode, all arguments will be as glob pattern
degit template -g "src/**/*.js" "package.json"
```
