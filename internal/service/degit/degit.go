package degit

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/sotvokun/emit/internal/service/log"
)

type WalkFunc func(path string, fs billy.Filesystem) error

type DegitService struct {
	remote     string
	authMethod transport.AuthMethod

	logger log.Logger
}

func NewDegitService(remote string) *DegitService {
	return &DegitService{
		remote: remote,
	}
}

func NewDegitServiceWithBasicAuth(remote string, username string, password string) *DegitService {
	auth := &http.BasicAuth{
		Username: username,
		Password: password,
	}
	return &DegitService{
		remote:     remote,
		authMethod: auth,
	}
}

func NewDegitServiceWithPublicKey(remote string, publicKeyPath string, auth ...string) (*DegitService, error) {
	username := ""
	password := ""
	if len(auth) >= 1 {
		username = auth[0]
	}
	if len(auth) >= 2 {
		password = auth[1]
	}

	authMethod, err := ssh.NewPublicKeysFromFile(username, publicKeyPath, password)
	if err != nil {
		return nil, err
	}
	return &DegitService{
		remote:     remote,
		authMethod: authMethod,
	}, nil
}

func (d *DegitService) SetLogger(logger log.Logger) {
	d.logger = logger
}

func (d *DegitService) Clone(ref string, destDir string, dryMode bool) error {
	refObj, err := d.getReference(ref)
	if err != nil {
		return err
	}

	memfs := memfs.New()
	_, err = git.Clone(memory.NewStorage(), memfs, &git.CloneOptions{
		URL:               d.remote,
		ReferenceName:     refObj.Name(),
		SingleBranch:      true,
		Depth:             1,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              d.authMethod,
		ShallowSubmodules: true,
	})
	if err != nil {
		return err
	}

	walker := NewWalker(destDir)
	walker.SetLogger(d.logger)
	walker.SetDryMode(dryMode)
	if err := d.walk(memfs, "", walker.WalkCopy); err != nil {
		return err
	}

	return nil
}

// getReference returns the reference with the given ref. The `ref` can be a branch, tag, or commit hash.
// The HEAD reference will be returned if no ref is provided. `nil` will be returned if the reference is not found.
func (d *DegitService) getReference(ref string) (*plumbing.Reference, error) {
	refs, err := d.references()
	if err != nil {
		return nil, err
	}

	d.log("found %d references", len(refs))

	for _, r := range refs {
		// Find HEAD reference for no ref provided
		if len(ref) == 0 {
			if r.Name() == plumbing.HEAD {
				return r, nil
			}
			continue
		}

		// Find reference by hash
		if strings.HasPrefix(r.Hash().String(), ref) {
			return r, nil
		}

		// Find reference by branch name
		if r.Name().IsBranch() && r.Name() == plumbing.NewBranchReferenceName(ref) {
			return r, nil
		}

		// Find reference by tag name
		if r.Name().IsTag() && r.Name() == plumbing.NewTagReferenceName(ref) {
			return r, nil
		}
	}

	return nil, fmt.Errorf("reference '%s' not found", ref)
}

// references returns the references of the remote repository
func (d *DegitService) references() ([]*plumbing.Reference, error) {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{d.remote},
	})
	return remote.List(&git.ListOptions{
		Auth: d.authMethod,
	})
}

// walk walks the filesystem and calls the given function for each non-directory entry
func (d *DegitService) walk(fs billy.Filesystem, root string, fn WalkFunc) error {
	entries, err := fs.ReadDir(root)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(root, entry.Name())
		if entry.IsDir() {
			if err := d.walk(fs, fullPath, fn); err != nil {
				return err
			}
			continue
		}
		if err := fn(fullPath, fs); err != nil {
			return err
		}
	}
	return nil
}

func (d *DegitService) log(msg string, a ...any) {
	if d.logger == nil {
		return
	}
	d.logger.Printf(msg+"\n", a...)
}
