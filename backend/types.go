package main

import "time"

type PubKey struct {
	Modulus  string `json:"modulus"`
	Exponent string `json:"exponent"`
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserId   int64  `json:"userId"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	Role     int64  `json:"role"`
	BuddyNum string `json:"buddynum"`
}

type TaskInfo struct {
	User                string `json:"user" form:"user"`
	Username            string `json:"username" form:"username"`
	Password            string `json:"password" form:"password"`
	CaptchaApi          string `json:"captcha_api" form:"captcha_api"`
	TaskId              string `json:"task_id" form:"task_id"`
	BuddyId             string `json:"buddy_id" form:"buddy_id"`
	BuddyUserId         string `json:"buddy_user_id" form:"buddy_user_id"`
	BuddyNum            string `json:"buddy_num" form:"buddy_num"`
	VenueSiteId         string `json:"venue_site_id" form:"venue_site_id"`
	ReservationDate     string `json:"reservation_date" form:"reservation_date"`
	ReservationTime     string `json:"reservation_time" form:"reservation_time"`
	UserPhone           string `json:"user_phone" form:"user_phone"`
	CaptchaVerification string `json:"captcha_verification" form:"captcha_verification"`
	SiteName            string `json:"site_name" form:"site_name"`
	SiteId              string `json:"site_id" form:"site_id"`
	TimeId              string `json:"time_id" form:"time_id"`
	OrderToken          string `json:"order_token" form:"order_token"`
	TradeNo             string `json:"trade_no" form:"trade_no"`
	OrderId             string `json:"order_id" form:"order_id"`
	InstantReservation  bool   `json:"instant_reservation" form:"instant_reservation"`
}

// 定义 TaskInfo 结构体，映射到数据库中的 task_info 表
type TaskInfoDb struct {
	ID                  uint      `gorm:"primaryKey"`
	User                string    `gorm:"not null"`                // 必填
	Username            string    `gorm:"not null"`                // 必填
	Password            string    `gorm:"not null"`                // 必填
	UserPhone           string    `gorm:"not null"`                // 必填
	CaptchaAPI          string    `gorm:"not null"`                // 必填
	BuddyUserID         string    `gorm:"not null"`                // 必填
	BuddyNum            string    `gorm:"not null"`                // 必填
	VenueSiteID         string    `gorm:"not null"`                // 必填
	ReservationDate     string    `gorm:"not null"`                // 必填
	ReservationTime     string    `gorm:"not null"`                // 必填
	SiteName            string    `gorm:"not null"`                // 必填
	TaskID              string    `gorm:"not null"`                // 必填
	CreateTime          time.Time `gorm:"autoCreateTime;not null"` // 必填
	IsFinished          bool      `gorm:"not null"`                // 必填
	InstantReservation  bool      `gorm:"not null"`                // 必填
	BuddyId             string
	SiteId              string
	TimeId              string
	OrderToken          string
	CaptchaVerification string
	TradeNo             string
	OrderId             string
	Autocancel          bool // 是否自动取消
	ReservationStatus   bool // 正常和已取消
}

type UserInfoDb struct {
	ID        uint   `gorm:"primaryKey"`
	User      string `gorm:"not null"`
	Lable     string `gorm:"not null"`
	Username  string `gorm:"not null;unique"` // 必填且唯一
	Password  string `gorm:"not null"`        // 必填
	Phone     string
	Name      string
	Sex       string
	Dept      string
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	IsDelete  bool      `gorm:"not null;default:false"` // 必填
}

type UserDb struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Username   string    `gorm:"not null;unique" json:"username"` // 必填且唯一
	Password   string    `gorm:"not null" json:"password"`
	CaptchaAPI string    `gorm:"not null" json:"captcha_api"`
	IsAdmin    bool      `gorm:"not null;default:false" json:"is_admin"`    // 必填
	IsDelete   bool      `gorm:"not null;default:false" json:"is_delete"`   // 必填
	CreatedAt  time.Time `gorm:"autoCreateTime;not null" json:"created_at"` // 自动设置创建时间
}

type ChangePassword struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

type ReservationAccount struct {
	User     string `json:"user"`
	Lable    string `json:"lable"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TyysAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
