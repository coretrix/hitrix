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
	beeorm.ORM `orm:"table=translation_texts;log=log_db_pool;redisCache;redisSearch=search_pool"`
	ID         uint64 `orm:"sortable"`
	Lang       string `orm:"required;unique=Lang_Key:1;searchable"`
	Key        string `orm:"required;unique=Lang_Key:2;searchable"`
	Status     string `orm:"required;enum=entity.TranslationStatusAll;searchable"`
	Text       string `orm:"length=max"`
	Vars       []string

	CachedQueryLangKey *beeorm.CachedQuery `queryOne:":Lang = ? AND :Key = ?"`
}
