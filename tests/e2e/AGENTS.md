## Build the Docker image for E2E testing

- Goal: Build the local image `crynux-relay-wallet:e2e` from the current repository source code.

- Workflow:
  1. Ensure the current working directory is the repository root.
  2. Build the image from `Dockerfile` with tag `crynux-relay-wallet:e2e`.
  3. Verify `crynux-relay-wallet:e2e` exists locally.

## Prepare the mount folder for the Docker image

- Goal: Prepare wallet host-side files in a Docker mount workspace for local e2e runs.

- Workspace rule:
  - Choose one mount workspace root directory and use it consistently in Docker volume mappings.
  - Reuse existing files in that workspace when possible instead of recreating everything.

- Required wallet structure:

```text
<mount-root>/
  config/
    config.yml
    blockchain_privkey.txt
    relay_api_privkey.txt
```

- Workflow:
  1. Create the `config` folder under `<mount-root>` if it does not already exist.
  2. Copy `tests/e2e/config.e2e.yml` to `<mount-root>/config/config.yml`.
  3. Mount `<mount-root>/config` to `/app/config` in the wallet container.

## Prepare the database

- Goal: Prepare a database instance for the wallet service during e2e runs.

- Database requirements:
  - A database instance must be available for the wallet container.
  - Configure the database connection in `<mount-root>/config/config.yml`.

## Prepare the system wallet account

- Goal: Prepare the blockchain private key and account information required by the wallet container.

- Required wallet file:

```text
<mount-root>/config/blockchain_privkey.txt
```

- System wallet account requirements:
  - `blockchains.<network>.account.address` in `<mount-root>/config/config.yml` must be set to the system wallet address used for e2e for each configured blockchain network.
  - The private key for that address must be stored in `<mount-root>/config/blockchain_privkey.txt` as a single-line hex value.
  - A `0x` prefix is allowed in `<mount-root>/config/blockchain_privkey.txt`.
  - The key material in `<mount-root>/config/blockchain_privkey.txt` must not contain trailing whitespace.
  - `blockchains.<network>.account.private_key_file` in `<mount-root>/config/config.yml` must point to `/app/config/blockchain_privkey.txt`.

## Prepare the Relay API key

- Goal: Prepare the private key used to authenticate Relay API requests.

- Required key file:

```text
<mount-root>/config/relay_api_privkey.txt
```

- Relay API key requirements:
  - The private key must be stored in `<mount-root>/config/relay_api_privkey.txt` as a single-line hex value.
  - A `0x` prefix is allowed in `<mount-root>/config/relay_api_privkey.txt`.
  - The key material in `<mount-root>/config/relay_api_privkey.txt` must not contain trailing whitespace.
  - `relay.api.private_key_file` in `<mount-root>/config/config.yml` must point to `/app/config/relay_api_privkey.txt`.

## Top up the system wallet with native tokens

- Goal: Ensure the system wallet has enough native token balance on every configured blockchain network before starting e2e. The balance covers both gas fees and withdrawal payouts.

- Minimum balance: 1000 tokens per blockchain network.

- Workflow:
  1. For each blockchain network configured in `<mount-root>/config/config.yml`, read the system wallet address from `blockchains.<network>.account.address`.
  2. Check the current native token balance for that address on the target blockchain.
  3. If the balance is below 1000 tokens, transfer enough native tokens to bring it to at least 1000.
  4. Confirm every configured network has at least 1000 tokens before running e2e.
