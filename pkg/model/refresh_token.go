package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Token      string    `json:"token" gorm:"type:text;not null;unique;index"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	ExpiresAt  int64     `json:"expires_at" gorm:"not null;index"`
	IsRevoked  bool      `json:"is_revoked" gorm:"default:false;index"`
	DeviceInfo string    `json:"device_info" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"not null"`
	User       User      `json:"-" gorm:"foreignKey:UserID"`
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return nil
}
