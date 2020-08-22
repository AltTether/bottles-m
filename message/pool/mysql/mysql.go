package mysql

import (
	"fmt"
	"time"
	"math/rand"
	
	"github.com/jinzhu/gorm"
	"github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/message/pool"
	"github.com/message/utils"
)


type Pool struct {
	client *gorm.DB
}

type MessageRecord struct {
	Id        int64   `json:"id"`
	Text      *string `json:"message"`
 	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	DeletedAt *int64  `json:"deleted_at"`
}

type Message struct {
	text *string `json:"message"`
}

type MessageRecords []MessageRecord

type Config struct {
	Dsn       string
	ReconnDelay int
}

func NewMessageRecord(t *string) *MessageRecord {
	return &MessageRecord{
		Text:      t,
		CreatedAt: *timestamp(),
		UpdatedAt: *timestamp(),
	}
}

func NewMessage(t *string) *Message {
	return &Message{
		text: t,
	}
}

func (m *Message) Text() *string {
	return m.text
}

func (m *MessageRecord) SetUpdatedAt(t *int64) {
	m.UpdatedAt = *t
}

func (m *MessageRecord) SetDeletedAt(t *int64) {
	m.DeletedAt = t
}

func DefaultConfig() (*Config) {
	user := utils.GetEnvStr("MYSQL_ROOT_USER")
	passwd := utils.GetEnvStr("MYSQL_PASSWORD")
	net := utils.GetEnvStr("MYSQL_PROTOCOL")
	host := utils.GetEnvStr("MYSQL_HOST")
	port := utils.GetEnvStr("MYSQL_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	dbName := utils.GetEnvStr("MYSQL_DATABASE")

	dsn := (&mysql.Config{
		User: user,
		Passwd: passwd,
		Net: net,
		Addr: addr,
		DBName:	dbName,
		AllowNativePasswords: true,
		ParseTime: true,
	}).FormatDSN()

	reconnDelay := utils.GetEnvInt("DB_RECONNECTION_SEC")

	return &Config{
		Dsn: dsn,
		ReconnDelay: reconnDelay,
	}
}

func New() *Pool {
	config := DefaultConfig()
	return NewWithConfig(config)
}

func NewWithConfig(config *Config) *Pool {
	client, err := gorm.Open("mysql", config.Dsn)
	for err != nil {
		time.Sleep(time.Duration(config.ReconnDelay) * time.Second)
		client, err = gorm.Open("mysql", config.Dsn)
	}

	return &Pool{
		client: client,
	}
}

func (p *Pool) Get() (pool.Message, error){
	var records MessageRecords
	p.client.Table("messages").Find(&records, "deleted_at is null")
	n := len(records)
	if n == 0 {
		return &Message{}, fmt.Errorf("Not Found Unable Message")
	}

	idx := rand.Intn(n)
	record := &records[idx]
	
	record.SetUpdatedAt(timestamp())
	record.SetDeletedAt(timestamp())
	p.client.Table("messages").Save(&record)

	message := NewMessage(record.Text)
	return message, nil
}

func (p *Pool) Post(m pool.Message) (error) {
	r := NewMessageRecord(m.Text())
	p.client.Table("messages").Save(&r)
	return nil
}

func timestamp() *int64 {
	timestamp := time.Now().UTC().Unix()
	return &timestamp
}
