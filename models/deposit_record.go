package models

import "gorm.io/gorm"

type DepositRecord struct {
	gorm.Model
	Network             string `json:"network" gorm:"type:string;size:191;not null;uniqueIndex:idx_deposit_record_network_tx_hash"`
	TxHash              string `json:"tx_hash" gorm:"type:string;size:191;not null;uniqueIndex:idx_deposit_record_network_tx_hash"`
	DepositAddress      string `json:"deposit_address" gorm:"type:string;size:191;not null"`
	FromAddress         string `json:"from_address" gorm:"type:string;size:191;not null;index"`
	Amount              BigInt `json:"amount" gorm:"type:string;size:255;not null"`
	RelayAccountEventID uint   `json:"relay_account_event_id" gorm:"not null;index"`
}

func (DepositRecord) TableName() string {
	return "deposit_records"
}
