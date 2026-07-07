package domain

import "time"

type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformWeb     Platform = "web"
)

func (p Platform) Valid() bool {
	switch p {
	case PlatformIOS, PlatformAndroid, PlatformWeb:
		return true
	}

	return false
}

type BrowserNotification struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterDevice struct {
	DeviceToken string
	Platform    Platform
}
