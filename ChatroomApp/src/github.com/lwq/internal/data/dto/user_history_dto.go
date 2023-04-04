package dto

import (
	"time"
)

type UserHistoryDto struct {
	Role     string
	Message  string
	Datetime time.Time
}
