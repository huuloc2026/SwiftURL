package model

type URL struct {
	ID      int64  `db:"id"`
	Code    string `db:"code"`
	LongURL string `db:"long_url"`
	Clicks  int64  `db:"clicks"`
}
