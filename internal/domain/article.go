package domain

import "time"

// Article 聚合根
type Article struct {
	ID      uint64
	Title   string
	Content string
	Ctime   time.Time
	Utime   time.Time

	Author  Author
}