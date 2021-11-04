package featureflag

import (
	"errors"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/latolukasz/beeorm"
)

type serviceFeatureFlag struct {
}

func NewFeatureFlagService() ServiceFeatureFlagInterface {
	return &serviceFeatureFlag{}
}

func (s *serviceFeatureFlag) IsActive(ormService *beeorm.Engine, name string) bool {
	if name == "" {
		panic("name cannot be empty")
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)
	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return false
	}

	return featureFlagEntity.IsActive
}

func (s *serviceFeatureFlag) FailIfIsNotActive(ormService *beeorm.Engine, name string) error {
	isActive := s.IsActive(ormService, name)
	if !isActive {
		return errors.New("feature is not active")
	}

	return nil
}

func (s *serviceFeatureFlag) Activate(ormService *beeorm.Engine, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)
	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return errors.New("feature cannot be found")
	}

	featureFlagEntity.IsActive = true
	ormService.Flush(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) DeActivate(ormService *beeorm.Engine, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)
	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return errors.New("feature cannot be found")
	}

	featureFlagEntity.IsActive = false
	ormService.Flush(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) Create(ormService *beeorm.Engine, clockService clock.IClock, name string, isActive bool) error {
	if name == "" {
		panic("name cannot be empty")
	}

	featureFlagEntity := &entity.FeatureFlagEntity{
		Name:      name,
		IsActive:  isActive,
		UpdatedAt: nil,
		CreatedAt: clockService.Now(),
	}

	ormService.Flush(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) Delete(ormService *beeorm.Engine, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)
	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return errors.New("feature cannot be found")
	}

	ormService.Delete(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) GetAll(ormService *beeorm.Engine, pager *beeorm.Pager) []*entity.FeatureFlagEntity {
	query := beeorm.NewRedisSearchQuery()
	var featureFlagEntities []*entity.FeatureFlagEntity
	ormService.RedisSearch(&featureFlagEntities, query, pager)

	return featureFlagEntities
}
