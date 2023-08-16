package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
)

type FakeTranslationService struct {
	mock.Mock
}

func (f *FakeTranslationService) GetText(_ *datalayer.DataLayer, _ entity.TranslationTextLang, key entity.TranslationTextKey) string {
	args := f.Called()
	if args.Get(0) == nil || args.Get(1) == nil {
		return string(key)
	}

	return args.String(0)
}

func (f *FakeTranslationService) GetTextWithVars(_ *datalayer.DataLayer,
	_ entity.TranslationTextLang,
	key entity.TranslationTextKey,
	_ map[string]interface{},
) string {
	args := f.Called()
	if args.Get(0) == nil || args.Get(1) == nil {
		return string(key)
	}

	return args.String(0)
}
