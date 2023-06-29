package licenseplaterecognizer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type PlateRecognizer struct {
	APIKey string
}

func NewPlateRecognizer(apiKey string) LicensePlateRecognizer {
	return &PlateRecognizer{APIKey: apiKey}
}

func (pr *PlateRecognizer) RecognizeFromImage(base64image string) ([]string, error) {
	form := url.Values{}
	form.Add("upload", base64image)
	form.Add("regions", "bg")

	req, err := http.NewRequest("POST", "https://api.platerecognizer.com/v1/plate-reader/", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", pr.APIKey))

	httpClient := http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status code %d but service returned  %d", http.StatusCreated, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &plates{}

	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	returnData := make([]string, len(result.Plates))

	for i, plate := range result.Plates {
		if plate.Plate == "" {
			return nil, fmt.Errorf("service returned empty license plate")
		}

		returnData[i] = strings.ToUpper(plate.Plate)
	}

	return returnData, nil
}

type plates struct {
	Plates []*plate `json:"results"`
}

type plate struct {
	Plate string `json:"plate"`
}
