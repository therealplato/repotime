package main

import (
	"context"
	"encoding/json"
	"fmt"
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
}
