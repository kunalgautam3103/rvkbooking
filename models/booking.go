package models

import "gorm.io/gorm"

type Booking struct {
	Name      *string    `json:"name"`
	Smartcard  *string  `json:"smartcard"`
	Category   *string  `json:"category"`
	Gateno		int    `json:"gateno"`
	Qrcode		*string  `json:"qr_code"`
}

func MigrateBooking(db *gorm.DB) error {
	err := db.AutoMigrate(&Booking{})
	return err
}