package whatever

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"

	"github.com/mobroco/whatever/pkg/kind"
)

func (h Handler) GetAuthLogin(w http.ResponseWriter, r *http.Request) {
	var authenticator Authenticator
	providerURL := h.seed.Login.ProviderURL
	provider, err := oidc.NewProvider(r.Context(), providerURL)
	if err != nil {
		panic(err)
	}
	authenticator = Authenticator{
		Provider: provider,
		Config: oauth2.Config{
			ClientID:     h.seed.Login.ClientID,
			ClientSecret: h.seed.Login.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  h.seed.Login.RedirectURL,
			Scopes:       []string{oidc.ScopeOpenID, "profile"},
		},
	}

	state, err := generateRandomState()
	if err != nil {
		panic(err)
	}

	expire := time.Now().AddDate(0, 0, 1)
	http.SetCookie(w, &http.Cookie{
		Name:       "state",
		Value:      state,
		Path:       "/",
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     86400,
		Secure:     h.seed.Cookie.Secure,
		HttpOnly:   true,
		SameSite:   http.SameSiteLaxMode,
	})
	http.Redirect(w, r, authenticator.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (h Handler) GetAuthCallback(w http.ResponseWriter, r *http.Request) {
	var authenticator Authenticator
	providerURL := h.seed.Login.ProviderURL
	provider, err := oidc.NewProvider(r.Context(), providerURL)
	if err != nil {
		panic(err)
	}
	authenticator = Authenticator{
		Provider: provider,
		Config: oauth2.Config{
			ClientID:     h.seed.Login.ClientID,
			ClientSecret: h.seed.Login.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  h.seed.Login.RedirectURL,
			Scopes:       []string{oidc.ScopeOpenID, "profile"},
		},
	}

	stateCookie, err := r.Cookie("state")
	if err != nil {
		panic(err)
	}

	if stateCookie.Value != r.URL.Query().Get("state") {
		fmt.Println(cmp.Diff(stateCookie.Value, r.URL.Query().Get("state")))
		panic("state did not match")
	}

	token, err := authenticator.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		panic("failed to exchange an authorization code for a token - " + err.Error())
	}

	idToken, err := authenticator.VerifyIDToken(r.Context(), token)
	if err != nil {
		panic("failed to verify ID token - " + err.Error())
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		panic(err)
	}

	_, jwt, err := h.GetJWTAuth().Encode(profile)
	if err != nil {
		panic("failed to encode profile - " + err.Error())
	}

	var me kind.Me
	if myself, err := h.ygg.GrabMe(r.Context(), fmt.Sprint(profile["name"])); err == nil {
		me = myself
	}
	me.Profile = profile

	err = h.ygg.StashMe(r.Context(), fmt.Sprint(profile["name"]), me)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error getting me"))
		fmt.Println(err)
		return
	}

	expire := time.Now().AddDate(0, 0, 1)
	http.SetCookie(w, &http.Cookie{
		Name:       "jwt",
		Value:      jwt,
		Path:       "/",
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     86400,
		Secure:     h.seed.Cookie.Secure,
		HttpOnly:   true,
		SameSite:   http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
