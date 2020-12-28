package friend

import (
	"fmt"
	"time"

	"github.com/lz1998/ecust_im/model"
	"gorm.io/gorm/clause"
)

type EcustFriend struct {
	ID        int64     `gorm:"column:id" json:"id" form:"id"`
	UserA     int64     `gorm:"column:user_a" json:"user_a" form:"user_a"`
	UserB     int64     `gorm:"column:user_b" json:"user_b" form:"user_b"`
	Status    int64     `gorm:"column:status" json:"status" form:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
}

// 添加/删除好友
func SaveFriend(f *EcustFriend) (*EcustFriend, error) {
	if f.UserA == f.UserB {
		return nil, fmt.Errorf("userA==userB")
	}
	if f.UserA > f.UserB {
		f.UserA, f.UserB = f.UserB, f.UserA
	}
	friend := &EcustFriend{
		UserA:  f.UserA,
		UserB:  f.UserB,
		Status: f.Status,
	}

	if err := model.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_a"}, {Name: "user_b"}},
		DoUpdates: clause.AssignmentColumns([]string{"status"}),
	}).Create(friend).Error; err != nil {
		return nil, err
	}

	return friend, nil
}

// 查询好友列表
func ListFriend(userId int64) ([]int64, error) {
	var friends []*EcustFriend
	q := model.Db.Model(&EcustFriend{})
	q = q.Where("status = 1")
	q = q.Where("user_a = ?",userId).Or("user_b = ?",userId)
	if err := q.Find(&friends).Error; err != nil {
		return nil, err
	}

	friendIds := make([]int64, 0)
	for _, friend := range friends {
		friendIds = append(friendIds, func() int64 {
			if friend.UserA != userId {
				return friend.UserA
			} else {
				return friend.UserB
			}
		}())
	}
	return friendIds, nil
}

func IsFriend(userA int64, userB int64) bool {
	if userA == userB {
		return false
	}
	if userA > userB {
		userA, userB = userB, userA
	}
	q := model.Db.Model(&EcustFriend{})
	q = q.Where("status = 1")
	q = q.Where("user_a = ?", userA)
	q = q.Where("user_b = ?", userB)
	ecustFriend := &EcustFriend{}
	if err := q.First(ecustFriend).Error; err != nil {
		return false
	} else {
		return true
	}
}
