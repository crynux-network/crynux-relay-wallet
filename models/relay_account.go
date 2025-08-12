package models

type RelayAccount struct {
	ID                uint64 `json:"id" gorm:"primarykey"`
	CreatedAt         uint64 `json:"created_at"`
	UpdatedAt         uint64 `json:"updated_at"`
	Address           string `json:"address" gorm:"uniqueIndex"`
	Balance           BigInt `json:"balance" gorm:"type:string;size:255"`
	BeneficialAddress string `json:"beneficial_address" gorm:"size:255"`
}
