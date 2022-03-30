package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/coretrix/hitrix/service/component/config"
)

// RapidAPIInstagramSimple https://rapidapi.com/premium-apis-premium-apis-default/api/instagram85/
type RapidAPIInstagramSimple struct {
	client       *http.Client
	apiKey       string
	apiHost      string
	apiBaseURL   string
	providerName string
}

func NewRapidAPIInstagramSimple(configService config.IConfig) (IProvider, error) {
	return &RapidAPIInstagramSimple{
		client:       &http.Client{},
		apiKey:       configService.MustString("instagram.api.rapid_api_token"),
		apiHost:      "instagram85.p.rapidapi.com",
		apiBaseURL:   "https://instagram85.p.rapidapi.com",
		providerName: "RapidAPIInstagram85",
	}, nil
}

func (i *RapidAPIInstagramSimple) GetName() string {
	return i.providerName
}

func (i *RapidAPIInstagramSimple) APIKey() string {
	return i.apiKey
}

func (i *RapidAPIInstagramSimple) GetAccount(account string) (*Account, error) {
	response := struct {
		Data struct {
			ID       int64  `json:"id"`
			Fullname string `json:"full_name"`
			Bio      string `json:"biography"`
			Website  string `json:"website"`
			Picture  struct {
				HD string `json:"hd"`
			} `json:"profile_picture"`
			Figures struct {
				Posts      int64 `json:"posts"`
				Followers  int64 `json:"followers"`
				Followings int64 `json:"followings"`
			} `json:"figures"`
			IsPrivate bool `json:"is_private"`
		} `json:"data"`
		Code int64 `json:"code"`
	}{}
	res, err := sendRapidRequest(i, fmt.Sprintf("%v/account/%v/info", i.apiBaseURL, account))
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, &response); err != nil {
		errResponse := struct {
			Code int64 `json:"code"`
		}{}
		if err2 := json.Unmarshal(res, &errResponse); err2 != nil {
			return nil, err2
		}

		if errResponse.Code == 302 {
			return nil, errors.New("internal_server_error status code 302")
		}
		return nil, err
	}

	if response.Code == 404 {
		return nil, errors.New("internal_server_error status code 404")
	} else if response.Code != 200 {
		return nil, errors.New("internal_server_error status code not 200")
	}

	return &Account{
		AccountID: response.Data.ID,
		FullName:  response.Data.Fullname,
		Bio:       response.Data.Bio,
		Posts:     response.Data.Figures.Posts,
		Followers: response.Data.Figures.Followers,
		Following: response.Data.Figures.Followings,
		Picture:   response.Data.Picture.HD,
		IsPrivate: response.Data.IsPrivate,
		Website:   response.Data.Website,
	}, nil
}

func (i *RapidAPIInstagramSimple) GetFeed(accountID int64, nextPageToken string) ([]*Post, string, error) {
	response := struct {
		Data []RapidAPIInstagramSimplePost `json:"data"`
		Meta struct {
			HasNext  bool   `json:"has_next"`
			NextPage string `json:"next_page"`
		} `json:"meta"`
		Code int64 `json:"code"`
	}{}

	instagramURL := fmt.Sprintf("%s/account/%d/feed", i.apiBaseURL, accountID)
	if nextPageToken != "" {
		instagramURL += fmt.Sprintf("?pageId=%s", nextPageToken)
	}

	res, err := sendRapidRequest(i, instagramURL)

	if err != nil {
		return nil, "", err
	}

	if err = json.Unmarshal(res, &response); err != nil {
		return nil, "", err
	}

	if response.Code == 404 {
		return nil, "", errors.New("internal_server_error status code 404")
	} else if response.Code != 200 {
		return nil, "", errors.New("internal_server_error status code not 200")
	}

	var posts []*Post

	for i := range response.Data {
		postData := response.Data[i]
		if postData.Type != "image" && postData.Type != "sidecar" {
			continue
		}

		if post := postData.ToPost(); post != nil {
			posts = append(posts, post)
		}
	}

	nextPageToken = ""

	if response.Meta.HasNext {
		nextPageToken = response.Meta.NextPage
	}

	return posts, nextPageToken, nil
}

func (i *RapidAPIInstagramSimple) APIHost() string {
	return i.apiHost
}

func (i *RapidAPIInstagramSimple) HTTPClient() *http.Client {
	return i.client
}

func (p RapidAPIInstagramSimplePost) ToPost() *Post {
	post := &Post{
		ID:        p.ID,
		Title:     p.Caption,
		CreatedAt: p.CreatedTime.Unix,
	}

	var images []string
	if p.Type == "image" {
		images = []string{
			p.Images.Original.High,
		}
	} else {
		for j := range p.Sidecar {
			if p.Sidecar[j].Type == "image" {
				images = append(images, p.Sidecar[j].Images.Original.High)
			}
		}
	}

	post.Images = images

	return post
}

type RapidAPIInstagramSimplePost struct {
	ID          string `json:"id"`
	CreatedTime struct {
		Unix int64 `json:"unix"`
	} `json:"created_time"`
	Caption string `json:"caption"`
	Type    string `json:"type"`
	Images  struct {
		Original struct {
			High string `json:"high"`
		} `json:"original"`
	} `json:"images"`
	Videos struct {
		Standard string `json:"standard"`
	} `json:"videos"`
	Sidecar []RapidAPIInstagramSimplePost `json:"sidecar"`
}

type rapidAPIProvider interface {
	HTTPClient() *http.Client
	APIKey() string
	APIHost() string
}

func sendRapidRequest(provider rapidAPIProvider, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-key", provider.APIKey())
	req.Header.Add("x-rapidapi-host", provider.APIHost())

	res, err := provider.HTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
