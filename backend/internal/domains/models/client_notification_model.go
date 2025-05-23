package models

type ClientNotification struct {
	ID             uint `json:"id" gorm:"primaryKey"`
	NotificationID uint `json:"notification_id"`
	UserID         uint `json:"user_id"`
	*Notification  `json:"notification" gorm:"foreignKey:NotificationID;references:ID"`
	*User          `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Seen           bool `json:"seen"`
}
