## Coding Requirements

### Full Picture

Before making any changes, consult `./README.md` and `./docs/architecture.md` for a high-level project overview. All modifications must be consistent with the existing architecture and design.

### Data Integrity

Because the Relay Wallet handles funds, the correctness of financial data must be strictly guaranteed under all circumstances.

Processing of logs and withdrawals may be delayed, but any data that has been processed must remain correct and consistent. In particular, during unexpected exceptions and shutdown, ensure that in-flight operations do not stop at a point that leaves data in an inconsistent state.

### Proper Error Handling

All function errors must be propagated up the call stack until handled. Any unhandled error reaching the `main` function must be logged and trigger an alert to operators.

### Proper Logging

Add sufficient logging at appropriate points in the code with the correct log levels, so both operators and developers can identify and diagnose issues easily.

### Clean Comment

Do not explain code changes in comments, such as "added xxx" or "removed xxx because xxx". Only describe the functionality of the final code. Keep comments concise and only add them for complex or non-obvious logic.

Do not use comments to delete code; directly remove the code without adding explanations about what was deleted.

## E2E Test

For E2E test execution instructions, use `tests/e2e/AGENTS.md` as the source of truth.
