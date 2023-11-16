package model

import "time"

type Repayment struct {
	ID           uint      `gorm:"primaryKey"`
	Date         time.Time `gorm:"not null"`
	Fee          float64   `gorm:"default=0"`
	FeeStampDuty float64   `gorm:"default=0"`
	Interest     float64   `gorm:"not null"`
	Principal    float64   `gorm:"not null"`
	AccountID    uint      `gorm:"not null"`
	LendingID    uint      `gorm:"not null"`
	Paid         bool      `gorm:"default=false"`
}

func (Repayment) TableName() string {
	return "repayment"
}
