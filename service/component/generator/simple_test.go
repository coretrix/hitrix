package generator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/service/component/generator"
)

func TestGenerateRandomRangeNumber(t *testing.T) {
	generatorService := &generator.SimpleGenerator{}

	code := generatorService.GenerateRandomRangeNumber(1000, 9999)

	assert.True(t, code >= 1000 && code <= 9999)
}
