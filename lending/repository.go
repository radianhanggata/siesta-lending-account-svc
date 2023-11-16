package lending

import (
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

func (r *Repository) InsertLending(data *model.Lending) error {
	err := r.db.Create(&data).Error
	if err != nil {
		return internalerror.ErrInternalServer
	}
	return nil
}

func (r *Repository) InsertRepayment(data []*model.Repayment) error {
	err := r.db.Create(&data).Error
	if err != nil {
		return internalerror.ErrInternalServer
	}
	return nil
}

func (r *Repository) GetUnpaidRepayment(accountID uint) ([]model.Repayment, error) {
	result := make([]model.Repayment, 0)
	err := r.db.Model(&model.Repayment{}).
		Where("account_id = ? and paid=false", accountID).
		Order("date asc").
		Find(&result).Error
	if err != nil {
		return result, internalerror.ErrInternalServer
	}
	return result, nil
}
