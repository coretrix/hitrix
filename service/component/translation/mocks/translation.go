package mocks

import (
	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/pkg/entity"
)

type FakeTranslationService struct {
	mock.Mock
}

func (f *FakeTranslationService) GetText(_ *beeorm.Engine, _ entity.TranslationTextLang, key entity.TranslationTextKey) (string, bool) {
	args := f.Called()
	if args.Get(0) == nil || args.Get(1) == nil {
		return string(key), false
	}

	return args.Get(0).(string), args.Bool(1)
}

func (f *FakeTranslationService) GetTextWithVars(_ *beeorm.Engine,
	_ *beeorm.Engine,
	_ entity.TranslationTextLang,
	key entity.TranslationTextKey,
	_ map[string]interface{},
) (string, bool) {
	args := f.Called()
	if args.Get(0) == nil || args.Get(1) == nil {
		return string(key), false
	}

	return args.Get(0).(string), args.Bool(1)
}
