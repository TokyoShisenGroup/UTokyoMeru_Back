package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Good struct {
	gorm.Model
	Title       string         `gorm:"not null"`
	Description string         `gorm:"not null"`
	Images      pq.StringArray `gorm:"type:text[]"`
	Price       float64        `gorm:"not null"`
	Views       uint           `gorm:"default:0"`
	Likes       uint           `gorm:"default:0"`
	IsInvisible bool           `gorm:"default:false"`
	IsDeleted   bool           `gorm:"default:false"`
	IsBought    bool           `gorm:"default:false"`
	Tags        pq.StringArray `gorm:"type:text[]"`
	SellerID    uint           `gorm:"not null"`
	BuyerID     uint           `gorm:"default:null"`
	Comments    []Comment      `gorm:"foreignKey:GoodID"`
	FavoUsers   []User         `gorm:"many2many:user_likes"`
	Seller      User           `gorm:"foreignKey:SellerID"`
	Buyer       User           `gorm:"foreignKey:BuyerID"`
}

type User struct {
	gorm.Model
	Name        string    `gorm:"not null"`
	MailAddress string    `gorm:"unique;not null"`
	Password    string    `gorm:"not null"`
	Avatar      string    `gorm:"type:text"`
	Sales       []Good    `gorm:"foreignKey:SellerID"`
	Buys        []Good    `gorm:"foreignKey:BuyerID"`
	FavoList    []Good    `gorm:"many2many:user_likes"`
	Comments    []Comment `gorm:"foreignKey:UserID"`
	IsDeleted   bool      `gorm:"not null;default:false"`
	IsBanned    bool      `gorm:"not null;default:false"`
	UserClass   string    `gorm:"not null;default:user"`
	Gender      string    `gorm:"not null;default:Others"`
	Birthday    time.Time `gorm:"default:null"`
	PhoneNumber string    `gorm:"default:null"`
	MailCode    string    `gorm:"default:null"`
	Address     string    `gorm:"default:null"`
	Rating      float32   `gorm:"default:0"`
	RatingCount uint      `gorm:"default:0"`
	Bio         string    `gorm:"default:null"`
	Token       string
}

type Comment struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Content     string `gorm:"not null"`
	ReplyTo     uint   `gorm:"default:null"`
	IsInvisible bool   `gorm:"default:false"`
	IsDeleted   bool   `gorm:"default:false"`
	UserID      uint
	GoodID      uint
}

type Message struct {
	gorm.Model
	From      uint   `gorm:"not null" json:"from"`
	To        uint   `gorm:"not null" json:"to"`
	Content   string `gorm:"not null" json:"content"`
	ReplyTo   uint   `gorm:"default:null" json:"reply_to"`
	Status    string `gorm:"not null;default:Received" json:"status"` //Sent:已发送 Received:已接收 Read:已读 Deleted:已删除
	Sender    User   `gorm:"foreignKey:From" json:"sender"`
	Receiver  User   `gorm:"foreignKey:To" json:"receiver"`
}

type MessageQueue struct {
	gorm.Model
	TargetUserID uint   `json:"user_id"`
	MessageID    uint   `json:"message_id"`
	Message      Message `gorm:"foreignKey:MessageID"`
	Status       string `json:"status"` // Pending:待发送 Sent:已发送 Received:已接收 Read:已读 Deleted:已删除
}
