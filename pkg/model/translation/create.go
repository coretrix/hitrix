package translation

import (
	"context"

	"github.com/coretrix/hitrix/pkg/dto/translation"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/errors"
	"github.com/coretrix/hitrix/service"
)

func Create(ctx context.Context, request *translation.RequestCreateTranslation) (*translation.ResponseTranslation, error) {
	ormService := service.DI().OrmEngineForContext(ctx)

	newTranslationEntity := &entity.TranslationTextEntity{
		Lang:   request.Lang.String(),
		Key:    request.Key.String(),
		Status: entity.TranslationStatusTranslated.String(),
		Text:   request.Text,
	}

	err := ormService.FlushWithCheck(newTranslationEntity)
	if err != nil {
		return nil, errors.HandleFlushWithCheckError(
			err,
			errors.HandleCustomErrors(map[string]string{"Lang": "text with this lang and key already exists"}),
		)
	}

	return &translation.ResponseTranslation{
		ID:   newTranslationEntity.ID,
		Lang: newTranslationEntity.Lang,
		Key:  newTranslationEntity.Key,
		Text: newTranslationEntity.Text,
	}, nil
}
