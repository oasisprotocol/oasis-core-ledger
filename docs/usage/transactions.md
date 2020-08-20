# Generating and Signing Transactions

As described in the [Stake Management] part of [Oasis Docs], you need to set
the appropriate **base flags** and **signer flags** for each transaction you
want to generate.

Make sure you set the following environment variables:

- `LEDGER_SIGNER_PATH`: Location of the `ledger-signer` binary.
  See [Setup] for more details.
- `LEDGER_WALLET_ID`: ID of the Ledger wallet to use.
  See [Identifying Ledger Devices] for more details.
- `LEDGER_INDEX`: Index (0-based) of the account on the Ledger device to use.

For convenience, you can set the `TX_FLAGS` environment variable like below:

```bash
TX_FLAGS=(--genesis.file /localhostdir/genesis.json
  --signer.dir entity-$LEDGER_INDEX
  --signer.backend plugin
  --signer.plugin.name ledger
  --signer.plugin.path $LEDGER_SIGNER_PATH
  --signer.plugin.config "wallet_id:$LEDGER_WALLET_ID,index:$LEDGER_INDEX"
)
```

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
[Common Transaction Flags] section of the [Stake Management] docs.
{% endhint %}

Next, verify the transaction's fields on your Ledger device's screen.

After you've confirmed the transaction's fields are correct, sign the
transaction on your Ledger device by double-pressing the _Sign transaction_
screen.

<!-- markdownlint-disable line-length -->
[Stake Management]:
  https://docs.oasis.dev/general/operator-docs/stake-management
[Oasis Docs]: https://docs.oasis.dev/
[Common Transaction Flags]:
  https://docs.oasis.dev/general/operator-docs/stake-management#common-transaction-flags
[Setup]: setup.md#remembering-path-to-ledger-signer-plugin
[Identifying Ledger Devices]: devices.md
<!-- markdownlint-enable line-length -->
