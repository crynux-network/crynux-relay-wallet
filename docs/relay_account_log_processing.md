# Relay Account Log Processing Specification

## Scope

This document defines how the Relay Wallet processes Relay Account logs and how withdrawal request ingestion is gated by log progress.

## Log Source and Schema

The wallet MUST fetch account logs from `/v1/relay_account/logs`.

Each log record MUST include:

- `id`
- `created_at`
- `address`
- `amount`
- `type`

## Log Type Handling Rules

For local `relay_accounts.balance`, the wallet MUST apply each log type as follows:

- `TaskIncome` (`type=0`): increase balance (`+amount`)
- `DaoTaskShare` (`type=1`): increase balance (`+amount`)
- `WithdrawalFeeIncome` (`type=2`): increase balance (`+amount`)
- `Deposit` (`type=3`): increase balance (`+amount`)
- `TaskPayment` (`type=4`): decrease balance (`-amount`)
- `TaskRefund` (`type=5`): increase balance (`+amount`)
- `Withdraw` (`type=6`): ignore in log sync
- `WithdrawRefund` (`type=7`): ignore in log sync

Under current implementation behavior, any type other than `TaskPayment`, `Withdraw`, and `WithdrawRefund` SHALL be merged as a positive delta.

## Validation Rules Before Balance Apply

For each fetched batch, the wallet MUST validate:

- `amount` MUST be parseable as an integer string.
- Per-log max amount threshold MUST be enforced for all types except `Deposit`.
- Per-address log count threshold MUST be enforced using only non-ignored logs.
- New-address count threshold MUST be enforced using only non-ignored logs.

If a batch contains only ignored log types (`Withdraw` and `WithdrawRefund`), validation SHALL pass and balance merge result SHALL be empty.

## Checkpoint Progress Rules

After a batch is accepted, the wallet MUST advance:

- `LatestTaskFeeLogID`
- `LatestTaskFeeLogTimestamp`

to the last fetched log in that batch, including batches where all logs are ignored for balance updates.

## Withdrawal Request Gate Rule

Each withdrawal request includes `relay_account_event_id`.

The wallet SHALL ingest a request only when:

`relay_account_event_id <= LatestTaskFeeLogID`

This guarantees that account logs up to the withdrawal anchor event are already synchronized before local withdrawal processing starts.

## Withdrawal Balance Ownership Rule

The wallet SHALL keep withdrawal balance ownership in withdrawal execution flow:

- On confirmed chain transfer, local account balance is debited in withdrawal processing.
- On failed or rejected withdrawal, log sync does not apply withdrawal balance adjustment.
- Log sync MUST ignore `Withdraw` and `WithdrawRefund` to avoid double-application with wallet-side withdrawal accounting.

Complete withdrawal lifecycle behavior is specified in [withdrawal_processing.md](./withdrawal_processing.md). Implementations MUST follow that specification.

