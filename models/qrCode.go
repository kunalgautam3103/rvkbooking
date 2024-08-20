 package models

import (
	"time"
	"gorm.io/gorm"
)

type Qrdetails struct {
	Qrcode    			string    `gorm:"primary_key" json:"qrcode"`
	Count     			int    		`gorm:autoIncrement" json:"count"`
	First_scanned_at 	time.Time `json:"first_scanned_at"`
}

func MigrateQrcode(db *gorm.DB) error {
	err := db.AutoMigrate(&Qrdetails{})
	return err
}
