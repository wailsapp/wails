package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

// Sponsor is a single sponsorship, flattened from the GraphQL response.
type Sponsor struct {
	Login     string
	Name      string
	AvatarURL string
	URL       string
	// Monthly is the monthly dollar amount. Zero when the tier is not
	// visible to the querying token or the sponsorship is one-time.
	Monthly   int
	IsOneTime bool
}

// DisplayName returns the sponsor's name, falling back to their login.
func (s Sponsor) DisplayName() string {
	if s.Name != "" {
		return s.Name
	}
	return s.Login
}

const sponsorQuery = `query($login: String!, $after: String) {
  user(login: $login) {
    sponsorshipsAsMaintainer(first: 100, after: $after, activeOnly: true) {
      pageInfo { hasNextPage endCursor }
      nodes {
        tier { monthlyPriceInDollars isOneTime }
        sponsorEntity {
          ... on User { login name avatarUrl url }
          ... on Organization { login name avatarUrl url }
        }
      }
    }
  }
}`

type graphQLResponse struct {
	Data struct {
		User struct {
			SponsorshipsAsMaintainer struct {
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
				Nodes []struct {
					Tier *struct {
						MonthlyPriceInDollars int  `json:"monthlyPriceInDollars"`
						IsOneTime             bool `json:"isOneTime"`
					} `json:"tier"`
					SponsorEntity struct {
						Login     string `json:"login"`
						Name      string `json:"name"`
						AvatarURL string `json:"avatarUrl"`
						URL       string `json:"url"`
					} `json:"sponsorEntity"`
				} `json:"nodes"`
			} `json:"sponsorshipsAsMaintainer"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// FetchSponsors returns all active sponsors of login, sorted by monthly
// amount (descending) then login for deterministic output.
func FetchSponsors(client *http.Client, token, login string) ([]Sponsor, error) {
	var sponsors []Sponsor
	seen := map[string]bool{}
	after := ""

	for {
		variables := map[string]any{"login": login}
		if after != "" {
			variables["after"] = after
		}
		body, err := json.Marshal(map[string]any{
			"query":     sponsorQuery,
			"variables": variables,
		})
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, "https://api.github.com/graphql", bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "wails-sponsorkit")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("github graphql: %s: %s", resp.Status, truncate(string(data), 300))
		}

		var parsed graphQLResponse
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("parsing graphql response: %w", err)
		}
		if len(parsed.Errors) > 0 {
			return nil, fmt.Errorf("github graphql: %s", parsed.Errors[0].Message)
		}

		conn := parsed.Data.User.SponsorshipsAsMaintainer
		for _, node := range conn.Nodes {
			entity := node.SponsorEntity
			if entity.Login == "" || seen[entity.Login] {
				continue
			}
			seen[entity.Login] = true
			s := Sponsor{
				Login:     entity.Login,
				Name:      entity.Name,
				AvatarURL: entity.AvatarURL,
				URL:       entity.URL,
			}
			if node.Tier != nil {
				s.IsOneTime = node.Tier.IsOneTime
				if !node.Tier.IsOneTime {
					s.Monthly = node.Tier.MonthlyPriceInDollars
				}
			}
			sponsors = append(sponsors, s)
		}

		if !conn.PageInfo.HasNextPage {
			break
		}
		after = conn.PageInfo.EndCursor
		time.Sleep(100 * time.Millisecond)
	}

	sort.Slice(sponsors, func(i, j int) bool {
		if sponsors[i].Monthly != sponsors[j].Monthly {
			return sponsors[i].Monthly > sponsors[j].Monthly
		}
		return sponsors[i].Login < sponsors[j].Login
	})
	return sponsors, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
