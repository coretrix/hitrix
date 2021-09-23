package trustpilot

import (
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/latolukasz/beeorm"
)

const (
	authEndpoint                      = "https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/accesstoken"
	authRefreshEndpoint               = "https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/refresh"
	upsertProductsEndpoint            = "https://api.trustpilot.com/v1/private/business-units/%s/products"
	getProductReviewsSummaryEndpoint  = "https://api.trustpilot.com/v1/product-reviews/business-units/%s"
	getProductReviewsEndpoint         = "https://api.trustpilot.com/v1/private/product-reviews/business-units/%s/reviews"
	createProductReviewInvitationLink = "https://api.trustpilot.com/v1/private/product-reviews/business-units/%s/invitation-links"
)

type ITrustpilot interface {
	UpsertProducts(products []TrustpilotProduct) ([]TrustpilotProduct, error)
	GetProductReviewsSummary(product TrustpilotProduct) (*TrustpilotProductReviewSummary, error)
	GetProductReviews(productSKU string, page *int, perPage *int) ([]TrustpilotReview, error)
	CreateProductReviewInvitationLink(productID string, userEmail string, userName string, refID string) (*TrustpilotProductReviewInvitation, error)
}

type TrustpilotAPI struct {
	apiKey         string
	apiSecret      string
	username       string
	password       string
	BusinessUnitID string
	AccessToken    *AccessToken
	clockService   clock.IClock
	ormService     *beeorm.Engine
}

type TrustpilotProductReviewSummary struct {
	Total   string `json:"total"`
	Average string `json:"average"`
}

type TrustpilotProductReviewInvitation struct {
	ID  string `json:"reviewLinkId"`
	URL string `json:"reviewUrl"`
}

type TrustpilotProduct struct {
	ID          string `json:"id"`
	Sku         string `json:"sku"`
	Title       string `json:"title"` // in api
	Name        string `json:"name"`  // in webhook
	Link        string `json:"link"`
	Price       string `json:"price"`
	Description string `json:"description"`
}

type TrustpilotConsumer struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type TrustpilotReview struct {
	ID          string             `json:"id"`
	Stars       int                `json:"stars"`
	Content     string             `json:"content"`
	Product     *TrustpilotProduct `json:"product"`
	Consumer    TrustpilotConsumer `json:"consumer"`
	ReferenceId string             `json:"referenceId"`
}

func (tp *TrustpilotAPI) UpsertProducts(products []TrustpilotProduct) ([]TrustpilotProduct, error) {
	if err := tp.authenticate(); err != nil {
		return nil, err
	}

	return tp.upsertProducts(products)
}

func (tp *TrustpilotAPI) GetProductReviewsSummary(product TrustpilotProduct) (*TrustpilotProductReviewSummary, error) {
	if err := tp.authenticate(); err != nil {
		return nil, err
	}

	return tp.getProductReviewsSummary(product)
}

func (tp *TrustpilotAPI) GetProductReviews(productSKU string, page *int, perPage *int) ([]TrustpilotReview, error) {
	if err := tp.authenticate(); err != nil {
		return nil, err
	}

	return tp.getProductReviews(productSKU, page, perPage)
}

func (tp *TrustpilotAPI) CreateProductReviewInvitationLink(productID string, userEmail string, userName string, refID string) (*TrustpilotProductReviewInvitation, error) {
	if err := tp.authenticate(); err != nil {
		return nil, err
	}

	return tp.createProductReviewInvitationLink(productID, userEmail, userName, refID)
}

func NewTrustpilot(apikey string, apiSecret string, username string, password string, businessUnitID string, ormService *beeorm.Engine, clockService clock.IClock) (ITrustpilot, error) {
	now := clockService.Now()

	accessToken, err := getSettingsAccessToken(ormService)
	if err != nil {
		return nil, err
	}

	if accessToken == nil || accessToken.NeedsRenew(now) {
		var err error
		if accessToken, err = getNewAccessToken(apikey, apiSecret, username, password); err != nil {
			return nil, err
		}
	} else if accessToken.NeedsRefresh(clockService.Now()) {
		if accessToken, err = refreshAccessToken(apikey, apiSecret, accessToken.RefreshToken); err != nil {
			return nil, err
		}
		err = setSettingsAccessToken(ormService, *accessToken)
		if err != nil {
			return nil, err
		}
	}

	return &TrustpilotAPI{
		apiKey:         apikey,
		apiSecret:      apiSecret,
		username:       username,
		password:       password,
		BusinessUnitID: businessUnitID,
		AccessToken:    accessToken,
		ormService:     ormService,
		clockService:   clockService,
	}, nil
}
