package models

import "time"

type ChatHistory struct {
    ID        uint      `gorm:"primaryKey;column:id" json:"id"`
    SessionID string    `gorm:"index;column:session_id" json:"session_id"`
    Message   string    `gorm:"column:message" json:"message"`
    CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (ChatHistory) TableName() string {
    return "n8n_chat_histories"
}
