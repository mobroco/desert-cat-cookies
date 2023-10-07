package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/hostrouter"
	heritage "github.com/mobroco/heritage/pkg"

	"github.com/mobroco/whatever/pkg"
)

var (
	//go:embed templates
	templates embed.FS
)

func main() {
	if len(os.Args) <= 1 {
		panic("too few arguments")
	}
	action := os.Args[1]
	flags := make(map[string]string)
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			argParts := strings.Split(arg, "=")
			key := strings.TrimPrefix(argParts[0], "--")

			var value string
			if len(argParts) > 1 {
				value = strings.Join(argParts[1:], "=")
			} else {
				value = "true"
			}
			flags[key] = value
		}
	}
	secretKey := os.Getenv("KEY")
	if flags["secret-key"] != "" {
		secretKey = flags["secret-key"]
	}
	var localMode bool
	if _, found := flags["local-mode"]; found {
		localMode = true
	}
	switch action {
	case "frontend":
		checkpoint := time.Now()
		for {
			out, err := exec.Command("npm", "run", "dev").CombinedOutput()
			if err != nil {
				panic(string(out) + err.Error())
			}
			fmt.Print("🌼")
			if now := time.Now(); now.Sub(checkpoint) > 5*time.Minute {
				checkpoint = now
				fmt.Println()
			}
		}
	case "sprout":
		enforceMain()

		packet := whatever.OpenPacket()

		_ = sprout()
		setup(packet)
	case "check":
		enforceMain()

		packet := whatever.OpenPacket()

		_ = sprout()
		setup(packet)

		diffBackend(packet)
	case "grow":
		enforceMain()

		var all bool
		if _, found := flags["all"]; found {
			all = true
		}
		packet := whatever.OpenPacket()
		t := sprout()
		setup(packet)
		if _, found := flags["build-backend"]; found || all {
			buildBackend(packet, t)
		}
		if _, found := flags["build-frontend"]; found || all {
			buildFrontend()
		}
		if _, found := flags["deploy-frontend"]; found || all {
			deployFrontend(packet)
		}
		if _, found := flags["deploy-backend"]; found || all {
			deployBackend(packet)
		}
		fmt.Println("done!")
	case "bump":
		enforceMain()

		out, err := exec.Command("git", "fetch", "--tags").CombinedOutput()
		if err != nil {
			panic(string(out) + err.Error())
		}

		out, err = exec.Command("git", "tag", "-l").Output()
		if err != nil {
			panic(string(out) + err.Error())
		}

		var versions []*semver.Version
		for _, tag := range strings.Split(string(out), "\n") {
			version, err := semver.NewVersion(tag)
			if err != nil || version == nil {
				continue
			}

			versions = append(versions, version)
		}
		sort.Sort(semver.Collection(versions))

		latest := versions[len(versions)-1]
		var bumpedVersion semver.Version
		fmt.Printf("latest version is [%s]\n", latest)
		fmt.Print("what kind of bump? (major/minor/patch) ")
		var bump string
		_, _ = fmt.Scanln(&bump)

		switch bump {
		case "major":
			bumpedVersion = latest.IncMajor()
		case "minor":
			bumpedVersion = latest.IncMinor()
		case "patch":
			bumpedVersion = latest.IncPatch()
		default:
			panic("unknown kind of bump")
		}

		fmt.Printf("tag commit as [%s]? (y/n) ", bumpedVersion.String())
		var confirm string
		_, _ = fmt.Scanln(&confirm)
		if !strings.EqualFold(confirm, "y") {
			panic("did not confirm to run")
		}

		out, err = exec.Command("git", "tag", "v"+bumpedVersion.String()).CombinedOutput()
		if err != nil {
			panic(string(out) + err.Error())
		}

		out, err = exec.Command("git", "push", "origin", "main", "--tags").CombinedOutput()
		if err != nil {
			panic(string(out) + err.Error())
		}
	case "encrypt":
		fmt.Println(whatever.Encrypt(os.Args[2], secretKey))
	case "decrypt":
		fmt.Println(whatever.Decrypt(os.Args[2], secretKey))
	case "serve":
		var packet heritage.Packet
		if local := flags["local"]; local == "" && action != "local-backend" {
			packet = whatever.OpenPacket()
		} else {
			packet = whatever.OpenLocalPacket(strings.Split(local, ",")...)
		}
		if _, found := flags["peek"]; found {
			whatever.PeekAt(packet)
		}
		r := chi.NewRouter()
		r.Use(middleware.RequestID)
		r.Use(middleware.Recoverer)
		r.Use(middleware.StripSlashes)
		r.Use(middleware.Compress(5))
		r.Use(middleware.Timeout(60 * time.Second))

		hr := hostrouter.New()

		for _, seed := range packet.Seeds {
			switch seed.Name {
			// Live
			case "desert-cat-cookies":
				h := whatever.NewHandler(templates, seed, secretKey, localMode)

				router := chi.NewRouter()
				router.Use(middleware.Logger)
				// backend
				router.Post("/x/estimates", h.CreateEstimate)
				// frontend
				router.Get("/dist/bundle.js", h.GetBundleJS)
				router.Get("/public/*", h.GetPublic)
				router.Get("/", h.GetIndex)
				// support
				router.Post("/report/csp", h.ReceiveContentSecurityPolicyReport)
				router.Post("/health", h.GetHealth)
				// not found
				router.Get("/*", h.Get404)
				hr.Map(seed.FQDN, router)

			case "greasy-shadows":
				h := whatever.NewHandler(templates, seed, secretKey, localMode)
				router := chi.NewRouter()
				router.Use(middleware.Logger)
				// frontend
				router.Get("/dist/bundle.js", h.GetBundleJS)
				router.Get("/public/*", h.GetPublic)
				router.Get("/", h.GetIndex)
				// support
				router.Post("/report/csp", h.ReceiveContentSecurityPolicyReport)
				router.Post("/health", h.GetHealth)
				// not found
				router.Get("/*", h.Get404)
				hr.Map(seed.FQDN, router)
			}

			// everything else
			routes := chi.NewRouter()
			routes.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("OK"))
			})
			routes.Get("/*", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("..."))
			})
			hr.Map("*", routes)
			r.Mount("/", hr)
			log.Println(http.ListenAndServe(":3000", r))
		}
	}
}

func sprout() string {
	now := time.Now().Format("20060102-150405")
	fmt.Print("🌱 sprout ", now, " ")
	_ = os.WriteFile("sprout.txt", []byte(now), os.ModePerm)
	fmt.Println("🌱")
	return now
}

func setup(packet heritage.Packet) {
	out, err := exec.Command("aws", "ecr", "get-login-password", "--region", packet.Region).CombinedOutput()
	if err != nil {
		panic(string(out) + err.Error())
	}
	fmt.Print("☀️")
	out, err = exec.Command("docker", "login", "--username", "AWS", "--password", string(out), packet.Account+".dkr.ecr."+packet.Region+".amazonaws.com").CombinedOutput()
	if err != nil {
		panic(string(out) + err.Error())
	}
	fmt.Println("☀️")
}

func enforceMain() {
	out, err := exec.Command("git", "symbolic-ref", "HEAD").CombinedOutput()
	if err != nil {
		panic(string(out) + err.Error())
	}
	if strings.TrimSpace(string(out)) != "refs/heads/main" {
		panic("this can only be run on the main branch")
	}
}

func buildBackend(packet heritage.Packet, t string) {
	fmt.Print("🍄 build the backend ")
	out, err := exec.Command("docker", "build", "-t", packet.Name+":local", ".").CombinedOutput()
	if err != nil {
		panic(string(out) + err.Error())
	}
	fmt.Print("🍄")

	for _, tagPart := range []string{"latest", t} {
		tag := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s", packet.Account, packet.Region, packet.Name, tagPart)
		out, err = exec.Command("docker", "tag", packet.Name+":local", tag).CombinedOutput()
		if err != nil {
			panic(string(out) + err.Error())
		}
		fmt.Print("🍄")
	}

	for _, tagPart := range []string{"latest", t} {
		tag := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s", packet.Account, packet.Region, packet.Name, tagPart)
		out, err = exec.Command("docker", "push", tag).CombinedOutput()
		if err != nil {
			panic(string(out) + err.Error())
		}
		fmt.Print("🍄")
	}
	fmt.Println("🍄")
}

func deployBackend(packet heritage.Packet) {
	fmt.Print("🍄 deploy the backend ")
	cmd := exec.Command("cdk", "deploy", packet.Name+"-*")
	cmd.Dir = "../garden"
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(string(out) + err.Error())
	}
	fmt.Println("🍄")
}

func diffBackend(packet heritage.Packet) {
	fmt.Println("🍄 diff the backend ")
	cmd := exec.Command("cdk", "diff", packet.Name+"-*")
	cmd.Dir = "../garden"
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(string(out) + err.Error())
	}
	fmt.Println(string(out))
	fmt.Println("🍄")
}

func buildFrontend() {
	fmt.Print("🌼 build the frontend ")
	out, err := exec.Command("npm", "run", "build").CombinedOutput()
	if err != nil {
		panic(string(out) + err.Error())
	}
	fmt.Println("🌼")
}

func deployFrontend(packet heritage.Packet) {
	fmt.Print("🌼 deploy the frontend ")
	for _, seed := range packet.Seeds {
		if seed.CDN.Bucket == "" {
			continue
		}
		fmt.Print("🌼")
		out, err := exec.Command("aws", "s3", "sync", "--delete", "dist", fmt.Sprintf("s3://%s/dist", seed.CDN.Bucket)).CombinedOutput()
		if err != nil {
			panic(string(out) + err.Error())
		}

		fmt.Print("🌼")
		out, err = exec.Command("aws", "s3", "sync", "--delete", "public", fmt.Sprintf("s3://%s/public", seed.CDN.Bucket)).CombinedOutput()
		if err != nil {
			panic(string(out) + err.Error())
		}

		fmt.Print("🌼")
		out, err = exec.Command("aws", "s3", "sync", "--delete", "public", fmt.Sprintf("s3://%s/public", seed.CDN.Bucket)).CombinedOutput()
		if err != nil {
			panic(string(out) + err.Error())
		}

		fmt.Print("🌼")
		out, err = exec.Command("aws", "cloudfront", "create-invalidation", "--distribution-id=E30072L8WN05EW", "--paths=/*").CombinedOutput()
		if err != nil {
			panic(string(out) + err.Error())
		}
	}
	fmt.Println("🌼")
}
