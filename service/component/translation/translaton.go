package translation

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/helper"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
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
	errorLoggerService errorlogger.ErrorLogger
}

func NewTranslationService(errorLoggerService errorlogger.ErrorLogger) ITranslationService {
	return &translationService{errorLoggerService}
}

func (u *translationService) GetText(ormService *beeorm.Engine, lang entity.TranslationTextLang, key entity.TranslationTextKey) string {
	translationTextEntity := &entity.TranslationTextEntity{}

	found := ormService.CachedSearchOne(
		translationTextEntity,
		"CachedQueryLangKey",
		lang.String(),
		key.String())

	if !found {
		translationTextEntity.Status = entity.TranslationStatusNew.String()
		translationTextEntity.Lang = lang.String()
		translationTextEntity.Key = key.String()

		ormService.Flush(translationTextEntity)

		return key.String()
	}

	return translationTextEntity.Text
}

func (u *translationService) GetTextWithVars(
	ormService *beeorm.Engine,
	lang entity.TranslationTextLang,
	key entity.TranslationTextKey,
	variables map[string]interface{},
) string {
	keys := make([]string, 0, len(variables))

	for k := range variables {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	translationTextEntity := &entity.TranslationTextEntity{}

	found := ormService.CachedSearchOne(
		translationTextEntity,
		"CachedQueryLangKey",
		lang.String(),
		key.String())

	if !found {
		translationTextEntity.Status = entity.TranslationStatusNew.String()
		translationTextEntity.Lang = lang.String()
		translationTextEntity.Key = key.String()
		translationTextEntity.Vars = keys

		ormService.Flush(translationTextEntity)

		return key.String()
	}

	if !helper.EqualString(translationTextEntity.Vars, keys) {
		translationTextEntity.Vars = keys
		ormService.FlushLazy(translationTextEntity)
	}

	if translationTextEntity.Status == entity.TranslationStatusNew.String() {
		return key.String()
	}

	text := translationTextEntity.Text

	for paramName, value := range variables {
		text = strings.Replace(text, fmt.Sprintf("[[%s]]", paramName), fmt.Sprintf("%v", value), -1)
	}

	re := regexp.MustCompile(`\[\[(.*?)\]\]`)
	subMatchAll := re.FindAllString(text, -1)

	if len(subMatchAll) > 0 {
		u.errorLoggerService.LogError(
			fmt.Sprintf(
				"not assigned vars (%s) for translation key `%s`",
				strings.Join(subMatchAll, ", "),
				key.String(),
			),
		)
	}

	return text
}
