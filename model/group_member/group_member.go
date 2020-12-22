package group_member

import (
	"time"

	"github.com/lz1998/ecust_im/model"
	"gorm.io/gorm/clause"
)

type EcustGroupMember struct {
	ID        int64     `gorm:"column:id" json:"id" form:"id"`
	GroupId   int64     `gorm:"column:group_id" json:"group_id" form:"group_id"`
	UserId    int64     `gorm:"column:user_id" json:"user_id" form:"user_id"`
	Status    int64     `gorm:"column:status" json:"status" form:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
}

// 创建/修改 群成员状态
func SaveGroupMember(m *EcustGroupMember) (*EcustGroupMember, error) {
	member := &EcustGroupMember{
		GroupId: m.GroupId,
		UserId:  m.UserId,
		Status:  m.Status,
	}
	if err := model.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "group_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status"}),
	}).Create(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}

// 查询群内成员
func ListGroupMember(groupId int64) ([]int64, error) {
	var members []*EcustGroupMember

	q := model.Db.Model(&EcustGroupMember{})
	q = q.Where("status = 1")
	q = q.Where("group_id = ?", groupId)
	if err := q.Find(&members).Error; err != nil {
		return nil, err
	}
	userIds := make([]int64, 0)
	for _, member := range members {
		userIds = append(userIds, member.UserId)
	}
	return userIds, nil
}

// 查询一个用户加了哪些群
func ListGroup(userId int64) ([]int64, error) {
	var members []*EcustGroupMember

	q := model.Db.Model(&EcustGroupMember{})
	q = q.Where("status = 1")
	q = q.Where("user_id = ?", userId)
	if err := q.Find(&members).Error; err != nil {
		return nil, err
	}
	groupIds := make([]int64, 0)
	for _, member := range members {
		groupIds = append(groupIds, member.GroupId)
	}
	return groupIds, nil
}

func IsInGroup(groupId int64, userId int64) bool {
	q := model.Db.Model(&EcustGroupMember{})
	q = q.Where("status = 1")
	q = q.Where("group_id = ?", groupId)
	q = q.Where("user_id = ?", userId)
	ecustGroupMember := &EcustGroupMember{}
	if err := q.First(ecustGroupMember).Error; err != nil {
		return false
	} else {
		return true
	}
}
