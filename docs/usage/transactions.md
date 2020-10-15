# Generating and Signing Transactions

{% hint style="info" %}
Before following the instructions below, make sure your Ledger wallet is
unlocked and the Oasis App is open.
{% endhint %}

{% hint style="warning" %}
While the Oasis App is available in _Developer mode_, opening the App brings
up the "Pending Ledger review" screen.
You need to press both buttons at once to close that screen and transition to
the _ordinary_ "Oasis Ready" screen where the Oasis App is ready to be used.
{% endhint %}

As described in the [Use Your Tokens' Setup] document of the general
[Oasis Docs], you need to set the appropriate [Base and Signer CLI flags] for
each transaction you want to generate.

Make sure you set the following environment variables:

- `GENESIS_FILE`: Location of the genesis file.
- `LEDGER_SIGNER_PATH`: Location of the `ledger-signer` binary.
  See [Setup] for more details.

For convenience, you can set the `TX_FLAGS` environment variable like below:

```bash
TX_FLAGS=(--genesis.file "$GENESIS_FILE"
  --signer.dir entity
  --signer.backend plugin
  --signer.plugin.name ledger
  --signer.plugin.path "$LEDGER_SIGNER_PATH"
)
```

Make sure you replace `entity` with the name of the directory that contains the
`entity.json` file for you Ledger wallet's account.
See [Exporting Public Key to Entity] for more details.

{% hint style="info" %}
In case you will have more than one Ledger wallet connected, you will need to
specify which wallet to use by setting the `wallet_id` configuration key in
the `--signer.plugin.config` flag, i.e.

```
--signer.plugin.config "wallet_id:<LEDGER-WALLET-ID>"
```

where `<LEDGER-WALLET-ID>` is replaced with the ID of your Ledger wallet.
See [Identifying Wallets] for more details.
{% endhint %}

{% hint style="info" %}
If you want to use different account index for the same Ledger wallet, you
will need to specify it by setting the `index` configuration key in the
`--signer.plugin.config` flag, i.e.

```
--signer.plugin.config "index:<LEDGER-ACCOUNT-INDEX>"
```

where `<LEDGER-ACCOUNT-INDEX>` is replaced with the account index you want to
use.

If you need to specify multiple configuration keys in the
`--signer.plugin.config` flag, you can separate them with a comma (`,`), e.g.

```
--signer.plugin.config "wallet_id:1fc3be,index:5"
```

{% endhint %}

Then, you can generate and sign a transaction, e.g. a transfer transaction, by
running:

```bash
oasis-node stake account gen_transfer \
  "${TX_FLAGS[@]}" \
  --stake.amount <AMOUNT-TO-TRANSFER> \
  --stake.transfer.destination <DESTINATION-ACCOUNT-ADDRESS> \
  --transaction.file tx_transfer.json \
  --transaction.nonce 1 \
  --transaction.fee.gas 1000 \
  --transaction.fee.amount 2000
```

where `<AMOUNT-TO-TRANSFER>` and `<DESTINATION-ACCOUNT-ADDRESS>` are replaced
with the amount of tokens to transfer and the address of the transfer's
destination account, respectively.

{% hint style="info" %}
For a more detailed explanation of the transaction flags that were set, see
[Common Transaction Flags] section of the [Use Your Tokens' Setup] doc.
{% endhint %}

Next, verify the transaction's fields on your Ledger wallet's screen.

After you've confirmed the transaction's fields are correct, sign the
transaction on your Ledger wallet by double-pressing the _Sign transaction_
screen.

<!-- markdownlint-disable line-length -->
[Use Your Tokens' Setup]: https://docs.oasis.dev/general/use-your-tokens/setup
[Oasis Docs]: https://docs.oasis.dev/
[Base and Signer CLI flags]:
  https://docs.oasis.dev/general/use-your-tokens/setup#common-cli-flags
[Common Transaction Flags]:
  https://docs.oasis.dev/general/use-your-tokens/setup#common-transaction-flags
[Setup]: setup.md#remembering-path-to-ledger-signer-plugin
[Exporting Public Key to Entity]: entity.md
[Identifying Wallets]: wallets.md
<!-- markdownlint-enable line-length -->
