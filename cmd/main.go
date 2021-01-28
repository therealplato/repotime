package main

import (
	"os"

	"github.com/therealplato/repotime/githubauth"
)

func main() {
	githubauth.MustAuthorize(os.Stdout)
}
