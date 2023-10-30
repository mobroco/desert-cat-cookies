package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/go-chi/chi/v5"
)

type Estimate struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	CookieTheme    string `json:"cookie_theme"`
	CookieQuantity string `json:"cookie_quantity"`
	PickupDate     string `json:"pickup_date"`
	AnythingElse   string `json:"anything_else"`
}

func Execute() error {
	sendGrid := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	from := mail.NewEmail("Clarice", "clarice@em2928.desertcatcookies.com")
	to := mail.NewEmail("Stacy", "stacymulhern@gmail.com")

	r := chi.NewRouter()
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	FileServer(r, "/public", filesDir)

	r.Get("/dist/bundle.js", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("dist/bundle.js")
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(file)
	})

	r.Post("/estimates", func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var estimate Estimate
		_ = json.Unmarshal(raw, &estimate)

		var sb strings.Builder
		sb.WriteString("First Name: " + estimate.FirstName + "\n")
		sb.WriteString("Last Name: " + estimate.LastName + "\n")
		sb.WriteString("Email: " + estimate.Email + "\n")
		sb.WriteString("Phone Number: " + estimate.PhoneNumber + "\n")
		sb.WriteString("Cookie Theme: " + estimate.CookieTheme + "\n")
		sb.WriteString("Cookie Quantity: " + estimate.CookieQuantity + "\n")
		sb.WriteString("Pickup Date: " + estimate.PickupDate + "\n")
		sb.WriteString("Anything Else: " + estimate.AnythingElse + "\n")
		message := mail.NewSingleEmail(from, "Estimate Request", to, sb.String(), "")
		response, err := sendGrid.Send(message)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println("email sent", response.StatusCode)
		}
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("public/index.html")
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(file)
	})
	return http.ListenAndServe(":3000", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
