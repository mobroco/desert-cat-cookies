package whatever

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"html/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/jwtauth/v5"
	heritage "github.com/mobroco/heritage/pkg"
	stripeclient "github.com/stripe/stripe-go/v74/client"
	mail "github.com/xhit/go-simple-mail/v2"
	"golang.org/x/oauth2"
)

type Handler struct {
	seed       heritage.Seed
	awsConfig  aws.Config
	smtpServer *mail.SMTPServer
	index      *template.Template
	tokenAuth  *jwtauth.JWTAuth
	localMode  bool
}

type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}
	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

func NewHandler(fs embed.FS, seed heritage.Seed, secretKey string, localMode bool) Handler {
	var smtpServer *mail.SMTPServer
	if seed.SMTP.Host != "" {
		smtpServer = mail.NewSMTPClient()
		smtpServer.Host = seed.SMTP.Host
		smtpServer.Port = seed.SMTP.Port
		smtpServer.Username = Decrypt(seed.SMTP.Username, secretKey)
		smtpServer.Password = Decrypt(seed.SMTP.Password, secretKey)
		smtpServer.Encryption = mail.EncryptionSTARTTLS
		smtpServer.KeepAlive = false
		smtpServer.ConnectTimeout = 10 * time.Second
		smtpServer.SendTimeout = 10 * time.Second
		smtpServer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	awsConfig, _ := config.LoadDefaultConfig(context.TODO())
	index := template.Must(template.ParseFS(fs, "templates/index.html"))
	seed.Login.ClientID = Decrypt(seed.Login.ClientID, secretKey)
	seed.Login.ClientSecret = Decrypt(seed.Login.ClientSecret, secretKey)

	stripeClient := &stripeclient.API{}
	if seed.Checkout.Stripe.Enabled {
		stripeSecretKey := Decrypt(seed.Checkout.Stripe.SecretKey, secretKey)
		stripeClient.Init(stripeSecretKey, nil)
	}
	return Handler{
		seed:       seed,
		smtpServer: smtpServer,
		awsConfig:  awsConfig,
		index:      index,
		tokenAuth:  jwtauth.New("HS256", []byte(seed.JWT.Secret), nil),
		localMode:  localMode,
	}
}

func (h Handler) GetJWTAuth() *jwtauth.JWTAuth {
	return h.tokenAuth
}
