package geocoding

import "github.com/coretrix/hitrix/pkg/entity"

const (
	LanguageEnglish   = Language("en")
	LanguageBulgarian = Language("bg")
)

type Language string

var languageToEnumMapping = map[Language]string{
	LanguageEnglish:   entity.LanguageEnglish,
	LanguageBulgarian: entity.LanguageBulgarian,
}
