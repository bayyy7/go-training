package model

type Auth struct {
	AuthID    int64  `json:"auth_id" gorm:"primaryKey;autoIncrement;<-:false"`
	AccountID int64  `json:"account_id" gorm:"autoIncrement;<-:false"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (Auth) TableName() string {
	return "auth"
}
