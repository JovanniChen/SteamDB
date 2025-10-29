package Dao

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/JovanniChen/SteamDB/Steam/Model"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db     *sql.DB
	dbOnce sync.Once
	dbErr  error
	dbPath = "steam.db" // 默认数据库路径
)

func SetDBPath(path string) {
	dbPath = path
}

// ensureDB 确保数据库已初始化
func ensureDB() error {
	dbOnce.Do(func() {
		dbErr = initDB(dbPath)
	})
	return dbErr
}

// initDB 内部初始化函数
func initDB(path string) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %w", err)
	}

	// 创建表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS game_update_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		game_id INTEGER NOT NULL,
		unique_id TEXT NOT NULL,
		app_id INTEGER NOT NULL,
		start_time INTEGER NOT NULL,
		event_name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(game_id, unique_id)
	);

	CREATE INDEX IF NOT EXISTS idx_game_id ON game_update_events(game_id);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}

	fmt.Println("✅ [数据库初始化] 初始化完成")
	return nil
}

// InitDB 手动初始化数据库连接（保留向后兼容）
func InitDB(path string) error {
	dbPath = path
	return ensureDB()
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetLatestUpdateEvent 获取指定游戏的最新更新事件
func GetLatestUpdateEvent(gameID int) (*Model.UpdateEventInfo, error) {
	// 确保数据库已初始化
	if err := ensureDB(); err != nil {
		return nil, err
	}

	query := `
		SELECT unique_id, app_id, start_time, event_name
		FROM game_update_events
		WHERE game_id = ?
		ORDER BY start_time DESC
		LIMIT 1
	`

	var event Model.UpdateEventInfo
	err := db.QueryRow(query, gameID).Scan(
		&event.UniqueID,
		&event.AppID,
		&event.StartTime,
		&event.EventName,
	)

	if err == sql.ErrNoRows {
		return nil, nil // 没有记录，返回 nil
	}

	if err != nil {
		return nil, fmt.Errorf("查询最新事件失败: %w", err)
	}

	return &event, nil
}

// SaveUpdateEvent 保存更新事件到数据库
func SaveUpdateEvent(gameID int, event *Model.UpdateEventInfo) error {
	// 确保数据库已初始化
	if err := ensureDB(); err != nil {
		return err
	}

	insertSQL := `
		INSERT OR REPLACE INTO game_update_events (game_id, unique_id, app_id, start_time, event_name)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.Exec(insertSQL, gameID, event.UniqueID, event.AppID, event.StartTime, event.EventName)
	if err != nil {
		return fmt.Errorf("保存更新事件失败: %w", err)
	}

	return nil
}

// CheckAndSaveUpdateEvent 检查并保存更新事件
// 返回值：needUpdate - 是否有新的更新
func CheckAndSaveUpdateEvent(gameID int, newEvent *Model.UpdateEventInfo) (needUpdate bool, err error) {
	// 查询数据库中的最新事件
	latestEvent, err := GetLatestUpdateEvent(gameID)
	if err != nil {
		return false, err
	}

	// 如果数据库中没有记录，直接保存
	if latestEvent == nil {
		if err := SaveUpdateEvent(gameID, newEvent); err != nil {
			return false, err
		}
		return true, nil
	}

	// 对比 UniqueID（如果不一致说明有新的更新）
	if latestEvent.UniqueID != newEvent.UniqueID {
		if err := SaveUpdateEvent(gameID, newEvent); err != nil {
			return false, err
		}
		return true, nil
	}

	// UniqueID 一致，说明没有新更新
	return false, nil
}
