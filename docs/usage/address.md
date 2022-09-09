# Obtaining Account Address

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

Each account on your Ledger wallet has a corresponding
[staking account address].

To obtain a staking account address that corresponds to the account with index
0 on your Ledger wallet, use:

```bash
oasis-core-ledger show_address
```

This will display your wallet's address and show it on your Ledger's screen for
confirmation.

To skip showing your wallet's address on your Ledger's screen, pass the
`--skip-device` flag in the command above.

:::info

In case you will have more than one Ledger wallet connected, you'll need to
specify which wallet to use by passing the `--wallet_id <LEDGER-WALLET-ID>` flag
to the command above, replacing `<LEDGER-WALLET-ID>` with the ID of your Ledger
wallet.
See [Identifying Wallets] for more details.

:::

:::info

You can obtain as many staking account addresses as needed for the same Ledger
wallet by passing the `--index` flag and specifying a different account index in
the command above.

Account index specifies the `address_index` part of the [BIP32] path conforming
to the [BIP44] specification.

:::

<!-- markdownlint-disable line-length -->
[staking account address]:
  https://github.com/oasisprotocol/docs/blob/main/docs/general/manage-tokens/terminology.md#address
[Identifying Wallets]: wallets.md
[BIP32]: https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
[BIP44]: https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
<!-- markdownlint-enable line-length -->
