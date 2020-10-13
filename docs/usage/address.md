# Obtaining Account Address

Each account on your Ledger wallet has a corresponding
[staking account address].

To obtain a staking account address that corresponds to the account with index
0 on your Ledger wallet, use:

```bash
oasis-core-ledger show_address --wallet_id <LEDGER-WALLET-ID>
```

where `<LEDGER-WALLET-ID>` is replaced with the ID of your Ledger wallet.

See [Identifying Ledger Devices] for more details.

This will display your wallet's address and show it on your Ledger's screen for
confirmation.

To just display your wallet's address without showing it on your Ledger's
screen, pass the `--no-show` flag in the command above.

{% hint style="info" %}
You can obtain as many staking account addresses as needed for the same Ledger
wallet by passing the `--index` flag and specifying a different account index in
the command above.
{% endhint %}

[staking account address]:
  https://docs.oasis.dev/general/use-your-tokens/account/address
[Identifying Ledger Devices]: devices.md