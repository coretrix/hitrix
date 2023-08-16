package mocks

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/service/component/geocoding"
)

type FakeGeocoding struct {
	mock.Mock
}

func (f *FakeGeocoding) Geocode(_ context.Context, _ *datalayer.DataLayer, address string, language geocoding.Language) (*geocoding.Address, error) {
	args := f.Called(address, language)

	return args.Get(0).(*geocoding.Address), args.Error(1)
}

func (f *FakeGeocoding) ReverseGeocode(
	_ context.Context,
	_ *datalayer.DataLayer,
	latLng *geocoding.LatLng,
	language geocoding.Language,
) (*geocoding.Address, error) {
	args := f.Called(latLng, language)

	return args.Get(0).(*geocoding.Address), args.Error(1)
}

func (f *FakeGeocoding) CutCoordinates(float float64, precision int) (float64, error) {
	asString := fmt.Sprintf("%.8f", float)
	asStringParts := strings.Split(asString, ".")

	return strconv.ParseFloat(asString[0:len(asStringParts[0])+1+precision], 64)
}
