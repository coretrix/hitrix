package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/coretrix/hitrix/service/component/config"
)

// RapidAPIInstagram28 https://rapidapi.com/yuananf/api/instagram28/
type RapidAPIInstagram28 struct {
	client       *http.Client
	apiKey       string
	apiHost      string
	apiBaseURL   string
	providerName string
}

func NewRapidAPIInstagram28(configService config.IConfig) (IProvider, error) {
	return &RapidAPIInstagram28{
		client:       &http.Client{},
		apiKey:       configService.MustString("instagram.api.rapid_api_token"),
		apiHost:      "instagram28.p.rapidapi.com",
		apiBaseURL:   "https://instagram28.p.rapidapi.com",
		providerName: "RapidAPIInstagram28",
	}, nil
}

func (i *RapidAPIInstagram28) GetName() string {
	return i.providerName
}

func (i *RapidAPIInstagram28) APIKey() string {
	return i.apiKey
}

func (i *RapidAPIInstagram28) GetAccount(account string) (*Account, error) {
	body, err := sendRapidRequest(
		i,
		fmt.Sprintf("%s/%s", i.apiBaseURL, fmt.Sprintf("user_info?user_name=%s", url.QueryEscape(account))),
	)

	if err != nil {
		return nil, err
	}

	response := struct {
		Status string
		Data   struct {
			User struct {
				ID       string
				Fullname string `json:"full_name"`
				Bio      string
				Website  string `json:"external_url"`
				Picture  string `json:"profile_pic_url_hd"`
				Posts    struct {
					Count int64 `json:"count"`
				} `json:"edge_owner_to_timeline_media"`
				Followers struct {
					Count int64 `json:"count"`
				} `json:"edge_followed_by"`
				Following struct {
					Count int64 `json:"count"`
				} `json:"edge_follow"`
				IsPrivate bool `json:"is_private"`
			}
		}
	}{}

	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, err
	}

	if response.Data.User.ID == "" {
		return nil, errors.New("internal_server_error status code 404")
	}

	userID, err := strconv.ParseInt(response.Data.User.ID, 10, 64)

	if err != nil {
		return nil, err
	}

	return &Account{
		AccountID: userID,
		FullName:  response.Data.User.Fullname,
		Bio:       response.Data.User.Bio,
		Posts:     response.Data.User.Posts.Count,
		Followers: response.Data.User.Followers.Count,
		Following: response.Data.User.Following.Count,
		Picture:   response.Data.User.Picture,
		IsPrivate: response.Data.User.IsPrivate,
		Website:   response.Data.User.Website,
	}, nil
}

func (i *RapidAPIInstagram28) GetFeed(accountID int64, nextPageToken string) ([]*Post, string, error) {
	reqNextPageToken := ""
	if len(nextPageToken) > 0 {
		reqNextPageToken = "&next_cursor=" + nextPageToken
	}
	body, err := sendRapidRequest(
		i,
		fmt.Sprintf("%s/%s", i.apiBaseURL, fmt.Sprintf("medias?user_id=%d&batch_size=40%s", accountID, reqNextPageToken)),
	)

	if err != nil {
		return nil, "", err
	}

	response := struct {
		Status string
		Data   struct {
			User struct {
				Timeline struct {
					PageInfo struct {
						HasNextPage bool   `json:"has_next_page"`
						EndCursor   string `json:"end_cursor"`
					} `json:"page_info"`

					Edges []struct {
						Node instagram40Post
					}
				} `json:"edge_owner_to_timeline_media"`
			}
		}
	}{}

	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, "", err
	}

	var posts []*Post

	for _, edge := range response.Data.User.Timeline.Edges {
		postData := edge.Node
		if postData.Type != "GraphSidecar" && postData.Type != "GraphImage" {
			continue
		}

		pi := postData.ToPost()
		posts = append(posts, pi)
	}

	nextPageToken = ""

	if response.Data.User.Timeline.PageInfo.HasNextPage {
		nextPageToken = response.Data.User.Timeline.PageInfo.EndCursor
	}

	return posts, nextPageToken, nil
}

type instagram40Post struct {
	ID                 string
	Type               string `json:"__typename"`
	DisplayURL         string `json:"display_url"`
	IsVideo            bool   `json:"is_video"`
	CreatedTime        int64  `json:"taken_at_timestamp"`
	EdgeMediaToCaption struct {
		Edges []struct {
			Node struct {
				Text string
			}
		}
	} `json:"edge_media_to_caption"`
	EdgeSidecarToChildren struct {
		Edges []struct {
			Node instagram40Post
		}
	} `json:"edge_sidecar_to_children"`
}

func (p instagram40Post) ToPost() *Post {
	post := &Post{
		ID:        p.ID,
		CreatedAt: p.CreatedTime,
	}

	// post caption
	if len(p.EdgeMediaToCaption.Edges) > 0 {
		post.Title = p.EdgeMediaToCaption.Edges[0].Node.Text
	}

	var images []string
	switch p.Type {
	case "GraphSidecar":
		for _, edge := range p.EdgeSidecarToChildren.Edges {
			if edge.Node.Type == "GraphImage" {
				images = append(images, edge.Node.DisplayURL)
			}
		}
	case "GraphImage":
		images = []string{
			p.DisplayURL,
		}
	}

	post.Images = images

	return post
}

func (i *RapidAPIInstagram28) APIHost() string {
	return i.apiHost
}

func (i *RapidAPIInstagram28) HTTPClient() *http.Client {
	return i.client
}
