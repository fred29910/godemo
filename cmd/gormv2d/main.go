package main

import (
	"context"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	dsn := "root:12345#lxikm@tcp(localhost:3307)/dbv?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var t uint8 = 1
	var name string = "test"
	var desc string = "test"
	var status uint8 = 0
	var day uint32 = 1
	var startTime time.Time = time.Now()
	err = UpsertRewardConfig(db, context.Background(), &RewardEditCfgReq{
		ID:        nil,
		Type:      &t,
		Name:      &name,
		Desc:      &desc,
		StartTime: &startTime,
		Status:    &status,
		Day:       &day,
		CurrencyMap: map[int64]int64{
			1: 100,
			2: 200,
			3: 300,
		},
	})

	if err != nil {
		log.Println("UpsertRewardConfig err ", err)
		return
	}

	list, err := GetRewardConfig(db, context.Background(), &RewardConfigList{
		Type:     nil,
		Status:   nil,
		Page:     1,
		PageSize: 20,
	})
	if err != nil {
		log.Println("GetRewardConfig err ", err)
		return
	}
	log.Println("list ", list)
}

type RewardCfg struct {
	ID   int64 `json:"id" gorm:"column:id;primary_key"`
	Type uint8 `json:"type" gorm:"column:type"`

	Name         string    `json:"name" gorm:"column:name"`
	Desc         string    `json:"desc" gorm:"column:desc"`
	StartTime    time.Time `json:"start_date" gorm:"column:start_date"`
	Status       uint8     `json:"status" gorm:"column:status"` // 0 active 1 inactive
	Day          uint32    `json:"day" gorm:"column:day"`
	ClaimedCount int       `json:"claimed_count" gorm:"column:claimed_count"`
}

func (b *RewardCfg) TableName() string {
	return "rewards_cfg"
}

type RewardCurrencyCfg struct {
	ID       int64 `json:"id" gorm:"column:id;primary_key"`
	Currency int64 `json:"currency" gorm:"column:currency"`

	CfgID int64 `json:"cfg_id" gorm:"column:cfg_id"`

	Amount int64 `json:"amount" gorm:"column:amount"`
}

func (b *RewardCurrencyCfg) TableName() string {
	return "rewards_currency_cfg"
}

type RewardConfigList struct {
	Type *uint8 `json:"type,omitempty"`

	Status *uint8 `json:"status,omitempty"`

	Page     uint64 `json:"page,omitempty"`
	PageSize uint64 `json:"page_size,omitempty"`
}

type RewardConfigListResp struct {
	Total int64 `json:"total"`

	List []*RewardCfg `json:"list"`
}

func GetRewardConfig(db *gorm.DB, ctx context.Context, params *RewardConfigList) (*RewardConfigListResp, error) {

	page := params.Page
	pagesize := params.PageSize

	if page <= 0 {
		page = 1
	}

	if pagesize <= 0 {
		pagesize = 20
	}

	start := (page - 1) * pagesize
	// end := pagesize

	var list []*RewardCfg
	query := db.Model(&RewardCfg{})
	if params.Type != nil {
		query = query.Where("type = ?", *params.Type)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	var count int64

	err := query.Count(&count).Error
	if err != nil {
		return nil, err
	}
	err = query.Limit(pagesize).Offset(start).Find(&list).Error

	if err != nil {
		return nil, err
	}

	return &RewardConfigListResp{
		Total: count,
		List:  list,
	}, nil
}

type RewardCurrencyInfoResp struct {
	RewardCfg
	CurrencyMap map[int64]int64 `json:"currency_map"`
}

func GetRewardConfigInfo(db *gorm.DB, ctx context.Context, id int64) (*RewardCurrencyInfoResp, error) {
	var data RewardCfg
	err := db.Model(&RewardCfg{}).Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	var list []*RewardCurrencyCfg

	err = db.Model(&RewardCurrencyCfg{}).Where("cfg_id = ?", id).Find(&list).Error

	if err != nil {
		return nil, err
	}
	var rst RewardCurrencyInfoResp = RewardCurrencyInfoResp{
		RewardCfg: data,
	}
	rst.CurrencyMap = make(map[int64]int64)
	for _, v := range list {
		rst.CurrencyMap[v.Currency] = v.Amount
	}

	return &rst, nil
}

type RewardEditCfgReq struct {
	ID *int64 `json:"id" gorm:"column:id;primary_key"` // just update

	Type *uint8 `json:"type" gorm:"column:type"` // just create

	Name      *string    `json:"name" gorm:"column:name"`
	Desc      *string    `json:"desc" gorm:"column:desc"`
	StartTime *time.Time `json:"start_date" gorm:"column:start_date"`
	Status    *uint8     `json:"status" gorm:"column:status"` // 0 active 1 inactive
	Day       *uint32    `json:"day" gorm:"column:day"`

	CurrencyMap map[int64]int64 `json:"currency_map"`
}

func UpsertRewardConfig(db *gorm.DB, ctx context.Context, data *RewardEditCfgReq) error {
	if data.ID == nil {
		tx := db.Begin()
		var err error
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
		dataSave := RewardCfg{
			Name:      *data.Name,
			Desc:      *data.Desc,
			StartTime: *data.StartTime,
			Status:    *data.Status,
			Day:       *data.Day,
		}
		err = tx.Create(&dataSave).Error
		if err != nil {
			return err
		}
		for currency, amount := range data.CurrencyMap {
			err = tx.Create(&RewardCurrencyCfg{
				Currency: currency,
				Amount:   amount,
				CfgID:    dataSave.ID,
			}).Error
			if err != nil {
				return err
			}
		}
		return nil
	}
	upsMap := make(map[string]interface{})
	if data.Name != nil {
		upsMap["name"] = *data.Name
	}
	if data.Desc != nil {
		upsMap["desc"] = *data.Desc
	}
	if data.StartTime != nil {
		upsMap["start_date"] = *data.StartTime
	}
	if data.Status != nil {
		upsMap["status"] = *data.Status
	}
	if data.Day != nil {
		upsMap["day"] = *data.Day
	}

	tx := db.Begin()
	var err error
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Model(&RewardCfg{}).Where("id = ?", *data.ID).Updates(upsMap).Error
	if err != nil {
		return err
	}
	if data.CurrencyMap != nil {
		for currency, amount := range data.CurrencyMap {
			err = tx.Exec(`insert into rewards_currency_cfg (cfg_id, currency, amount)values (?, ?, ?) on duplicate key update amount = ?;`, *data.ID, currency, amount, amount).Error
			if err != nil {
				return err
			}
			// tx.Model(&RewardCurrencyCfg{}).Where("cfg_id = ? and currency = ?", *data.ID, currency).Updates(map[string]interface{}{"amount": amount})
		}
	}
	return err
}
