package trustpilot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// https://documentation-apidocumentation.trustpilot.com/private-products-api#batch-upsert-products
func (tp *TrustpilotAPI) upsertProducts(products []TrustpilotProduct) ([]TrustpilotProduct, error) {
	endpoint := fmt.Sprintf(upsertProductsEndpoint, tp.BusinessUnitID)

	requestData := map[string]interface{}{
		"products":                               products,
		"skuSameAsGoogleMerchantCenterProductId": true,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := makeTrustpilotAuthenticatedRequest(tp.AccessToken.AccessToken, http.MethodPost, endpoint, nil, body)
	if err != nil {
		return nil, err
	}

	client, err := makeTrustpilotAuthenticatedClient(tp.AccessToken.AccessToken)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trustpilot upsertProducts request failed with status code ( %d )", resp.StatusCode)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	type response struct {
		Products []TrustpilotProduct `json:"products"`
	}
	var responseData response
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	return responseData.Products, nil
}

// https://developers.trustpilot.com/product-reviews-api#get-product-reviews-summary
func (tp *TrustpilotAPI) getProductReviewsSummary(product TrustpilotProduct) (*TrustpilotProductReviewSummary, error) {
	endpoint := fmt.Sprintf(getProductReviewsSummaryEndpoint, tp.BusinessUnitID)

	params := url.Values{}
	params.Add("sku", product.Sku)

	req, err := makeTrustpilotAuthenticatedRequest(tp.AccessToken.AccessToken, http.MethodGet, endpoint, params, nil)
	if err != nil {
		return nil, err
	}

	client, err := makeTrustpilotAuthenticatedClient(tp.AccessToken.AccessToken)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trustpilot getProductReviewsSummary request failed with status code ( %d )", resp.StatusCode)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	type response struct {
		NumberOfReviews struct {
			Total      int `json:"total"`
			OneStar    int `json:"oneStar"`
			TwoStars   int `json:"twoStars"`
			ThreeStars int `json:"threeStars"`
			FourStars  int `json:"fourStars"`
			FiveStars  int `json:"fiveStars"`
		} `json:"numberOfReviews"`
		StarsAverage float64 `json:"starsAverage"`
	}

	var responseData response
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	var reviewsSummary TrustpilotProductReviewSummary
	reviewsSummary.Total = strconv.Itoa(responseData.NumberOfReviews.Total)
	reviewsSummary.Average = fmt.Sprintf("%.1f", responseData.StarsAverage)

	return &reviewsSummary, nil
}

// https://developers.trustpilot.com/product-reviews-api#get-private-product-reviews
func (tp *TrustpilotAPI) getProductReviews(productSKU string, page *int, perPage *int) ([]TrustpilotReview, error) {
	endpoint := fmt.Sprintf(getProductReviewsEndpoint, tp.BusinessUnitID)

	reqPage := 1
	if page != nil {
		reqPage = *page
	}

	reqPerPage := 100
	if perPage != nil {
		reqPage = *perPage
	}

	params := url.Values{}
	params.Add("page", strconv.Itoa(reqPage))
	params.Add("perPage", strconv.Itoa(reqPerPage))
	params.Add("sku", productSKU)

	req, err := makeTrustpilotAuthenticatedRequest(tp.AccessToken.AccessToken, http.MethodPost, endpoint, params, nil)
	if err != nil {
		return nil, err
	}

	client, err := makeTrustpilotAuthenticatedClient(tp.AccessToken.AccessToken)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trustpilot getProductReviews request failed with status code ( %d )", resp.StatusCode)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	type response struct {
		Reviews []struct {
			ID          string            `json:"id"`
			Content     string            `json:"content"`
			Stars       int               `json:"stars"`
			ReferenceId string            `json:"referenceId"`
			Product     TrustpilotProduct `json:"product"`
			Consumer    struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"consumer"`
		} `json:"productReviews"`
		Total int `json:"total"`
	}

	var responseData response
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	reviews := make([]TrustpilotReview, 0)

	for _, v := range responseData.Reviews {
		review := TrustpilotReview{
			ID:      v.ID,
			Stars:   v.Stars,
			Content: v.Content,
			Product: &v.Product,
			Consumer: TrustpilotConsumer{
				ID:          v.Consumer.ID,
				DisplayName: v.Consumer.Name,
				Email:       v.Consumer.Email,
			},
			ReferenceId: v.ReferenceId,
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// https://developers.trustpilot.com/product-reviews-api#create-product-review-invitation-link
func (tp *TrustpilotAPI) createProductReviewInvitationLink(productID string, userEmail string, userName string, refID string) (*TrustpilotProductReviewInvitation, error) {
	endpoint := fmt.Sprintf(createProductReviewInvitationLink, tp.BusinessUnitID)

	requestData := map[string]interface{}{
		"consumer": map[string]interface{}{
			"email": userEmail,
			"name":  userName,
		},
		"referenceId": refID,
		"locale":      "en-US",
		"redirectUri": "https://www.loveyourself.co.uk/", // TODO trustpilot redirect catcher to close mobile webview
		"productIds": []string{
			productID,
		},
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := makeTrustpilotAuthenticatedRequest(tp.AccessToken.AccessToken, http.MethodPost, endpoint, nil, body)
	if err != nil {
		return nil, err
	}

	client, err := makeTrustpilotAuthenticatedClient(tp.AccessToken.AccessToken)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trustpilot createProductReviewInvitationLink request failed with status code ( %d )", resp.StatusCode)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	type response struct {
		ID  string `json:"reviewLinkId"`
		URL string `json:"reviewUrl"`
	}

	var responseData response
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	var reviewInvitation TrustpilotProductReviewInvitation
	reviewInvitation.ID = responseData.ID
	reviewInvitation.URL = responseData.URL

	return &reviewInvitation, nil
}
