package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Email        *string   `gorm:"unique" json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Currency struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Code          string `gorm:"unique;not null" json:"code"`
	Name          string `gorm:"not null" json:"name"`
	Symbol        string `json:"symbol"`
	DecimalPlaces int    `gorm:"default:2" json:"decimal_places"`
}

type Account struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	Type        string    `gorm:"not null;check:type IN ('asset','liability','income','expense')" json:"type"`
	Balance     Amount    `gorm:"not null;default:0" json:"balance"` // in cents
	CurrencyID  uint      `gorm:"default:1" json:"currency_id"`
	Description *string   `json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Currency *Currency `gorm:"foreignKey:CurrencyID" json:"currency,omitempty"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	ParentID  *uint     `json:"parent_id"`
	Color     *string   `json:"color"`
	Icon      *string   `json:"icon"`
	Type      string    `gorm:"not null;check:type IN ('income','expense')" json:"type"`
	CreatedAt time.Time `json:"created_at"`

	Parent *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}

type Transaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	AccountID   uint      `gorm:"not null;index" json:"account_id"`
	ToAccountID *uint     `json:"to_account_id"` // for transfers
	CategoryID  *uint     `json:"category_id"`
	Amount      Amount    `gorm:"not null" json:"amount"` // always positive, in cents
	Type        string    `gorm:"not null;check:type IN ('income','expense','transfer')" json:"type"`
	Date        time.Time `gorm:"type:date;not null;index" json:"date"`
	Description string    `json:"description"`
	Notes       *string   `json:"notes"`
	Reconciled  bool      `gorm:"default:false" json:"reconciled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Account   *Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	ToAccount *Account  `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
	Category  *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags      []Tag     `gorm:"many2many:transaction_tags;" json:"tags,omitempty"`
}

type Budget struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	CategoryID uint      `gorm:"not null" json:"category_id"`
	Period     string    `gorm:"not null;index" json:"period"` // YYYY-MM
	Amount     Amount    `gorm:"not null" json:"amount"`       // in cents
	CreatedAt  time.Time `json:"created_at"`

	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	Color     *string   `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

// Many-to-many join table
type TransactionTag struct {
	TransactionID uint `gorm:"primaryKey"`
	TagID         uint `gorm:"primaryKey"`
}
