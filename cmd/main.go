package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/therealplato/repotime/githubauth"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	token := githubauth.MustAuthorize(os.Stdout)
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	client := github.NewClient(tc)
	_ = client

	/*
		repos, _, err := client.Repositories.List(ctx, "", nil)
		if err != nil {
			fmt.Print("failed to perform github repository list: %q\n", err)
			os.Exit(1)
		}
		bb, err := json.MarshalIndent(repos, "", "  ")
		if err != nil {
			fmt.Print("failed to marshal github repository list: %q\n", err)
			os.Exit(1)
		}
		fmt.Print(string(bb))
	*/
	apiServer := &apiServer{client: client}
	mux := http.NewServeMux()
	mux.Handle("/api/", apiServer)
	mux.Handle("/", http.FileServer(http.Dir("./public_http")))
	http.ListenAndServe("0.0.0.0:2992", mux)
}

type apiServer struct {
	client *github.Client
	mux    *http.ServeMux
}

func (s *apiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.mux == nil {
		s.mux = http.NewServeMux()
		s.mux.HandleFunc("/api/username", s.getUsername)
	}
	s.mux.ServeHTTP(w, r)
}

func (s *apiServer) getUsername(w http.ResponseWriter, r *http.Request) {
	fmt.Println("retrieving username")
	w.Write([]byte(`{"username":"@therealplato2"}`))
}
