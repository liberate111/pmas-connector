package model

import "time"

type KpiTag struct {
	Point       int       `gorm:"column:Point"`
	Group       string    `gorm:"column:Group"`
	Description string    `gorm:"column:Description"`
	KpiTag      string    `gorm:"column:KPI_Tag"`
	TagRaw      string    `gorm:"column:Tag_Raw"`
	Craft       string    `gorm:"column:Craft"`
	Status      string    `gorm:"column:Status"`
	Meter       string    `gorm:"column:Meter"`
	Location    string    `gorm:"column:Location"`
	Sfrom       string    `gorm:"column:SFrom"`
	TimeStamp   time.Time `gorm:"column:TimeStamp"`
}
