package models

import "time"

type Merchant struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"not null;uniqueIndex" json:"name"`
	Email            string    `gorm:"not null;uniqueIndex" json:"email"`
	ContactFirstName string    `gorm:"not null" json:"contactFirstName"`
	ContactLastName  string    `gorm:"not null" json:"contactLastName"`
	KeycloakUserID   string    `gorm:"not null;uniqueIndex" json:"-"`
	CreatedByUserID  string    `gorm:"not null" json:"-"`
	CreatedAt        time.Time `json:"createdAt"`
}
