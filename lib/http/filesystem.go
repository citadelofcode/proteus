package http

import (
	"time"
)

type File struct {
	Contents []byte
	ContentType string
	LastModifiedAt time.Time
}