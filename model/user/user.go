package user

import (
	"time"

	"github.com/lz1998/ecust_im/model"
)

type EcustUser struct {
	UserId    int64     `gorm:"column:user_id;primary_key;auto_increment;not_null" json:"user_id" form:"user_id"`
	Password  string    `gorm:"column:password" json:"password" form:"password"`
	Nickname  string    `gorm:"column:nickname" json:"nickname" form:"nickname"`
	Status    int64     `gorm:"column:status" json:"status" form:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
}

//// 自动创建表
//func init() {
//	if err := model.Db.AutoMigrate(&EcustUser{}); err != nil {
//		panic(err)
//	}
//}

func CreateUser(u *EcustUser) (*EcustUser, error) {
	user := &EcustUser{
		Password: u.Password,
		Nickname: u.Nickname,
		Status:   0,
	}
	if err := model.Db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func ListUser(userIds []int64) ([]*EcustUser, error) {
	var users []*EcustUser

	q := model.Db.Model(&EcustUser{})
	q = q.Where("status = 0")
	q = q.Where("user_id in ?", userIds)

	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(userId int64) (*EcustUser, error) {
	user := &EcustUser{}
	q := model.Db.Model(&EcustUser{})
	q = q.Where("status = 0")
	q = q.Where("user_id = ?", userId)
	if err := q.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(users []*EcustUser) error {
	for _, user := range users {
		if err := model.Db.Model(user).Updates(user).Error; err != nil {
			return err
		}
	}
	return nil
}
