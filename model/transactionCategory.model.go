package model

type TransactionCategory struct {
	TransactionCategoryID int64  `json:"transaction_category_id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name                  string `json:"name"`
}

func (TransactionCategory) TableName() string {
	return "transaction_category"
}
