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
	RatingCount uint   `gorm:"default:0"`
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
    From    uint      `gorm:"not null"`
    To      uint      `gorm:"not null"`
    Content   string    `gorm:"not null"`
	ReplyTo   uint      `gorm:"default:null"`
	Type      string    `gorm:"not null"`
    IsRead    bool      `gorm:"not null;default:false"`
	IsDeleted bool      `gorm:"not null;default:false"`
	Sender    User      `gorm:"foreignKey:From"`
	Receiver  User      `gorm:"foreignKey:To"`
}
