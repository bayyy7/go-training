package model

type Transaction struct {
	TransactionID         int64 `json:"transaction_id" gorm:"primaryKey;autoIncrement;<-:false"`
	TransactionCategoryID int64 `json:"transaction_category_id" gorm:"foreignKey:TransactionCategoryID;autoIncrement;<-:false"`
	AccountID             int64 `json:"account_id"`
	FromAccountID         int64 `json:"from_account_id"`
	ToAccountID           int64 `json:"to_account_id"`
	Amount                int64 `json:"amount"`
	TransactionDate       int64 `json:"transaction_date"`
}

func (Transaction) TableName() string {
	return "transaction"
}
