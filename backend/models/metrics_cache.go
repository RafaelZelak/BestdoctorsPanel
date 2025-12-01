package models

import (
    "time"

    "gorm.io/datatypes"
)

type MetricsCache struct {
    MetricKey       string         `gorm:"primaryKey;column:metric_key" json:"metric_key"`
    Payload         datatypes.JSON `gorm:"column:payload" json:"payload"`
    LastRefreshedAt time.Time      `gorm:"column:last_refreshed_at" json:"last_refreshed_at"`
}

func (MetricsCache) TableName() string {
    return "bestdoctors_metrics_cache"
}
