package tasks

import (
	"crynux_relay_wallet/relay_api"
	"errors"
	"testing"
)

func TestMergeTaskFeeLogsRejectsUnknownEventType(t *testing.T) {
	_, err := mergeTaskFeeLogs([]relay_api.TaskFeeLog{
		{
			Address: "0xabc",
			Amount:  "1",
			Type:    relay_api.TaskFeeLogType(100),
		},
	})
	if !errors.Is(err, ErrTaskFeeUnknownEventType) {
		t.Fatalf("expected ErrTaskFeeUnknownEventType, got %v", err)
	}
}

func TestMergeTaskFeeLogsHandlesVestingTypes(t *testing.T) {
	merged, err := mergeTaskFeeLogs([]relay_api.TaskFeeLog{
		{
			Address: "0xabc",
			Amount:  "0",
			Type:    relay_api.TaskFeeLogTypeVestingCreated,
		},
		{
			Address: "0xabc",
			Amount:  "25",
			Type:    relay_api.TaskFeeLogTypeVestingRelease,
		},
	})
	if err != nil {
		t.Fatalf("mergeTaskFeeLogs failed: %v", err)
	}
	if got := merged["0xabc"].String(); got != "25" {
		t.Fatalf("expected merged vesting release amount 25, got %s", got)
	}
}

func TestMergeTaskFeeLogsHandlesUserDelegation(t *testing.T) {
	merged, err := mergeTaskFeeLogs([]relay_api.TaskFeeLog{
		{
			Address: "0xabc",
			Amount:  "13",
			Type:    relay_api.TaskFeeLogTypeUserDelegation,
		},
	})
	if err != nil {
		t.Fatalf("mergeTaskFeeLogs failed: %v", err)
	}
	if got := merged["0xabc"].String(); got != "13" {
		t.Fatalf("expected merged user delegation amount 13, got %s", got)
	}
}

func TestParseVestingReleasePayloadOnlyRequiresVestingID(t *testing.T) {
	payload, err := parseVestingReleasePayload(relay_api.TaskFeeLog{
		Payload: `{"vesting_id":42}`,
	})
	if err != nil {
		t.Fatalf("parseVestingReleasePayload failed: %v", err)
	}
	if payload.VestingID != 42 {
		t.Fatalf("expected vesting id 42, got %d", payload.VestingID)
	}
}

func TestParseVestingReleasePayloadRejectsMissingVestingID(t *testing.T) {
	_, err := parseVestingReleasePayload(relay_api.TaskFeeLog{
		Payload: `{}`,
	})
	if !errors.Is(err, ErrTaskFeeVestingPayloadInvalid) {
		t.Fatalf("expected ErrTaskFeeVestingPayloadInvalid, got %v", err)
	}
}
