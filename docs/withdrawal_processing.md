# Withdrawal Processing Specification

## Scope

This document defines the end-to-end withdrawal handling flow in Relay Wallet, including request synchronization, validation, local persistence, blockchain execution, callback reporting, and timeout behavior.

## Entry Tasks

The wallet SHALL run two long-lived tasks for withdrawals:

- `StartSyncWithdrawalRequests`: pulls and stores requests from Relay.
- `StartProcessWithdrawalRequests`: executes locally stored requests and reports result.

## Request Synchronization Flow

`syncWithdrawalRequests` SHALL process withdrawals in batches using `LatestWithdrawalRequestID` as the pull cursor.

For each fetched batch:

1. Read current `LatestTaskFeeLogID` from task-fee checkpoint.
2. Apply gate rule:
   - A request is ingestible only if `relay_account_event_id <= LatestTaskFeeLogID`.
   - At first request that violates this rule, stop batch at that point.
3. If no ingestible requests remain, wait one sync interval and retry.
4. Validate ingestible requests with `checkWithdrawalRequests`.
5. Insert local `withdraw_records` with `OnConflict DoNothing`.
6. Update withdrawal checkpoint to the last ingested request in the same transaction.

## Validation Rules (`checkWithdrawalRequests`)

The wallet MUST enforce all rules below before storing a request:

- Amount MUST be parseable as integer and MUST be greater than or equal to configured minimum withdrawal amount.
- Request status MUST be `pending`.
- Aggregated per-network withdrawal amount in the batch MUST NOT exceed system wallet on-chain balance for that network.
- Every request address in the batch MUST already exist in local `relay_accounts`.
- Aggregated per-address withdrawal amount in the batch MUST NOT exceed local account balance.
- Benefit address fetched from chain (`GetBenefitAddress`) MUST equal request `benefit_address`.

Validation failure SHALL fail the sync attempt.

## Local Record Model and Status

The wallet stores each accepted request as `withdraw_records` with local status lifecycle:

- `pending` -> `success` or `failed` -> `finished`

`success` and `failed` represent local execution outcome before relay callback completion.
`finished` represents callback completion (`fulfill` or `reject`) and local finalization.

## Execution Flow (`processWithdrawalRecord`)

For each unfinished local record:

1. If no blockchain transaction is attached, queue send transaction (`QueueSendETH`) and store `blockchain_transaction_id`.
2. Poll transaction status until terminal (`confirmed` or `failed`) or context cancellation.
3. If confirmed:
   - Load local relay account by record address.
   - Verify local balance is sufficient for record amount.
   - Decrease local balance by record amount.
   - Update record status to `success` in the same transaction.
4. If failed and retries are exhausted, update record status to `failed`.
5. After leaving pending loop:
   - If status is `success`, call Relay `FulfillWithdrawalRequest` with tx hash.
   - Otherwise call Relay `RejectWithdrawalRequest`.
6. Set local record status to `finished`.

## Timeout Handling

Each record processing attempt SHALL run with a per-record deadline:

- deadline = `record.CreatedAt + ProcessWithdrawalRequests.Timeout`

If deadline is exceeded:

- Call Relay `RejectWithdrawalRequest`.
- Set local status to `finished`.
- Return timeout error for alerting path.

## Balance Ownership Rule

Local account balance adjustment for withdrawals SHALL remain owned by withdrawal execution flow:

- Balance is decreased only after confirmed chain transfer.
- Withdrawal-related Relay account logs (`Withdraw`, `WithdrawRefund`) are not used to adjust local balance in log sync.

