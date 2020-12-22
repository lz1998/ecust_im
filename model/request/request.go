package request

import (
	"time"

	"github.com/lz1998/ecust_im/model"
	"github.com/lz1998/ecust_im/util"
)

const (
	TFriend = 0
	TGroup = 1
)

type EcustRequest struct {
	ReqId     int64     `gorm:"column:req_id" json:"req_id" form:"req_id"`
	ReqType   int64     `gorm:"column:req_type" json:"req_type" form:"req_type"`
	FromId    int64     `gorm:"column:from_id" json:"from_id" form:"from_id"`
	ToId      int64     `gorm:"column:to_id" json:"to_id" form:"to_id"`
	Status    int64     `gorm:"column:status" json:"status" form:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
}

func CreateRequest(req *EcustRequest) (*EcustRequest, error) {
	request := &EcustRequest{
		ReqId:   util.GenerateId(),
		ReqType: req.ReqType,
		FromId:  req.FromId,
		ToId:    req.ToId,
		Status:  0,
	}
	if err := model.Db.Create(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func UpdateRequest(req *EcustRequest) error {
	request, err := GetRequest(req.ReqId)
	if err != nil {
		return err
	}

	request.Status = req.Status
	return model.Db.Save(request).Error
}

func GetRequest(reqId int64) (*EcustRequest, error) {
	request := &EcustRequest{}
	q := model.Db.Model(&EcustRequest{})
	q = q.Where("req_id = ?", reqId)
	if err := q.First(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}
