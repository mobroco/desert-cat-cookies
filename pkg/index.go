package whatever

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

func (h Handler) SecureGetIndex(w http.ResponseWriter, r *http.Request) {
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	token, err := h.GetJWTAuth().Decode(jwtCookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	me, err := h.ygg.GrabMe(r.Context(), fmt.Sprint(token.PrivateClaims()["name"]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error getting me"))
		fmt.Println(err)
		return
	}

	var extraBodyClasses string
	parts := strings.Split(r.RequestURI, "/")
	if len(parts) >= 2 {
		// /w/<main>/<sub>
		if parts[1] == "w" {
			if !me.Allowed(parts[2]) {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			main := parts[2]
			var sub string
			if len(parts) > 3 {
				sub = parts[3]
			}

			what, err := h.ygg.GrabWhat(r.Context(), main, sub)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("error getting what"))
				fmt.Println(err)
				return
			}
			extraBodyClasses += "mt_" + what.Kind
		}
	}

	nonce := makeNonce()
	if csp := h.seed.ContentSecurityPolicy; csp.Make(nonce) != "" {
		if csp.ReportOnly {
			w.Header().Set("Content-Security-Policy-Report-Only", csp.Make(nonce))
		} else {
			w.Header().Set("Content-Security-Policy", csp.Make(nonce))
		}
	}

	id := strings.TrimPrefix(strings.Join(parts, "-"), "-")
	if id == "" {
		id = "index"
	}
	site := h.seed.Site
	site.ID = id
	site.BodyClasses += " " + extraBodyClasses
	err = h.index.Execute(w, map[string]any{
		"name":  h.seed.Name,
		"site":  site,
		"nonce": nonce,
		"cdn":   h.seed.CDN.Use,
		"me":    me,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (h Handler) GetIndex(w http.ResponseWriter, r *http.Request) {
	nonce := makeNonce()
	if csp := h.seed.ContentSecurityPolicy; csp.Make(nonce) != "" {
		if csp.ReportOnly {
			w.Header().Set("Content-Security-Policy-Report-Only", csp.Make(nonce))
		} else {
			w.Header().Set("Content-Security-Policy", csp.Make(nonce))
		}
	}

	parts := strings.Split(r.RequestURI, "/")
	id := strings.TrimPrefix(strings.Join(parts, "-"), "-")
	if id == "" {
		id = "index"
	}
	site := h.seed.Site
	site.ID = id
	err := h.index.Execute(w, map[string]any{
		"name":  h.seed.Name,
		"site":  site,
		"nonce": nonce,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (h Handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	nonce := makeNonce()
	if csp := h.seed.ContentSecurityPolicy; csp.Make(nonce) != "" {
		if csp.ReportOnly {
			w.Header().Set("Content-Security-Policy-Report-Only", csp.Make(nonce))
		} else {
			w.Header().Set("Content-Security-Policy", csp.Make(nonce))
		}
	}
	h.seed.Site.HTMLClasses += " h-full"
	h.seed.Site.BodyClasses += " h-full"

	parts := strings.Split(r.RequestURI, "/")
	site := h.seed.Site
	site.ID = strings.TrimPrefix(strings.Join(parts, "-"), "-")
	err := h.index.Execute(w, map[string]any{
		"name":  h.seed.Name,
		"site":  site,
		"nonce": nonce,
		"cdn":   h.seed.CDN.Use,
	})
	if err != nil {
		log.Fatal(err)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func makeNonce() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
