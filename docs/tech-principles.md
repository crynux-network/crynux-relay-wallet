
## Exception Handling and Graceful Shutdown

Because the Relay Wallet handles funds, the correctness of financial data must be strictly guaranteed under all circumstances.

Processing of logs and withdrawals may be delayed, but any data that has been processed must remain correct and consistent. In particular, during unexpected exceptions and during shutdown, ensure that in-flight operations do not stop at a point that leaves data in an inconsistent state.

## Logging

Add sufficient logging at appropriate points in the code with the correct log levels, so both operators and developers can identify and diagnose issues easily.
