package setting

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"github.com/jimu-server/logger"
	"github.com/jimu-server/model"
	"github.com/jimu-server/redis/cache"
	"github.com/jimu-server/redis/redisUtil"
	"github.com/jimu-server/util/treeutils/tree"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
)

/*
系统内嵌初始化的设置项
*/

// Settings 本地配置列表 初始化配置模板
//
//go:embed setting_template/*-template.json
var Settings embed.FS

// GetUseSetting
// 获取用户设置项
// @param id 用户id
// @param set 设置项
func GetUseSetting[T any](id string, set string) (T, error) {
	var err error
	var data []tree.AnyNode
	var value T
	var get string
	if get, err = redisUtil.Get(fmt.Sprintf("%s:%s", USER_SETTING, id)); err != nil {
		return value, err
	}
	if err = jsoniter.Unmarshal([]byte(get), &data); err != nil {
		return value, err
	}
	var setItem map[string]any
	for i := range data {
		item := data[i]
		setItem = item.Entity.(map[string]any)
		settingName := setItem["name"].(string)
		if settingName == set {
			setValue := setItem["setting"].(string)
			if err = jsoniter.Unmarshal([]byte(setValue), &value); err != nil {
				return value, err
			}
		}

	}
	return value, err
}

func QueryUserSetting(id string) []tree.AnyNode {
	key := fmt.Sprintf("%s:%s", USER_SETTING, id)
	var err error
	var app_set string
	var data []tree.AnyNode
	if app_set, err = redisUtil.Get(key); err != nil && !errors.Is(err, redis.Nil) {
		logger.Logger.Error(err.Error())
		return nil
	}
	if err = jsoniter.Unmarshal([]byte(app_set), &data); err != nil {
		logger.Logger.Error(err.Error())
		return nil
	}
	return data
}

func UpdateUserSetting(id string, value any) error {
	var err error
	var settingValue string
	switch v := value.(type) {
	case string:
		settingValue = v
	default:
		if settingValue, err = jsoniter.MarshalToString(value); err != nil {
			return err
		}
	}
	if err = redisUtil.Del(fmt.Sprintf("%s:%s", USER_SETTING, id)); err != nil {
		return err
	}
	if err = redisUtil.SetEx(fmt.Sprintf("%s:%s", USER_SETTING, id), settingValue, cache.SettingCacheTime); err != nil {
		return err
	}
	return nil
}

// GetSettingTemplate 获取系统设置模板
func GetSettingTemplate() ([]model.AppSetting, error) {
	dir, err := Settings.ReadDir("setting_template")
	if err != nil {
		return nil, err
	}
	var buf []byte
	var arr []model.AppSetting
	for i := range dir {
		if buf, err = Settings.ReadFile(fmt.Sprintf("setting_template/%s", dir[i].Name())); err != nil {
			return nil, err
		}
		var data map[string]any
		if err = jsoniter.Unmarshal(buf, &data); err != nil {
			return nil, err
		}
		data["setting"], _ = jsoniter.MarshalToString(data["setting"])

		arr = append(arr, model.AppSetting{
			Name:    data["name"].(string),
			Pid:     data["pid"].(string),
			Value:   data["value"].(string),
			ToolId:  data["toolId"].(string),
			Setting: data["setting"].(string),
		})
	}
	return arr, nil
}
