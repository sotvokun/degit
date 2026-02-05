Emit
====

Emit is a powerful plain text project scaffolding tool. It offers a batch of commands to make plain text content fast.

## Usage
### Degit

Degit is a command to download a remote Git repository without their `.git` folder and history.

The simplest usage of degit is to download the main branch of a repository to the current working directory:
```sh
emit degit https://github.com/user/repo
emit degit git@gitlab.com:user/repo

# A shortcut for github
emit degit user/repo
```

**Degit into a specified path**
```sh
emit degit user/repo new-project-folder
```

**Specify a tag, branch (commit hash has not support yet)**
```sh
emit degit user/repo#branch # branch
emit degit user/repo#tag    # tag
```

For more information, please read the help message by
```sh
emit degit --help
```
