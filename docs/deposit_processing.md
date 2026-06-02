# Deposit Processing Specification

## Scope

This document defines how the Relay Wallet validates and applies `Deposit` relay account event logs.

## Input Data

Deposit logs are fetched from `/v1/relay_account/event_logs` as part of the task-fee log synchronization flow.

For `Deposit` logs (`type=3`), `payload` must be valid JSON and include:

- `tx_hash`
- `network`

## Processing Flow

For each synced batch, the wallet validates deposits before any balance update is committed:

1. Parse `amount` from the event log.
2. Parse deposit payload (`tx_hash`, `network`).
3. Query blockchain transaction receipt by `tx_hash`.
4. Verify receipt status is successful.
5. Query raw transaction transfer fields by `tx_hash` using `eth_getTransactionByHash`.
6. Verify raw transaction `to` equals configured `relay.deposit_address`.
7. Verify raw transaction `from` equals event `address`.
8. Verify raw transaction `value` equals event `amount`.
9. Verify raw transaction `input` is `0x`, so only ordinary native transfers are accepted as deposits.
10. Query the transaction block header and verify transaction age does not exceed configured `tasks.sync_task_fee_logs.deposit_max_age_seconds`.
11. Persist with unique key `(network, tx_hash)` and fail-fast on duplicate via database constraint.

If all validations pass, the deposit is marked as accepted for persistence.

The wallet MUST NOT decode the full transaction into a typed Ethereum transaction for deposit validation. The wallet MUST reject the deposit if raw transaction fields are unavailable, malformed, or inconsistent with the event log.

## Rejection Policy

Deposit validation follows fail-fast behavior:

- Any invalid deposit in the batch causes the sync attempt to fail.
- The task returns `TaskFeeError` and exits.
- Existing task-level alerting reports the failure for operator intervention.
- No balances are changed for that batch.
- The task-fee checkpoint is not advanced for that batch.

This behavior is intentional because abnormal deposits indicate either a serious bug or a potential attack.

## Persistence and Atomicity

Accepted deposits are recorded in `deposit_records` with at least:

- `network`
- `tx_hash`
- `deposit_address`
- `from_address`
- `amount`
- `relay_account_event_id`

The wallet persists accepted deposits, applies relay-account balance deltas, and updates `task_fee_checkpoints` in one database transaction.

If any step fails, the transaction is rolled back, so partial state is never committed.

## Idempotency Guarantee

`deposit_records` enforces a unique key on `(network, tx_hash)`.

This guarantees a previously accepted deposit transaction cannot be applied again, even if the same event appears in later sync attempts.
