package mocks

import (
	"context"

	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/geocoding"
)

type FakeGeocoding struct {
	mock.Mock
}

func (f *FakeGeocoding) Geocode(_ context.Context, _ *beeorm.Engine, address string, language geocoding.Language) (*geocoding.Address, error) {
	args := f.Called(address, language)

	return args.Get(0).(*geocoding.Address), args.Error(1)
}

func (f *FakeGeocoding) ReverseGeocode(
	_ context.Context,
	_ *beeorm.Engine,
	latLng *geocoding.LatLng,
	language geocoding.Language,
) (*geocoding.Address, error) {
	args := f.Called(latLng, language)

	return args.Get(0).(*geocoding.Address), args.Error(1)
}
