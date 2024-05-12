package mapper

import (
	"database/sql"
	"github.com/jimu-server/model"
)

type SettingMapper struct {
	SettingsList       func(any) ([]*model.AppSetting, error)
	GetUserInfoSetting func() (*model.AppSetting, error)
	UpdateSetting      func(any) error
	AddSetting         func(any, *sql.Tx) error
}
