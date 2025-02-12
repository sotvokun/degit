package git

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Auth struct {
	PrivateKey string // Path to private key
	Username   string
	Password   string
}

func (a *Auth) Transport() (transport.AuthMethod, error) {
	if a.PrivateKey != "" {
		_, err := os.Stat(a.PrivateKey)
		if err != nil {
			return nil, err
		}

		publicKeyUser := "git"
		publicKey, err := ssh.NewPublicKeysFromFile(publicKeyUser, a.PrivateKey, a.Password)

		if err != nil {
			return nil, err
		}
		return publicKey, nil
	} else if a.Username != "" && a.Password != "" {
		return &http.BasicAuth{
			Username: a.Username,
			Password: a.Password,
		}, nil
	} else {
		return nil, nil
	}
}
