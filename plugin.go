package main

import (
	"log"
	"time"

	"github.com/gotify/plugin-api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UserLog represents a database model for user activity logs
type UserLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"index"`
	Message   string    `gorm:"type:text"`
	Timestamp time.Time `gorm:"autoCreateTime"`
}

// Plugin structure
type Plugin struct {
	plugin.Messenger
	db *gorm.DB
}

// Init initializes the plugin
func (p *Plugin) Init() error {
	var err error

	// Connect to SQLite database
	p.db, err = gorm.Open(sqlite.Open("user_activity.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto-migrate the UserLog model
	err = p.db.AutoMigrate(&UserLog{})
	if err != nil {
		return err
	}

	log.Println("User Activity Tracker plugin initialized")
	return nil
}

// SendMessage logs the user activity when a message is sent
func (p *Plugin) SendMessage(message plugin.Message) {
	userID := message.Extras["user_id"].(string)

	// Log the user activity
	userLog := UserLog{
		UserID:    userID,
		Message:   message.Message,
		Timestamp: time.Now(),
	}
	p.db.Create(&userLog)

	log.Printf("Message logged for user: %s at %s\n", userID, userLog.Timestamp)
}

// CheckUserActivity retrieves user logs for the last 24 hours
func (p *Plugin) CheckUserActivity(userID string) ([]UserLog, error) {
	var logs []UserLog
	err := p.db.Where("user_id = ? AND timestamp >= ?", userID, time.Now().Add(-24*time.Hour)).Find(&logs).Error
	return logs, err
}

// Close closes the plugin
func (p *Plugin) Close() {
	sqlDB, err := p.db.DB()
	if err == nil {
		sqlDB.Close()
	}
	log.Println("User Activity Tracker plugin closed")
}

func main() {
	plugin.Serve(&Plugin{})
}
