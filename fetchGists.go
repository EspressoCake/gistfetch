package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/google/go-github/v41/github"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"os"
)

func checkExistence(desiredLocation string) bool {
	_, err := os.Stat(desiredLocation)
	if errors.Is(err, os.ErrNotExist) {
		return true
	}
	return false
}

func writeDataToFile(gistArray []github.GistFile) error {
	bar := progressbar.Default(int64(len(gistArray)))

	for _, gist := range gistArray {
		newfile := *gist.Filename
		currentURL := *gist.RawURL

		path, err := os.Getwd()
		if err != nil {
			log.Panic(err)
		}

		currentFullPath := fmt.Sprintf("%s/%s", path, newfile)

		if !checkExistence(currentFullPath) {
			currentFullPath = fmt.Sprintf("%s/%04d_%s", path, *gist.Size, newfile)
		}

		file, err := os.Create(currentFullPath)
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()

		resp, err := http.Get(currentURL)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		file.Write(body)
		bar.Add(1)
	}

	return nil
}

func main() {
	ghAccessToken := flag.String("token", "none", "Github API access token")
	flag.Parse()

	repo := &github.GistListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *ghAccessToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var allGists []github.GistFile

	for {
		gists, resp, err := client.Gists.List(ctx, "", repo)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range gists {
			for _, gistinfo := range entry.Files {
				allGists = append(allGists, gistinfo)
			}
		}

		if resp.NextPage == 0 {
			break
		}

		repo.Page = resp.NextPage
	}

	writeDataToFile(allGists)
}
