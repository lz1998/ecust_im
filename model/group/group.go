package group

import "time"

type EcustGroup struct {
	GroupId   int64     `gorm:"column:group_id" json:"group_id" form:"group_id"`
	GroupName string    `gorm:"column:group_name" json:"group_name" form:"group_name"`
	OwnerId   int64     `gorm:"column:owner_id" json:"owner_id" form:"owner_id"`
	Status    int64     `gorm:"column:status" json:"status" form:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
}

func CreateGroup
