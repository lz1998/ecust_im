package group

import (
	"time"

	"github.com/lz1998/ecust_im/model"
)

type EcustGroup struct {
	GroupId   int64     `gorm:"column:group_id;primary_key;auto_increment;not_null" json:"group_id" form:"group_id"`
	GroupName string    `gorm:"column:group_name" json:"group_name" form:"group_name"`
	OwnerId   int64     `gorm:"column:owner_id" json:"owner_id" form:"owner_id"`
	Status    int64     `gorm:"column:status" json:"status" form:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
}

func CreateGroup(g *EcustGroup) (*EcustGroup, error) {
	group := &EcustGroup{
		GroupName: g.GroupName,
		OwnerId:   g.OwnerId,
		Status:    0,
	}

	if err := model.Db.Create(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

func ListGroup(groupIds []int64) ([]*EcustGroup, error) {
	var groups []*EcustGroup

	q := model.Db.Model(&EcustGroup{})
	q = q.Where("status = 0")
	q = q.Where("group_id in ?", groupIds)

	if err := q.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// 在handler层需要检测权限，这里允许修改ownerId（转让群）
func UpdateGroup(groups []*EcustGroup) error {
	for _, groups := range groups {
		if err := model.Db.Model(groups).Updates(groups).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetGroup(groupId int64) (*EcustGroup, error) {
	group := &EcustGroup{}
	if err := model.Db.Model(&EcustGroup{}).Where("group_id = ?", groupId).First(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}
