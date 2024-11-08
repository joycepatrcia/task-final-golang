package model

import (
    "time"
)

type Transaction struct {
    TransactionID        int64     `gorm:"primaryKey;autoIncrement" json:"transaction_id"`
    TransactionCategoryID *int64    `json:"transaction_category_id"` // Adjusted for nullable foreign key
    AccountID            *int64    `json:"account_id"`
    FromAccountID        *int64    `json:"from_account_id"`
    ToAccountID          *int64    `json:"to_account_id"`
    Amount               int64     `json:"amount"`
    TransactionDate      time.Time `json:"transaction_date"`
}


// Menentukan nama tabel yang benar
func (Transaction) TableName() string {
    return "transaction" // Pastikan nama tabelnya sesuai
}

