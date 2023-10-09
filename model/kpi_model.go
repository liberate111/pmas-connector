package model

import "time"

type KpiTag struct {
	Point       int       `gorm:"column:POINT"`
	Group       string    `gorm:"column:Group"`
	Description string    `gorm:"column:DESCRIPTION"`
	KpiTag      string    `gorm:"column:KPI_TAG"`
	TagRaw      string    `gorm:"column:TAG_RAW"`
	Craft       string    `gorm:"column:CRAFT"`
	Status      string    `gorm:"column:STATUS"`
	Meter       string    `gorm:"column:METER"`
	Location    string    `gorm:"column:LOCATION"`
	Sfrom       string    `gorm:"column:SFROM"`
	TimeStamp   time.Time `gorm:"column:TimeStamp"`
}
