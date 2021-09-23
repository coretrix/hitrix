package mocks

import (
	"github.com/coretrix/hitrix/service/component/trustpilot"
	"github.com/stretchr/testify/mock"
)

type FakeTrustpilot struct {
	mock.Mock
}

func (m *FakeTrustpilot) UpsertProducts(products []trustpilot.TrustpilotProduct) ([]trustpilot.TrustpilotProduct, error) {
	return m.Called(products).Get(0).([]trustpilot.TrustpilotProduct), nil
}

func (m *FakeTrustpilot) GetProductReviewsSummary(product trustpilot.TrustpilotProduct) (*trustpilot.TrustpilotProductReviewSummary, error) {
	return m.Called(product).Get(0).(*trustpilot.TrustpilotProductReviewSummary), nil
}

func (m *FakeTrustpilot) GetProductReviews(productSKU string, page *int, perPage *int) ([]trustpilot.TrustpilotReview, error) {
	return m.Called(productSKU, page, perPage).Get(0).([]trustpilot.TrustpilotReview), nil
}

func (m *FakeTrustpilot) CreateProductReviewInvitationLink(productID string, userEmail string, userName string, refID string) (*trustpilot.TrustpilotProductReviewInvitation, error) {
	return m.Called(productID, userEmail, userName, refID).Get(0).(*trustpilot.TrustpilotProductReviewInvitation), nil
}
