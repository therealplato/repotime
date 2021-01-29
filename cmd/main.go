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
	apiServer := &apiServer{client: client}
	mux := http.NewServeMux()
	mux.Handle("/api/", apiServer)
	mux.Handle("/", http.FileServer(http.Dir("./public_http")))
	fmt.Println("It's repotime! Go to http://localhost:2992")
	http.ListenAndServe("0.0.0.0:2992", mux)
}

type apiServer struct {
	client           *github.Client
	mux              *http.ServeMux
	ownerLogin       string
	repoName         string
	chosenRepository *github.Repository
}

func (s *apiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.mux == nil {
		s.mux = http.NewServeMux()
		s.mux.HandleFunc("/api/username", s.getUsername)
		s.mux.HandleFunc("/api/repositories", s.getRepos)
		s.mux.HandleFunc("/api/commits", s.getCommits)
		s.mux.HandleFunc("/api/issues", s.getIssues)
		s.mux.HandleFunc("/api/set-repository", s.setRepo)
	}
	s.mux.ServeHTTP(w, r)
}

func (s *apiServer) getUsername(w http.ResponseWriter, r *http.Request) {
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
	repo := &github.Repository{}
	err := json.NewDecoder(r.Body).Decode(repo)
	if err != nil {
		fmt.Printf("setRepo decode error: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.chosenRepository = repo
	s.ownerLogin = *repo.GetOwner().Login
	s.repoName = *repo.Name
	fmt.Printf("repo has been set to %v\n", *repo.FullName)
}

func (s *apiServer) getCommits(w http.ResponseWriter, r *http.Request) {
	if s.chosenRepository == nil {
		fmt.Println("cannot get commits until repository is chosen")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cc, _, err := s.client.Repositories.ListCommits(r.Context(), s.ownerLogin, s.repoName, nil)
	if err != nil {
		fmt.Printf("failed to retrieve commits: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bb, err := json.MarshalIndent(cc, "", "  ")
	if err != nil {
		fmt.Printf("getCommits marshal error: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bb)
}

func (s *apiServer) getIssues(w http.ResponseWriter, r *http.Request) {
	if s.chosenRepository == nil {
		fmt.Println("cannot get issues until repository is chosen")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ii, _, err := s.client.Activity.ListIssueEventsForRepository(r.Context(), s.ownerLogin, s.repoName, nil)
	if err != nil {
		fmt.Printf("failed to retrieve issues: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bb, err := json.MarshalIndent(ii, "", "  ")
	if err != nil {
		fmt.Printf("getIssues marshal error: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bb)
}
