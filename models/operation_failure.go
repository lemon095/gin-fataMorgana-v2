package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONString 自定义JSON类型
type JSONString struct {
	Data []byte
}

// Value 实现driver.Valuer接口
func (j JSONString) Value() (driver.Value, error) {
	if j.Data == nil {
		return nil, nil
	}
	return string(j.Data), nil
}

// Scan 实现sql.Scanner接口
func (j *JSONString) Scan(value interface{}) error {
	if value == nil {
		j.Data = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		j.Data = v
	case string:
		j.Data = []byte(v)
	default:
		return json.Unmarshal([]byte(v.(string)), &j.Data)
	}
	return nil
}

// MarshalJSON 实现json.Marshaler接口
func (j JSONString) MarshalJSON() ([]byte, error) {
	if j.Data == nil {
		return []byte("null"), nil
	}
	return j.Data, nil
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (j *JSONString) UnmarshalJSON(data []byte) error {
	j.Data = data
	return nil
}

// OperationFailure 操作失败记录表
type OperationFailure struct {
	ID            uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	Uid           *string     `json:"uid" gorm:"size:8;index;comment:用户唯一ID（可为空，如注册时）"`
	OperationType string      `json:"operation_type" gorm:"not null;size:50;index;comment:操作类型"`
	RequestData   *JSONString `json:"request_data" gorm:"type:json;comment:请求数据（JSON格式）"`
	ResponseData  *JSONString `json:"response_data" gorm:"type:json;comment:响应数据（JSON格式）"`
	CreatedAt     time.Time   `json:"created_at" gorm:"autoCreateTime;index"`
}

// TableName 指定表名
func (OperationFailure) TableName() string {
	return "operation_failures"
}

// TableComment 表注释
func (OperationFailure) TableComment() string {
	return "操作失败记录表 - 记录用户操作失败信息"
}

// 操作类型常量
const (
	OperationTypeLogin          = "login"
	OperationTypeRegister       = "register"
	OperationTypeOrderCreate    = "order_create"
	OperationTypeWalletWithdraw = "wallet_withdraw"
	OperationTypeWalletRecharge = "wallet_recharge"
	OperationTypeBankCardBind   = "bank_card_bind"
	OperationTypeGroupBuyJoin   = "group_buy_join"
	OperationTypeSystemTask     = "system_task"
)

// OperationFailureResponse 操作失败记录响应
type OperationFailureResponse struct {
	ID            uint      `json:"id"`
	Uid           *string   `json:"uid"`
	OperationType string    `json:"operation_type"`
	RequestData   *string   `json:"request_data"`
	ResponseData  *string   `json:"response_data"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToResponse 转换为响应格式
func (f *OperationFailure) ToResponse() OperationFailureResponse {
	var requestDataStr, responseDataStr *string

	if f.RequestData != nil {
		reqStr := string(f.RequestData.Data)
		requestDataStr = &reqStr
	}

	if f.ResponseData != nil {
		respStr := string(f.ResponseData.Data)
		responseDataStr = &respStr
	}

	return OperationFailureResponse{
		ID:            f.ID,
		Uid:           f.Uid,
		OperationType: f.OperationType,
		RequestData:   requestDataStr,
		ResponseData:  responseDataStr,
		CreatedAt:     f.CreatedAt,
	}
}
