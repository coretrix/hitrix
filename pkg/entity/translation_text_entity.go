package entity

import (
	"github.com/latolukasz/beeorm"
)

type TranslationTextLang string

func (u TranslationTextLang) String() string {
	return string(u)
}

type TranslationTextKey string

func (u TranslationTextKey) String() string {
	return string(u)
}

type TranslationStatus string

func (u TranslationStatus) String() string {
	return string(u)
}

const (
	TranslationStatusNew        TranslationStatus = "new"
	TranslationStatusTranslated TranslationStatus = "translated"
)

type translationStatus struct {
	New        string
	Translated string
}

var TranslationStatusAll = translationStatus{
	New:        TranslationStatusNew.String(),
	Translated: TranslationStatusTranslated.String(),
}

type TranslationTextEntity struct {
	beeorm.ORM `orm:"table=translation_texts;log=log_db_pool;redisCache;localCache"`
	ID         uint64
	Lang       string `orm:"required;unique=Lang_Key:1"`
	Key        string `orm:"required;unique=Lang_Key:2"`
	Status     string `orm:"required;enum=entity.TranslationStatusAll"`
	Text       string `orm:"length=max"`
	Vars       []string

	CachedQueryLangKey *beeorm.CachedQuery `queryOne:":Lang = ? AND :Key = ?"`
}
