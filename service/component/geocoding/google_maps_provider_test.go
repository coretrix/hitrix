package geocoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"googlemaps.github.io/maps"
)

func TestExcludeFormattedAddressComponentsRemovesNeighborhood(t *testing.T) {
	formattedAddress := `Tsentar, ул. „Стефан Стамболов“ 44, 2140 Botevgrad, Bulgaria`
	addressComponents := []maps.AddressComponent{
		{LongName: "Tsentar", ShortName: "Tsentar", Types: []string{"neighborhood", "political"}},
	}

	actual := excludeFormattedAddressComponents(
		formattedAddress,
		addressComponents,
		excludedFormattedAddressComponentTypes,
	)

	assert.Equal(t, `ул. „Стефан Стамболов“ 44, 2140 Botevgrad, Bulgaria`, actual)
}

func TestExcludeFormattedAddressComponentsRemovesAdministrativeAreaLevel1(t *testing.T) {
	formattedAddress := `ул. „Стефан Стамболов“ 44, 2140 Botevgrad, Sofiyska oblast, Bulgaria`
	addressComponents := []maps.AddressComponent{
		{LongName: "Sofiyska oblast", ShortName: "Sofia Province", Types: []string{"administrative_area_level_1", "political"}},
	}

	actual := excludeFormattedAddressComponents(
		formattedAddress,
		addressComponents,
		excludedFormattedAddressComponentTypes,
	)

	assert.Equal(t, `ул. „Стефан Стамболов“ 44, 2140 Botevgrad, Bulgaria`, actual)
}
