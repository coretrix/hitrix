package translation

import (
	"fmt"
	"github.com/coretrix/hitrix/service/component/app"
	"strings"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

type ITranslationService interface {
	GetText(ormService *beeorm.Engine, lang entity.TranslationTextLang, key entity.TranslationTextKey) string
	GetTextWithVars(
		ormService *beeorm.Engine,
		lang entity.TranslationTextLang,
		key entity.TranslationTextKey,
		variables map[string]interface{},
	) string
}

type translationService struct {
	appService *app.App
}

func NewTranslationService(appService *app.App) ITranslationService {
	return &translationService{appService}
}

func (u *translationService) GetText(ormService *beeorm.Engine, lang entity.TranslationTextLang, key entity.TranslationTextKey) string {
	var found bool
	translationTextEntity := &entity.TranslationTextEntity{}

	found = ormService.CachedSearchOne(
		translationTextEntity,
		"CachedQueryKey",
		key)

	if !found {
		if !u.appService.IsInTestMode() {
			translationTextEntity.Status = entity.TranslationStatusNew.String()
			translationTextEntity.Lang = lang.String()
			translationTextEntity.Key = key.String()
		}

		ormService.Flush(translationTextEntity)

		return key.String()
	}

	return translationTextEntity.Text
}

func (u *translationService) GetTextWithVars(
	ormService *beeorm.Engine,
	lang entity.TranslationTextLang,
	key entity.TranslationTextKey,
	variables map[string]interface{}) string {
	text := u.GetText(ormService, lang, key)

	for paramName, value := range variables {
		text = strings.Replace(text, fmt.Sprintf("[[%s]]", paramName), fmt.Sprintf("%v", value), -1)
	}

	return text
}
