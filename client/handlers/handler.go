package handlers

import (
	"net/http"

	types "github.com/cynico/oauth/client/types"

	"github.com/hashicorp/go-hclog"
)

type CallbackHandler struct {
	log   hclog.Logger
	oauth *types.OAuth
}

func NewCallbackHandler(l hclog.Logger, o *types.OAuth) *CallbackHandler {
	return &CallbackHandler{
		log:   l,
		oauth: o,
	}
}

func (c *CallbackHandler) Handle(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state != c.oauth.State {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	c.log.Debug("got callback request", "code", code, "state", state)

	// tok, err := c.oauth.Conf.Exchange(r.Context(), code)
	// if err != nil {
	// 	http.Error(w, "couldn't get authorized", http.StatusBadGateway)
	// 	return
	// }
	// c.oauth.Token = tok

}
