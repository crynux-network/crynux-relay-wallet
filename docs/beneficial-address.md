## Beneficial Address

A beneficial address is an optional, dedicated address that receives all tokens associated with a node (stake refunds, emissions, and withdrawals). Its private key never needs to be online. This reduces the impact of a compromise of the node's operational key.

### The problem

- The node must keep its operational wallet private key online to sign requests and transactions.
- If that online key is leaked or the host is compromised, an attacker can attempt to redirect funds or capture future payouts and stake refunds.
- Protecting the operational key alone is not sufficient; we need a design that keeps tokens safe even if the node key is exposed.

### The solution (at a glance)

- Introduce a one‑time, immutable "beneficial address" binding for each node address.
- Keep the beneficial address key offline; it is never required for node operations.
- Route all tokens (stake refunds, token emissions, and Relay Wallet withdrawals) to the beneficial address instead of the node address.
- Bind withdrawal requests to the beneficial address and independently verify that binding against the blockchain.

### Design and rationale

- One‑time, immutable binding
  - The beneficial address is set via a dedicated on‑chain transaction and cannot be changed thereafter, even if the node later quits and rejoins.

  Why: Prevents an attacker with the node's online key from swapping the payout destination. Immutability eliminates the "change‑the‑address" attack.

- Offline custody of the beneficial key
  - The private key of the beneficial address is not needed for any node or Relay interaction and should be stored offline (cold storage).

  Why: Minimizes exposure. Even if the node machine is fully compromised, the attacker cannot spend funds held at the beneficial address.

- Unified routing of tokens to the beneficial address
  - Stake refunds after a node quits are paid to the beneficial address.
  - Ongoing token emissions accrue to the beneficial address.
  - Relay Wallet withdrawals are sent to the beneficial address rather than the node address.

  Why: Consolidates all token flows behind the offline‑key address, removing opportunities for an attacker with the node key to redirect or capture payouts.

- Withdrawal request binding and verification
  - Withdrawal requests must include the node's beneficial address and be signed by the node's operational key.
  - The Relay Wallet independently reads the on‑chain beneficial binding via its own blockchain node and cross‑checks it with the request.
  
  Why: Prevents tampering in transit or UI spoofing. Cross‑checking on chain also protects against a compromised blockchain node or frontend that might misstate the bound address.

### Operational notes

- Until the beneficial address is set, funds will be routed to the node address; this increases risk because the online key can spend those funds.
- Because the binding is immutable, choose the beneficial address carefully and ensure its key is backed up securely offline.
