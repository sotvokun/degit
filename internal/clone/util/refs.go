package util

import (
	"degit/internal/clone/git"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
)

func SearchRef(refs []*git.Reference, search ...string) (*git.Reference, error) {
	if len(search) == 0 || len(strings.TrimSpace(search[0])) == 0 {
		return DefaultRef(refs)
	}
	searchText := search[0]
	for _, ref := range refs {
		if ref.Hash().String() == searchText {
			return ref, nil
		}
		tagRef := plumbing.NewTagReferenceName(searchText)
		branchRef := plumbing.NewBranchReferenceName(searchText)
		if ref.Name() == tagRef || ref.Name() == branchRef {
			return ref, nil
		}
	}
	return nil, fmt.Errorf("could not find ref %s", searchText)
}

func DefaultRef(refs []*git.Reference) (*git.Reference, error) {
	for _, ref := range refs {
		if ref.Name() == plumbing.HEAD {
			return ref, nil
		}
	}
	return nil, fmt.Errorf("no default ref found")
}
