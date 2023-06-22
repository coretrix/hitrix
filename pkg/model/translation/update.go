package translation

import (
	"context"
	"fmt"

	"github.com/coretrix/hitrix/pkg/dto/translation"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/errors"
	"github.com/coretrix/hitrix/service"
)

func Update(ctx context.Context, request *translation.RequestUpdateTranslation, id uint64) (*translation.ResponseTranslation, error) {
	ormService := service.DI().OrmEngineForContext(ctx)

	translationTextEntity := &entity.TranslationTextEntity{}
	found := ormService.LoadByID(id, translationTextEntity)

	if !found {
		return nil,
			errors.HandleCustomErrors(map[string]string{"ID": fmt.Sprintf("City with id %v does not exists", id)})
	}

	translationTextEntity.Lang = request.Lang.String()
	translationTextEntity.Key = request.Key.String()
	translationTextEntity.Text = request.Text
	translationTextEntity.Status = entity.TranslationStatusTranslated.String()

	err := ormService.FlushWithCheck(translationTextEntity)
	if err != nil {
		return nil, errors.HandleFlushWithCheckError(
			err,
			errors.HandleCustomErrors(map[string]string{"Lang": "text with this lang and key already exists"}),
		)
	}

	return &translation.ResponseTranslation{
		ID:   translationTextEntity.ID,
		Lang: translationTextEntity.Lang,
		Key:  translationTextEntity.Key,
		Text: translationTextEntity.Text,
	}, nil
}
