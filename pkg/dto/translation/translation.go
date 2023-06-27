package translation

import (
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/crud"
)

type RequestCreateTranslation struct {
	Key  entity.TranslationTextKey  `binding:"required"`
	Lang entity.TranslationTextLang `binding:"required"`
	Text string                     `binding:"required"`
}

type RequestUpdateTranslation struct {
	Key  entity.TranslationTextKey  `binding:"required"`
	Lang entity.TranslationTextLang `binding:"required"`
	Text string                     `binding:"required"`
}

type RequestDTOTranslationID struct {
	ID uint64 `uri:"ID" binding:"required" example:"1"`
}

type ResponseTranslation struct {
	ID        uint64
	Lang      string
	Key       string
	Text      string
	Variables string
}

type ResponseDTOList struct {
	Rows    []*ListRow
	Total   int
	Columns []crud.Column
}

type ListRow struct {
	ID     uint64
	Status string
	Lang   string
	Key    string
	Text   string
}
