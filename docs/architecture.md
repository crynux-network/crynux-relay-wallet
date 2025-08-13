## System Architecture

### Relay API

The Relay Wallet interacts with Relay by calling four API endpoints:

- GetTaskFeeLogs
- GetWithdrawalRequests
- FullfillWithdrawalRequest
- RejectWithdrawalRequest

The first two fetch data from Relay, while the latter two report the processing status of withdrawals back to Relay.

The Relay API client code resides in `/relay_api` as wrapper functions for these calls. This directory also defines the data structures for task fee logs and withdrawal requests.

### Task Fee Log Sync

Task fee log synchronization starts from `tasks.StartSyncTaskFeeLogs`, which is launched as a goroutine in `/main.go`. It runs continuously in the background, retrieving the latest task fee logs from the Relay API and updating the locally stored balances of node relay accounts based on those logs.

`GetTaskFeeLogs` supports incremental, gap-free syncing via a pivot ID. On startup, the Relay Wallet reads the last synced task fee log ID from the local database (stored in `models/system.go`). If the record is empty, no prior sync has occurred, so syncing begins at ID=0. Otherwise, the stored ID is passed to the API, which returns records starting from the next ID.

After retrieving records, the wallet updates, within a single database transaction, both the node relay account balances and the latest task fee log ID to ensure they are always consistent.

### Withdrawal Request Processing

TODO

### Other Components

#### DB Models

- Uses gorm v2 for database operations.
- All models are under `/models`.
- Database initialization and configuration are in `/config/db.go`.

#### Logs

- A global logger is configured. Initialization code is in `/config/log.go`.

#### Tasks

- The application starts multiple long-running, concurrently executing background tasks (for example, syncing task fee logs and processing withdrawals).
- Task code resides in `/tasks` and is started from `/main.go` as goroutines.
- Graceful shutdown is supported via OS signal handling (e.g., Ctrl+C). From `/main.go`, signals are captured and used to cancel contexts passed into tasks. Tasks listen to `ctx.Done()` and exit cleanly after finishing in-flight work (committing any pending operations before returning).
- Task exceptions are isolated: if a task encounters an unrecoverable error or panic, only that task exits and an alert is sent via the standard alerting path. Other tasks continue running and remain unaffected. Implemented in `/main.go`.

#### DB Migration

- Model changes are applied via migrations to keep different deployment environments in sync.
- Any model change must first add a migration under `/migrate`. See `/migrate/migrate.go` for the mechanism.
- On startup, `/main.go` automatically runs database migrations.

#### Alert

- The project integrates with external services (such as AWS CloudWatch) to provide alerting. Related code is in `/alert`.
- To ensure the alerting path itself remains healthy and to avoid missing alerts due to bugs or unexpected shutdowns of the Relay Wallet, the system sends proactive heartbeats through the exact same alerting pathway (same service and same code). The external service is configured to trigger an alert if heartbeats are missing.
