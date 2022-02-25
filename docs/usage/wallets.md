# Identifying Wallets

:::info

Before following the instructions below, make sure your Ledger wallet is
unlocked and the Oasis App is open.

:::

:::caution

While the Oasis App is available in _Developer mode_, opening the App brings
up the "Pending Ledger review" screen.
You need to press both buttons at once to close that screen and transition to
the _ordinary_ "Oasis Ready" screen where the Oasis App is ready to be used.

:::

To list all Ledger wallets connected to a system, run:

```bash
oasis-core-ledger list_devices
```

If your Ledger wallet is properly connected, you should see an output similar to
the one below:

```text
- Wallet ID: 431fc6
  App version: 1.7.2
```

You can pass this ID when you need to specify which Ledger wallet you want to
connect to via `--wallet_id` CLI flag or `wallet_id` configuration key.
