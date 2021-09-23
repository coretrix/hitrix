package trustpilot

import (
	"encoding/json"
)

const (
	eventProductReviewCreatedName = "product-review-created"

	eventProductReviewStatePublished = "published"
)

func ProductReviewCreatedWebhookDecoder(body string) ([]TrustpilotReview, error) {
	// NOTE: Field names here(webhook) may be different from the ones in the APIs. Leave as-is
	type response struct {
		Events []struct {
			EventName string `json:"eventName"`
			EventData struct {
				ID          string            `json:"id"`
				Stars       int               `json:"stars"`
				Content     string            `json:"content"`
				Product     TrustpilotProduct `json:"product"`
				ReferenceId string            `json:"referenceId"`
				Consumer    struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"consumer"`
				State     string `json:"state"`
				CreatedAt string `json:"createdAt"`
			} `json:"eventData"`
		} `json:"events"`
	}

	var responseData response
	if err := json.Unmarshal([]byte(body), &responseData); err != nil {
		return nil, err
	}

	reviews := make([]TrustpilotReview, 0)
	for _, event := range responseData.Events {
		if event.EventName != eventProductReviewCreatedName {
			continue
		}
		if event.EventData.State != eventProductReviewStatePublished {
			continue
		}
		review := TrustpilotReview{
			ID:      event.EventData.ID,
			Stars:   event.EventData.Stars,
			Content: event.EventData.Content,
			Product: &event.EventData.Product,
			Consumer: TrustpilotConsumer{
				ID:          event.EventData.Consumer.ID,
				DisplayName: event.EventData.Consumer.Name,
				Email:       event.EventData.Consumer.Email,
			},
			ReferenceId: event.EventData.ReferenceId,
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
