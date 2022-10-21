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

// RapidAPIInstagram188 https://rapidapi.com/ahmetarpaci/api/instagram188/
type RapidAPIInstagram188 struct {
	client       *http.Client
	apiKey       string
	apiHost      string
	apiBaseURL   string
	providerName string
}

func NewRapidAPIInstagram188(configService config.IConfig) (IProvider, error) {
	return &RapidAPIInstagram188{
		client:       &http.Client{},
		apiKey:       configService.MustString("instagram.api.rapid_api_token"),
		apiHost:      "instagram188.p.rapidapi.com",
		apiBaseURL:   "https://instagram188.p.rapidapi.com",
		providerName: "RapidAPIInstagram188",
	}, nil
}

func (i *RapidAPIInstagram188) GetName() string {
	return i.providerName
}

func (i *RapidAPIInstagram188) APIKey() string {
	return i.apiKey
}

func (i *RapidAPIInstagram188) GetAccount(account string) (*Account, error) {
	body, err := sendRapidRequest(
		i,
		fmt.Sprintf("%s/%s/%s", i.apiBaseURL, "userinfo", url.PathEscape(account)),
	)

	if err != nil {
		return nil, err
	}

	response := struct {
		Success bool `json:"success"`
		Data    struct {
			ID            string `json:"id"`
			IsPrivate     bool   `json:"is_private"`
			FullName      string `json:"full_name"`
			Biography     string `json:"biography"`
			ProfilePicURL string `json:"profile_pic_url_hd"`
			Following     struct {
				Count int64 `json:"count"`
			} `json:"edge_follow"`
			Followers struct {
				Count int64 `json:"count"`
			} `json:"edge_followed_by"`
			Timeline struct {
				Count int64 `json:"count"`
			} `json:"edge_owner_to_timeline_media"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal %s : %+v", string(body), err)
	}

	if !response.Success {
		return nil, errors.New("internal_server_error")
	}

	if response.Data.ID == "" {
		return nil, errors.New("internal_server_error status code 404")
	}

	userID, err := strconv.ParseInt(response.Data.ID, 10, 64)

	if err != nil {
		return nil, err
	}

	return &Account{
		AccountID: userID,
		FullName:  response.Data.FullName,
		Bio:       response.Data.Biography,
		Posts:     response.Data.Timeline.Count,
		Followers: response.Data.Followers.Count,
		Following: response.Data.Following.Count,
		Picture:   response.Data.ProfilePicURL,
		IsPrivate: response.Data.IsPrivate,
		Website:   response.Data.Biography,
		BioLinks:  []*BioLink{{URL: response.Data.Biography}},
	}, nil
}

func (i *RapidAPIInstagram188) GetFeed(accountID int64, nextPageToken string) ([]*Post, string, error) {
	nextCursor := "{end_cursor}"

	if len(nextPageToken) > 0 {
		nextCursor = nextPageToken
	}

	body, err := sendRapidRequest(
		i,
		fmt.Sprintf("%s/%s/%d/40/%s", i.apiBaseURL, "userpost", accountID, url.PathEscape(nextCursor)),
	)

	if err != nil {
		return nil, "", err
	}

	response := struct {
		Success bool `json:"success"`
		Data    struct {
			HasNextPage bool   `json:"has_next_page"`
			EndCursor   string `json:"end_cursor"`
			Edges       []struct {
				Node instagram188Post `json:"node"`
			} `json:"edges"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, "", err
	}

	if !response.Success {
		return nil, "", errors.New("internal_server_error")
	}

	var posts []*Post

	for _, edge := range response.Data.Edges {
		postData := edge.Node
		if postData.Type != "GraphSidecar" && postData.Type != "GraphImage" {
			continue
		}

		pi := postData.ToPost()
		posts = append(posts, pi)
	}

	nextPageToken = ""

	if response.Data.HasNextPage {
		nextPageToken = response.Data.EndCursor
	}

	return posts, nextPageToken, nil
}

type instagram188Post struct {
	ID                 string `json:"id"`
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
			Node instagram188Post
		}
	} `json:"edge_sidecar_to_children"`
}

func (p instagram188Post) ToPost() *Post {
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

func (i *RapidAPIInstagram188) APIHost() string {
	return i.apiHost
}

func (i *RapidAPIInstagram188) HTTPClient() *http.Client {
	return i.client
}
