package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	// Conn 数据库链接
	Conn *gorm.DB
)

// InitDatabase - 初始化数据库连接
// func InitDatabase() {}

// Warings 模型
type Warings struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Title    string `gorm:"type:varchar(255);unique"`
	Link     string `gorm:"type:text;unique"`
	Time     time.Time
	CreateAt time.Time
	From     string `gorm:"type:varchar(255)"`
	Send     bool
}

// UpdateWaringsSend 变更发送状态
func UpdateWaringsSend(id uint) {
	var w Warings
	// Conn.Model(&Warings).Update("CreatedAt", time.Now())
	Conn.First(&w, id).Update("send", true)
}

// CreateWarings 添加记录
func CreateWarings(title, link, from, _timeFormat, _time string) Warings {
	t, err := time.Parse(_timeFormat, _time)
	if err != nil {
		log.Println(err.Error())
		t = time.Now()
	}
	var w Warings
	Conn.FirstOrCreate(&w, Warings{
		Title:    title,
		Link:     link,
		From:     from,
		Time:     t,
		CreateAt: time.Now(),
	})
	Conn.Save(&w)
	return w
}

func initDatabas() {
	if !Conn.HasTable(&Warings{}) {
		Conn.CreateTable(&Warings{})
		Conn.AutoMigrate(&Warings{})
	}
}

func init() {
	var err error
	Conn, err = gorm.Open("mysql", DSN)
	if err != nil {
		log.Fatalln("Open MySQL Failed : ", err)
		return
	}

	Conn.LogMode(false)

	Conn.DB().SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	Conn.DB().SetMaxOpenConns(100)                  //设置最大连接数
	Conn.DB().SetMaxIdleConns(16)                   //设置闲置连接数

	initDatabas()
}
