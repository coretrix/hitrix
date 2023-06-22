package translation

import (
	"context"
	"fmt"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
)

func Delete(ctx context.Context, id uint64) error {
	ormService := service.DI().OrmEngineForContext(ctx)

	translationTextEntity := &entity.TranslationTextEntity{}
	found := ormService.LoadByID(id, translationTextEntity)

	if !found {
		return fmt.Errorf("translation text with ID %v not found", id)
	}

	ormService.Delete(translationTextEntity)

	return nil
}
