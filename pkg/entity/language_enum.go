package entity

const (
	LanguageEnglish   = "en"
	LanguageBulgarian = "bg"
)

type languageValue struct {
	LanguageEnglish   string
	LanguageBulgarian string
}

var LanguageValueAll = languageValue{
	LanguageEnglish:   LanguageEnglish,
	LanguageBulgarian: LanguageBulgarian,
}
