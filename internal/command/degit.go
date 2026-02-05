package command

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/sotvokun/emit/internal/pkg/alflag"
	"github.com/sotvokun/emit/internal/service/degit"
)

var (
	DegitCommandRemoteGitHubShortcutRegexp = regexp.MustCompile(`^([a-zA-Z0-9\_\.\-]+)\/([a-zA-Z0-9\_\.\-]+)$`)
)

type DegitCommand struct {
	flagset *alflag.FlagSet

	help    *bool
	dryRun  *bool
	verbose *bool

	identity  *string
	username  *string
	secrets   *string
	noSecrets *bool
}

func NewDegitCommand() *DegitCommand {
	flagset := alflag.NewFlagSet("degit")
	help := flagset.Bool("h, help", false)

	identity := flagset.String("i", "")
	username := flagset.String("l", "")
	secrets := flagset.String("p", "")
	noSecrets := flagset.Bool("no-secrets", false)

	dryRun := flagset.Bool("dry-run", false)
	verbose := flagset.Bool("v, verbose", false)

	return &DegitCommand{
		flagset: flagset,

		help:    help,
		dryRun:  dryRun,
		verbose: verbose,

		identity:  identity,
		username:  username,
		secrets:   secrets,
		noSecrets: noSecrets,
	}
}

func (d *DegitCommand) Name() string {
	return "degit"
}

func (d *DegitCommand) Usage() string {
	return `
Usage: emit degit [OPTIONS] <remote>[#<ref>] [<destination>]

OPTIONS:
    -i <identity_file_path>    Path to the identity file to use for the SSH authentication
    -l <username>              Username to use for the authentication
    -p <secrets>               Password for the basic authentication, or the passphrase for the public key authentication
    --no-secrets               Skip the interactive secrets prompt for the authentication
    --dry-run                  Dry run the command, will not clone the repository
    -v, --verbose              Enable verbose output
    -h, --help                 Print this help message and exit

ARGUMENTS:
    <remote>                   The remote URL of a Git repository
    <ref>                      (OPTIONAL) The reference to clone (support: branch, tag, commit hash)
                               Use the HEAD reference if not specified
    <destination>              (OPTIONAL) The destination directory to clone the repository into
                               Use the current directory if not specified

AUTHENTICATION:
    Basic Authentication:
        Provide "-l" option with the username, will enable basic authentication.
    	The interactive password prompt will be shown when the "-p" or "--no-secrets" option is not provided.

    Public Key Authentication:
        Provide "-i" option with the path to the identity file, will enable public key authentication.
    	By default, the username is "git", and the interactive passphrase prompt will not be shown.
    	Once "-l" option is provided, the interactive passphrase prompt will be shown by default,
    	except the passphrase is provided with "-p" option or "--no-secrets" option is provided.
`
}

func (d *DegitCommand) Run(args []string) (int, error) {
	if err := d.flagset.Parse(args); err != nil {
		return ExitCodeInternalError, err
	}

	if *d.help {
		fmt.Fprintln(os.Stdout, strings.TrimSpace(d.Usage()))
		return ExitCodeSuccess, nil
	}

	if d.flagset.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "emit: missing remote")
		return ExitCodeArgumentError, nil
	}

	arg := d.flagset.Arg(0)
	remote, ref := d.parseArgument(arg)

	degitService, err := d.createService(remote)
	if err != nil {
		fmt.Fprintln(os.Stderr, "emit: failed to create degit service")
		return ExitCodeInternalError, err
	}

	if *d.verbose {
		logger := log.New(os.Stdout, "", log.LstdFlags)
		degitService.SetLogger(logger)
	}

	destDir := "."
	if d.flagset.NArg() >= 2 {
		destDir = d.flagset.Arg(1)
	}

	if err := degitService.Clone(ref, destDir, *d.dryRun); err != nil {
		fmt.Fprintln(os.Stderr, "emit: failed to clone the repository")
		return ExitCodeInternalError, err
	}

	return ExitCodeSuccess, nil
}

func (d *DegitCommand) createService(remote string) (*degit.DegitService, error) {
	var degitService *degit.DegitService
	if len(*d.identity) != 0 {
		if len(*d.username) == 0 {
			*d.username = "git"
		}
		passphrase := ""
		if len(*d.secrets) != 0 {
			passphrase = *d.secrets
		}
		if *d.username == "git" || (len(passphrase) == 0 && !*d.noSecrets) {
			fmt.Printf("Passphrase: ")
			fmt.Scanln(&passphrase)
		}

		var err error
		degitService, err = degit.NewDegitServiceWithPublicKey(remote, *d.identity, *d.username, passphrase)
		if err != nil {
			return nil, err
		}
	} else if len(*d.username) != 0 {
		password := ""
		if len(*d.secrets) != 0 {
			password = *d.secrets
		}
		if len(password) == 0 && !*d.noSecrets {
			fmt.Printf("Password: ")
			fmt.Scanln(&password)
		}
		degitService = degit.NewDegitServiceWithBasicAuth(remote, *d.username, password)
	} else {
		degitService = degit.NewDegitService(remote)
	}

	return degitService, nil
}

func (d *DegitCommand) parseArgument(arg string) (string, string) {
	ref := ""

	parts := strings.SplitN(arg, "#", 2)
	remote := parts[0]
	if len(parts) == 2 {
		ref = parts[1]
	}

	if DegitCommandRemoteGitHubShortcutRegexp.MatchString(remote) {
		remote = fmt.Sprintf("https://github.com/%s.git", remote)
	}

	return remote, ref
}
