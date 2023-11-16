package model

type Config struct {
	ID                   string  `gorm:"primaryKey"`
	Fee                  float64 `gorm:"not null"`
	Interest             float64 `gorm:"not null"`
	OutstandingThreshold float64 `gorm:"not null"`
	OutstandingFee       float64 `gorm:"not null"`
}

func (Config) TableName() string {
	return "config"
}
