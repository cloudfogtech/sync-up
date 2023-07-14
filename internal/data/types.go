package data

import (
	"time"
)

type Log struct {
	ID        string    `gorm:"primarykey,index"`
	ServiceId string    `json:"serviceId" gorm:"not null"`
	Data      string    `json:"data" gorm:"not null"`
	StartAt   time.Time `json:"startAt" gorm:"index,not null"`
	EndAt     time.Time `json:"endAt" gorm:"not null"`
	Status    int64     `json:"status" gorm:"not null"`
	CreateAt  time.Time `json:"createAt" gorm:"not null"`
}
