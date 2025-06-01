package logar

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/expr-lang/expr"
	"sadk.dev/logar/models"
)

type FeatureFlags interface {
	Common

	HasFeatureFlag(ctx context.Context, flag string) (bool, error)

	GetFeatureFlags() ([]models.FeatureFlag, error)
	GetFeatureFlagByName(name string) (models.FeatureFlag, error)
	GetFeatureFlag(id uint) (models.FeatureFlag, error)
	CreateFeatureFlag(flag *models.FeatureFlag) error
	UpdateFeatureFlag(flag *models.FeatureFlag) error
	DeleteFeatureFlag(id uint) error
}

type FeatureFlagsImpl struct {
	core *AppImpl
}

func (f *FeatureFlagsImpl) GetApp() App {
	return f.core
}

func (f *FeatureFlagsImpl) GetFeatureFlags() ([]models.FeatureFlag, error) {
	var flags []models.FeatureFlag
	err := f.core.db.Find(&flags).Error
	if err != nil {
		return nil, err
	}

	return flags, nil
}

func (f *FeatureFlagsImpl) GetFeatureFlagByName(name string) (models.FeatureFlag, error) {
	var flag models.FeatureFlag
	err := f.core.db.Where("name = ?", name).First(&flag).Error
	if err != nil {
		return models.FeatureFlag{}, err
	}

	return flag, nil
}

func (f *FeatureFlagsImpl) GetFeatureFlag(id uint) (models.FeatureFlag, error) {
	var flag models.FeatureFlag
	err := f.core.db.Where("id = ?", id).First(&flag).Error
	if err != nil {
		return models.FeatureFlag{}, err
	}

	return flag, nil
}

func (f *FeatureFlagsImpl) CreateFeatureFlag(flag *models.FeatureFlag) error {
	err := f.core.db.Create(flag).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *FeatureFlagsImpl) UpdateFeatureFlag(flag *models.FeatureFlag) error {
	err := f.core.db.Save(flag).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *FeatureFlagsImpl) DeleteFeatureFlag(id uint) error {
	err := f.core.db.Where("id = ?", id).Delete(&models.FeatureFlag{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *FeatureFlagsImpl) HasFeatureFlag(ctx context.Context, flag string) (bool, error) {
	values, ok := f.core.GetContextValues(ctx)
	if !ok {
		values = Map{}
	}

	featureflag, err := f.GetFeatureFlagByName(flag)
	if err != nil {
		return false, err
	}

	if !featureflag.Enabled {
		return false, nil
	}

	timeout, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	contextValues := map[string]any{}
	for k, v := range values {
		contextValues[k] = v
	}

	env := map[string]any{
		"context": contextValues,
		"println": fmt.Println,
		"ctx":     timeout,
	}

	program, err := expr.Compile(featureflag.Condition, expr.Env(env), expr.AsBool(), expr.WithContext("ctx"))
	if err != nil {
		return false, err
	}

	result, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}

	enabled, ok := result.(bool)
	if !ok {
		return false, errors.New("feature flag result is not a boolean")
	}

	return enabled, nil
}
