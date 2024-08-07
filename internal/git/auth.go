package git

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Auth struct {
	PrivateKeyPath string
	Username       string
	Password       string
}

func NewAuthMethod(auth Auth) (transport.AuthMethod, error) {
	if auth.PrivateKeyPath != "" {
		_, err := os.Stat(auth.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
		publicKeys, err := ssh.NewPublicKeysFromFile("git", auth.PrivateKeyPath, auth.Password)
		if err != nil {
			return nil, err
		}
		return publicKeys, nil
	} else if auth.Username != "" && auth.Password != "" {
		return &http.BasicAuth{
			Username: auth.Username,
			Password: auth.Password,
		}, nil
	}
	return nil, nil
}
