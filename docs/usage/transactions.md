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

Then, you can generate and sign a transaction by running:

```bash
oasis-node stake account gen_<TX-TYPE> \
  "${TX_FLAGS[@]}" \
  --transaction.file tx.json \
  --transaction.nonce <NONCE> \
  --transaction.fee.gas <GAS-LIMIT> \
  --transaction.fee.amount <FEE>
```

where:

- `<TX-TYPE>`: type of transaction, e.g. `transfer`, `escrow`, `reclaim_escrow`,
  ...
- `<NONCE>`: your account's current nonce,
- `<GAS-LIMIT>`: maximum amount of gas this transaction can spend,
- `<FEE>`: amount of tokens you will pay as a fee for this transaction.

Besides these common transaction flags, you will need to specify additional
transaction flags specific to the chosen transaction type. Run
`oasis-node stake account gen_<TX-TYPE> --help` for more details.

{% hint style="info" %}
For a more detailed explanation of the common transaction flags, see
[Common Transaction Flags] section of the [Use Your Tokens' Setup] doc.
{% endhint %}

For example, to generate and sign a transfer transaction of 100 tokens to an
account with address `oasis1qpcgnf84hnvvfvzup542rhc8kjyvqf4aqqlj5kqh`, run:

```bash
oasis-node stake account gen_transfer \
  "${TX_FLAGS[@]}" \
  --stake.amount 100000000000 \
  --stake.transfer.destination oasis1qpcgnf84hnvvfvzup542rhc8kjyvqf4aqqlj5kqh \
  --transaction.file tx.json \
  --transaction.nonce 1 \
  --transaction.fee.gas 2000 \
  --transaction.fee.amount 2000
```

{% hint style="info" %}
The amounts passed via the `--stake.amount` and `--transaction.fee.amount` flags
are specified in nROSE units, i.e. 1 ROSE equals 1,000,000,000 nROSE.
{% endhint %}

This will output a preview of the generated transaction:

```
You are about to sign the following transaction:
  Nonce:  1
  Fee:
    Amount: ROSE 0.000002
    Gas limit: 2000
    (gas price: ROSE 0.000000001 per gas unit)
  Method: staking.Transfer
  Body:
    To:     oasis1qpcgnf84hnvvfvzup542rhc8kjyvqf4aqqlj5kqh
    Amount: ROSE 100.0
Other info:
  Genesis document's hash: a245619497e580dd3bc1aa3256c07f68b8dcc13f92da115eadc3b231b083d3c4
```

and ask you to verify the transaction's fields on your Ledger wallet's screen.

After you've confirmed the transaction's fields are correct, sign the
transaction on your Ledger wallet by double-pressing the _Sign transaction_
screen.

{% hint style="info" %}
The next step after signing a transaction is to submit it to the network via
an online Oasis node by running:

```bash
oasis-node consensus submit_tx \
  -a $ADDR \
  --transaction.file tx_transfer.json
```

For more details, see the [Transfer Tokens] document of the general
[Oasis Docs].
{% endhint %}

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
[Transfer Tokens]:
  https://docs.oasis.dev/general/use-your-tokens/transfer-tokens
<!-- markdownlint-enable line-length -->
