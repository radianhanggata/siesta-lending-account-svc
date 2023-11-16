package config

import (
	"errors"

	"gorm.io/gorm"

	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/internalerror"
	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/model"
)

type Repository struct {
	db *gorm.DB
}

func SetupRepository(db *gorm.DB) Repository {
	return Repository{db}
}

func (r *Repository) GetByID(id string) (model.Config, error) {
	config := model.Config{}
	err := r.db.Where("id = ?", id).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return config, internalerror.ErrNotFound
	}

	return config, err
}
