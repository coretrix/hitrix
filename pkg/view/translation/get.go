package translation

import (
	"context"
	"fmt"
	"strings"

	"github.com/coretrix/hitrix/pkg/dto/translation"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/errors"
	"github.com/coretrix/hitrix/service"
)

func Get(ctx context.Context, id uint64) (*translation.ResponseTranslation, error) {
	ormService := service.DI().OrmEngineForContext(ctx)

	translationEntity := &entity.TranslationTextEntity{}

	found := ormService.LoadByID(id, translationEntity)

	if !found {
		return nil, errors.HandleCustomErrors(map[string]string{
			"ID": fmt.Sprintf("%d does not exists", id),
		})
	}

	return &translation.ResponseTranslation{
		ID:        translationEntity.ID,
		Lang:      translationEntity.Lang,
		Key:       translationEntity.Key,
		Text:      translationEntity.Text,
		Variables: strings.Join(translationEntity.Vars, " "),
	}, nil
}
