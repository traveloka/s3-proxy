package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/abihf/delta"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/traveloka/s3-proxy/auth"
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	client := ssm.NewFromConfig(cfg)
	res, err := client.GetParameters(context.Background(), &ssm.GetParametersInput{
		Names:          []string{"/tvlk-secret/webghr/web/github.secret", "/tvlk-secret/webghr/web/oauth.secret"},
		WithDecryption: true,
	})
	if err != nil {
		fmt.Printf("can read parameter store: %v", err)
		return
	}

	os.Setenv("GITHUB_CLIENT_SECRET", *res.Parameters[0].Value)
	os.Setenv("OAUTH_SECRET", *res.Parameters[1].Value)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/_healthz", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "ready")
	})

	mux.HandleFunc("/_auth", auth.HandleAuthCallback)
	mux.HandleFunc("/_logout", auth.HandleLogout)
	mux.Handle("/", auth.WithAuth(http.HandlerFunc(handleProxy)))
	port := "8080"
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}
	fmt.Printf("listening on http://127.0.0.0:%s/", port)
	panic(delta.ServeOrStartLambda(":"+port, mux))
}
