package types

import "golang.org/x/oauth2"

type OAuth struct {
	Conf  *oauth2.Config
	State string
	Token *oauth2.Token
}
