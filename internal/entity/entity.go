package entity

import "time"

type ShortURL struct {
	ID        int64      `db:"id"`
	ShortCode string     `db:"short_code"`
	LongURL   string     `db:"long_url"`
	CreatedAt time.Time  `db:"created_at"`
	ExpireAt  *time.Time `db:"expire_at"`
	Clicks    int64      `db:"clicks"`
}
type ClickLog struct {
	ShortCode  string
	ClickedAt  time.Time
	Referrer   string
	UserAgent  string
	DeviceType string
	OS         string
	Browser    string
	IPAddress  string
	Country    string
	City       string
}
