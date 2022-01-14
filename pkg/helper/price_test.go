package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/pkg/helper"
)

func TestNewPrice(t *testing.T) {
	price := helper.NewPrice(229.90)

	assert.Equal(t, 229.90, price.Float())
	assert.Equal(t, int64(229900), price.Units())
	assert.Equal(t, "229.90", price.String())
}

func TestNewTotalPrice(t *testing.T) {
	price := helper.NewTotalPrice(10.29, 3)

	assert.Equal(t, 30.87, price.Float())
	assert.Equal(t, int64(30870), price.Units())
	assert.Equal(t, "30.87", price.String())
}
