package sessions

import "golang.org/x/oauth2"

type Store interface {
	Save(accessToken string, token *oauth2.Token) error

	Get(name string) (*oauth2.Token, error)

	Delete(name string) error
}
