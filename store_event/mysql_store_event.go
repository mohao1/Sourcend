package store_event

import (
	"Sourcend/action"
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQLAction MySQL的存储结构和Action通用结构进行分离
type MySQLAction struct {
	id           int64        // 主键
	ActionName   string       // Action的名称
	RootData     string       // 初始化的数据 - json数据
	StartEventID int64        // 起始版本ID
	Events       []MySQLEvent // 指令版本记录 DATA:Event数据
	EventLen     int64        // 记录长度
}

// MySQLEvent MySQL的存储结构和Event通用结构进行分离
type MySQLEvent struct {
	CommandID  string            // CommandID - CommandID指令
	MutationID string            // MutationID - MutationID指令
	Event      string            // Event - 修改数据 - 序列化成JSON字符串的结构
	Params     map[string]string // Params - 其他扩展数据
}

// MySQLStore MYSQL的StoreEvent的实现
type MySQLStore struct {
	config      MySQLConfig
	gormDB      gorm.DB
	rooInitData string
}

func NewMySQLStore(config MySQLConfig, rooInitData string) *MySQLStore {
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("gorm open err:", err)
		return nil
	}

	store := &MySQLStore{
		config: config,
		gormDB: *db,
	}

	// 创建数据
	// 设置表名 - 更新创建数据库的结构
	err = db.Table(config.ActionTable).AutoMigrate(&action.Action{})
	if err != nil {
		fmt.Println("gorm migrate err:", err)
		return nil
	}

	return store
}

func (m *MySQLStore) Handler(ctx context.Context, data StoreEventInfo) error {
	dbTable := m.gormDB.Table(m.config.ActionTable)

	isKey := m.config.IsActionKey
	if isKey && m.config.ActionKey != "" {
		actionKey := data.Params[m.config.ActionKey]
		actionData := MySQLAction{}
		// 通过长度计算
		err := dbTable.Where("action_name = ?", actionKey).Order("start_event_id DESC").First(&actionData).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				actionData = MySQLAction{
					ActionName:   actionKey,
					RootData:     m.rooInitData,
					StartEventID: 0,
					Events:       []MySQLEvent{},
					EventLen:     0,
				}
			} else {
				fmt.Println("first action err:", err)
				return err
			}
		}

		if actionData.StartEventID+actionData.EventLen >= m.config.ActionMaxLen {
			// 获取新的rootData
			// 进行回放
			rootData := ""
			// 进行新的数据组装
			actionData = MySQLAction{
				ActionName:   actionKey,
				RootData:     rootData,
				StartEventID: actionData.StartEventID + actionData.EventLen,
				Events:       []MySQLEvent{},
				EventLen:     0,
			}
		}

		actionData.EventLen++
		actionData.Events = append(actionData.Events, MySQLEvent{
			CommandID:  data.CommandID,
			MutationID: data.MutationID,
			Event:      data.Event,
			Params:     data.Params,
		})

		// 存储
		dbTable.Save(actionData)

	} else {
		// 使用最小版本号来进行作为主键存储
		actionKey := "0"
	}

	//sqlAction := MySQLAction{
	//	ActionName:   "",
	//	RootData:     "",
	//	StartEventID: 0,
	//	Events:       nil,
	//	EventLen:     0,
	//}

	return nil
}
