package model

import "time"

type AchievementHistory struct {
	Status string    `json:"status"`
	At     time.Time `json:"at"`
}
