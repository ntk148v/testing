package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	// simple query starred repository
	// var query struct {
	// 	Viewer struct {
	// 		Login               githubv4.String
	// 		StarredRepositories struct {
	// 			Nodes []struct {
	// 				ID              githubv4.String
	// 				Name            githubv4.String
	// 				PrimaryLanguage struct {
	// 					ID    githubv4.String
	// 					Name  githubv4.String
	// 					Color githubv4.String
	// 				}
	// 			}
	// 			PageInfo struct {
	// 				EndCursor   githubv4.String
	// 				HasNextPage githubv4.Boolean
	// 			}
	// 		} `graphql:"starredRepositories(first:3)"`
	// 	}
	// }

	items, err := getStarredRepositoriesByUsername(ctx, client, "ntk148v", "")
	if err != nil {
		panic(err)
	}

	fmt.Println(items[0])
	printJSON(items)
}

type Repository struct {
	Name           string   `json:"name"`
	URL            string   `json:"url"`
	Description    string   `json:"description"`
	StargazerCount int      `json:"stargazer_count"`
	Language       string   `json:"language"`
	Topics         []string `json:"topics"`
}

func getStarredRepositoriesByUsername(ctx context.Context, client *githubv4.Client, username, after string) ([]Repository, error) {
	var query struct {
		User struct {
			StarredRepositories struct {
				TotalCount githubv4.Int
				Nodes      []struct {
					Name            githubv4.String
					NameWithOwner   githubv4.String
					Description     githubv4.String
					URL             githubv4.String
					StargazerCount  githubv4.Int
					ForkCount       githubv4.Int
					IsPrivate       githubv4.Boolean
					PushedAt        githubv4.DateTime
					UpdatedAt       githubv4.DateTime
					PrimaryLanguage struct {
						ID   githubv4.String
						Name githubv4.String
					}
					RepositoryTopics struct {
						Nodes []struct {
							Topic struct {
								Name           githubv4.String
								StargazerCount githubv4.Int
							}
						}
					} `graphql:"repositoryTopics(first: 100)"`
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage githubv4.Boolean
				}
			} `graphql:"starredRepositories(first: 100, after: $after, orderBy: {direction: DESC, field: STARRED_AT})"`
		} `graphql:"user(login:$username)"`
	}

	var items []Repository

	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"after":    githubv4.String(after),
	}
	if err := client.Query(ctx, &query, variables); err != nil {
		return items, err
	}

	hasNext := query.User.StarredRepositories.PageInfo.HasNextPage
	endCursor := query.User.StarredRepositories.PageInfo.EndCursor
	repos := query.User.StarredRepositories.Nodes

	for _, repo := range repos {
		tmpRepo := Repository{
			Name:           string(repo.Name),
			Description:    string(repo.Description),
			Language:       string(repo.PrimaryLanguage.Name),
			URL:            string(repo.URL),
			StargazerCount: int(repo.StargazerCount),
			Topics:         make([]string, 0),
		}

		for _, topic := range repo.RepositoryTopics.Nodes {
			tmpRepo.Topics = append(tmpRepo.Topics, string(topic.Topic.Name))
		}

		items = append(items, tmpRepo)
	}

	if hasNext {
		tmpItems, err := getStarredRepositoriesByUsername(ctx, client, username, string(endCursor))
		if err != nil {
			return items, err
		}
		items = append(items, tmpItems...)
	}

	return items, nil
}

// printJSON prints v as JSON encoded with indent to stdout. It panics on any error.
func printJSON(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "\t")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
