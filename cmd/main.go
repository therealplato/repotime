package main

import (
	"context"
	"encoding/json"
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
	client           *github.Client
	mux              *http.ServeMux
	chosenRepository *github.Repository
}

func (s *apiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.mux == nil {
		s.mux = http.NewServeMux()
		s.mux.HandleFunc("/api/username", s.getUsername)
		s.mux.HandleFunc("/api/repositories", s.getRepos)
		s.mux.HandleFunc("/api/set-repository", s.setRepo)
	}
	s.mux.ServeHTTP(w, r)
}

func (s *apiServer) getUsername(w http.ResponseWriter, r *http.Request) {
	fmt.Println("retrieving username")
	u, _, err := s.client.Users.Get(r.Context(), "")
	if err != nil {
		fmt.Printf("getUsername error: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bb, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		fmt.Printf("getUsername marshal error: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bb)
}

func (s *apiServer) getRepos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("retrieving repos")
	rr, _, err := s.client.Repositories.List(r.Context(), "", nil)
	if err != nil {
		fmt.Printf("getRepos error: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bb, err := json.MarshalIndent(rr, "", "  ")
	if err != nil {
		fmt.Printf("getRepos marshal error: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bb)
}

func (s *apiServer) setRepo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("setting repo")
	repo := &github.Repository{}
	err := json.NewDecoder(r.Body).Decode(repo)
	if err != nil {
		fmt.Printf("setRepo decode error: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.chosenRepository = repo
	fmt.Printf("repo has been set to %v\n", *repo.FullName)
}
