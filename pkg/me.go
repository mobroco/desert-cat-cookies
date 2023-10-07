package whatever

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mobroco/whatever/pkg/kind"
)

func (h Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	var tokenRaw string
	if bearerToken := r.Header.Get("Authorization"); strings.HasPrefix(bearerToken, "Bearer ") {
		tokenRaw = strings.TrimSpace(strings.TrimPrefix(bearerToken, "Bearer "))
	}
	if jwtCookie, err := r.Cookie("jwt"); err == nil {
		tokenRaw = jwtCookie.Value
	}
	me := kind.Me{
		Profile: map[string]any{
			"name":  "What?",
			"email": "nothing@whatever.dev",
		},
		Allows: nil,
	}
	if token, err := h.GetJWTAuth().Decode(tokenRaw); err == nil {
		if myself, err := h.ygg.GrabMe(r.Context(), fmt.Sprint(token.PrivateClaims()["name"])); err == nil {
			me = myself
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(me.JSON())
}
