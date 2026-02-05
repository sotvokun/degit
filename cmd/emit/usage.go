package main

func usage() string {
	return `
Usage: emit [-h | --help] [-v | --version] <command> [<arguments>]

Options:
    -h, --help          Print this help message and exit
    -v, --version       Print the version of the program

Commands:
    version             Display version information about emit
    degit               Clone a repository from a remote URL
`
}
