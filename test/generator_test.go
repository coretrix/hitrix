package main

import (
	"testing"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
)

func TestGenerateRandomRangeNumber(t *testing.T) {
	t.Run("generate otp code", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.GeneratorService(),
		)

		generatorService, _ := service.DI().GeneratorService()
		code := generatorService.GenerateRandomRangeNumber(1000, 9999)
		if code < 1000 || code > 9999 {
			t.Fail()
		}
	})
}
