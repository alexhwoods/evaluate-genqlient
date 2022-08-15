package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/alexhwoods/evaluate-genqlient/generated"
)

// go:generate go run github.com/Khan/genqlient genqlient.yaml

const GithubGraphqlAPI = "https://api.github.com/graphql"

type authedTransport struct {
	key     string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.String() == GithubGraphqlAPI {
		req.Header.Set("Authorization", "Bearer "+t.key)
	}
	return t.wrapped.RoundTrip(req)
}

func main() {
	key := "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"

	ctx := context.Background()

	httpClient := http.Client{
		Transport: &authedTransport{
			key:     key,
			wrapped: http.DefaultTransport,
		},
	}
	client := graphql.NewClient(GithubGraphqlAPI, &httpClient)
	resp, err := generated.GetRepository(ctx, client, "kubernetes", "kubernetes")

	if err != nil {
		fmt.Printf("error: %s", err)
	}

	useRepository(resp.Repository)
}

func useRepository(repository generated.GetRepositoryRepository) {
	fmt.Printf("repository.Name: %v\n", repository.Name)
	fmt.Printf("repository.Description: %v\n", repository.Description)
	fmt.Printf("repository.StargazerCount: %v\n", repository.StargazerCount)
	fmt.Printf("repository.Url: %v\n", repository.Url)
}
