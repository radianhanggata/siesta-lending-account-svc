package model

import "time"

type Lending struct {
	ID        uint      `gorm:"primaryKey"`
	Date      time.Time `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	Tenor     int       `gorm:"not null"`
	Fee       float64   `gorm:"not null"`
	Interest  float64   `gorm:"not null"`
	AccountID uint      `gorm:"not null"`
}

func (Lending) TableName() string {
	return "lending"
}
