package models

import (
	"fmt"
	"math/big"
	"time"

	"gorm.io/gorm"
)

type VestingStatus int8

const (
	VestingStatusActive VestingStatus = iota
	VestingStatusCompleted
)

type VestingRecord struct {
	gorm.Model
	RelayVestingID uint          `json:"relay_vesting_id" gorm:"not null;uniqueIndex"`
	Address        string        `json:"address" gorm:"not null;index"`
	TotalAmount    BigInt        `json:"total_amount" gorm:"not null"`
	ReleasedAmount BigInt        `json:"released_amount" gorm:"not null"`
	StartTime      time.Time     `json:"start_time" gorm:"not null;index"`
	DurationDays   uint          `json:"duration_days" gorm:"not null"`
	Source         string        `json:"source" gorm:"not null;size:64;index"`
	ExternalID     string        `json:"external_id" gorm:"not null;size:128;index"`
	AdminSignature string        `json:"admin_signature" gorm:"not null;size:255"`
	Status         VestingStatus `json:"status" gorm:"not null;default:0;index"`
}

type VestingSignPayload struct {
	Address      string `json:"address"`
	TotalAmount  string `json:"total_amount"`
	StartTime    int64  `json:"start_time"`
	DurationDays uint   `json:"duration_days"`
	Source       string `json:"source"`
	ExternalID   string `json:"external_id"`
}

func BuildVestingSignMessage(payload VestingSignPayload) string {
	return fmt.Sprintf(
		"Crynux Relay Vesting\nAddress: %s\nTotalAmount: %s\nStartTime: %d\nDurationDays: %d\nSource: %s\nExternalID: %s",
		payload.Address,
		payload.TotalAmount,
		payload.StartTime,
		payload.DurationDays,
		payload.Source,
		payload.ExternalID,
	)
}

func ComputeVestingShouldReleased(totalAmount *big.Int, startTime time.Time, durationDays uint, now time.Time) *big.Int {
	if totalAmount == nil || totalAmount.Sign() <= 0 || durationDays == 0 {
		return big.NewInt(0)
	}

	startUTC := startTime.UTC()
	nowUTC := now.UTC()
	if nowUTC.Before(startUTC) {
		return big.NewInt(0)
	}

	startDay := time.Date(startUTC.Year(), startUTC.Month(), startUTC.Day(), 0, 0, 0, 0, time.UTC)
	nowDay := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)
	elapsedDays := uint(nowDay.Sub(startDay) / (24 * time.Hour))
	if elapsedDays >= durationDays {
		return new(big.Int).Set(totalAmount)
	}

	shouldReleased := big.NewInt(0).Mul(totalAmount, big.NewInt(0).SetUint64(uint64(elapsedDays)))
	shouldReleased.Div(shouldReleased, big.NewInt(0).SetUint64(uint64(durationDays)))
	return shouldReleased
}
