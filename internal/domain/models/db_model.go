package models

type DbBaseModel struct {
	Id          string `json:"id" gorm:"primaryKey;size:32"`
	DataVersion int64  `json:"data_version"`
	DataStatus  int    `json:"data_status"`
	CreateTime  int64  `json:"create_time" gorm:"autoCreateTime:milli"`
	UpdateTime  int64  `json:"update_time" gorm:"autoUpdateTime:milli"`
}

func NewDbBaseModel(id string) DbBaseModel {
	return DbBaseModel{
		Id:          id,
		DataVersion: 1,
		DataStatus:  1,
	}
}

type TenantBaseModel struct {
	TenantId string `json:"tenant_id" gorm:"primaryKey;size:32"`
	DbBaseModel
}

func NewTenantBaseModel(tenantId string, id string) TenantBaseModel {
	return TenantBaseModel{
		TenantId:    tenantId,
		DbBaseModel: NewDbBaseModel(id),
	}
}
