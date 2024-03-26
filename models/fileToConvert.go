package models

import "time"

type FileToConvert struct {
	Id            int       `json:"id" gorm:"primaryKey"`
	FileNameAsMp4 string    `json:"fileNameAsMp4"`
	FileNameAsMp3 string    `json:"fileNameAsMp3"`
	Status        int       `json:"status"`
	InsertedAt    time.Time `json:"insertedAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
