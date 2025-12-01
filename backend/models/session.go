package models

import "time"

type SessionPhone struct {
    SessionID     string    `gorm:"primaryKey;column:session_id" json:"session_id"`
    Phone         string    `gorm:"column:phone" json:"phone"`
    AIActive      bool      `gorm:"column:ai_active" json:"ai_active"`
    CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
    LastMessageAt time.Time `gorm:"column:last_message_at" json:"last_message_at"` 
    LeadName      string    `gorm:"column:lead_name" json:"lead_name"`
}

func (SessionPhone) TableName() string {
    return "session_phones"
}
