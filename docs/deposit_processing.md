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
5. Select validation rules by configured network token type.
6. For native-token networks, query raw transaction transfer fields by `tx_hash` using `eth_getTransactionByHash`.
7. For native-token networks, verify raw transaction `to` equals configured `relay.deposit_address`.
8. For native-token networks, verify raw transaction `from` equals event `address`.
9. For native-token networks, verify raw transaction `value` equals event `amount`.
10. For native-token networks, verify raw transaction `input` is `0x`, so only ordinary native transfers are accepted as deposits.
11. For ERC20-token networks, verify the receipt contains a `Transfer(address,address,uint256)` log emitted by the configured token contract.
12. For ERC20-token networks, verify the indexed `to` address equals configured `relay.deposit_address`.
13. For ERC20-token networks, verify the indexed `from` address equals event `address`.
14. For ERC20-token networks, verify the indexed `from` address is not the zero address.
15. For ERC20-token networks, verify the transfer amount equals event `amount`.
16. Query the transaction block header and verify transaction age does not exceed configured `tasks.sync_task_fee_logs.deposit_max_age_seconds`.
17. Persist with unique key `(network, tx_hash)` and fail-fast on duplicate via database constraint.

If all validations pass, the deposit is marked as accepted for persistence.

The wallet MUST NOT decode the full transaction into a typed Ethereum transaction for deposit validation. For native deposits, the wallet MUST reject the deposit if raw transaction fields are unavailable, malformed, or inconsistent with the event log. For ERC20 deposits, the wallet MUST reject the deposit if receipt log fields are unavailable, malformed, or inconsistent with the event log.

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
