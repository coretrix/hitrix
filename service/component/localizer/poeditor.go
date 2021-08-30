package localizer

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type PoeditorSource struct {
	apiKey    string
	projectId string
	language  string
}

type termListResponse struct {
	Response struct {
		Status  string
		Code    string
		Message string
	}
	Result struct {
		Terms []struct {
			Term        string
			Translation struct {
				Content string
			}
		}
	}
}

func (l *PoeditorSource) Pull() (pairs map[string]string, err error) {
	termList := termListResponse{}
	params := url.Values{}
	params.Add("api_token", l.apiKey)
	params.Add("id", l.projectId)
	params.Add("language", l.language)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://api.poeditor.com/v2/terms/list", body)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&termList)
	if err != nil {
		log.Fatal(err)
		return
	}

	if pairs == nil {
		pairs = make(map[string]string)
	}

	for _, term := range termList.Result.Terms {
		pairs[term.Term] = term.Translation.Content
	}

	return
}

func (l *PoeditorSource) Push(terms []string) (err error) {
	termsModel := []map[string]string{}
	for _, k := range terms {
		termsModel = append(termsModel, map[string]string{
			"term": k,
		})
	}
	b, err := json.Marshal(termsModel)
	if err != nil {
		log.Fatal(err)
		return
	}
	params := url.Values{}
	params.Add("api_token", l.apiKey)
	params.Add("id", l.projectId)
	params.Add("language", l.language)
	params.Add("data", string(b))
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://api.poeditor.com/v2/projects/sync", body)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	return
}

func NewPoeditorSource(apiKey string, projectId string, language string) *PoeditorSource {
	return &PoeditorSource{
		apiKey:    apiKey,
		projectId: projectId,
		language:  language,
	}
}
