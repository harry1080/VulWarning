package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	conn *gorm.DB
)

// Warings 模型
type Warings struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Link     string `gorm:"type:text;unique"` // 情报链接
	Index    string `gorm:"type:varchar(255)"`
	Title    string `gorm:"type:varchar(255)"`
	From     string `gorm:"type:varchar(255)"` // 情报平台
	Desc     string `gorm:"type:text"`         // 情报描述/简介
	Time     time.Time
	CreateAt time.Time
	Send     bool
}

func initDatabase() {
	var err error
	conn, err = gorm.Open("mysql", DSN)
	if err != nil {
		log.Fatalln("Open MySQL Failed : ", err)
		return
	}

	conn.LogMode(false)

	conn.DB().SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	conn.DB().SetMaxOpenConns(100)                  //设置最大连接数
	conn.DB().SetMaxIdleConns(16)                   //设置闲置连接数

	// conn.DropTableIfExists(&Warings{})
	if !conn.HasTable(&Warings{}) {
		conn.CreateTable(&Warings{})
		conn.AutoMigrate(&Warings{})
	}
}

func addWarings(ws []*Warings) (err error) {
	tx := conn.Begin()
	for _, w := range ws {
		var out Warings
		if tx.First(&out, Warings{Link: w.Link}).RecordNotFound() {
			w.Send = true
			if err = tx.Create(w).Error; err != nil {
				logger.Errorln(err)
			}
			text := fmt.Sprintf(`%s
Time : %v
Url  : %s  
From : %s  `, w.Desc, w.Time, w.Link, w.From)
			logger.Debugln(text)
			pushFeishuMessage(w.Title, text)
		}
	}
	tx.Commit()
	return nil
}

// CreateWarings 添加记录
func CreateWarings(title, link, from, _timeFormat, _time string) Warings {
	t, err := time.Parse(_timeFormat, _time)
	if err != nil {
		log.Println(err.Error())
		t = time.Now()
	}
	var w Warings
	conn.FirstOrCreate(&w, Warings{
		Title:    title,
		Link:     link,
		From:     from,
		Time:     t,
		CreateAt: time.Now(),
	})
	conn.Save(&w)
	return w
}
