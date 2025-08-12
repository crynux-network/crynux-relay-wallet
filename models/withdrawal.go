package models

type Withdrawal struct {
	ID                uint64 `json:"id" gorm:"primarykey"`
	CreatedAt         uint64 `json:"created_at"`
	UpdatedAt         uint64 `json:"updated_at"`
	RequestID         uint64 `json:"request_id"`
	Blockchain        string `json:"blockchain" gorm:"size:255"`
	Address           string `json:"address" gorm:"size:255"`
	BeneficialAddress string `json:"beneficial_address" gorm:"size:255"`
	Amount            BigInt `json:"amount" gorm:"type:string;size:255"`
	Timestamp         int64  `json:"timestamp"`
	Nonce             uint64 `json:"nonce" gorm:"uniqueIndex"`
	Status            string `json:"status" gorm:"index"`
	ErrorMessage      string `json:"error_message" gorm:"size:255"`
	TxHash            string `json:"tx_hash" gorm:"size:255"`
	TxStatus          string `json:"tx_status" gorm:"size:255"`
}
