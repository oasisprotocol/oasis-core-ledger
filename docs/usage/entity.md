# Exporting Public Key to Entity

Before you can sign anything, you need to export the public key from a Ledger
device and use it to generate an entity.

Make sure you set the following environment variables:

- `LEDGER_SIGNER_PATH`: Location of the `ledger-signer` binary.
  See [Setup] for more details.
- `LEDGER_WALLET_ID`: ID of the Ledger wallet to use.
  See [Identifying Ledger Devices] for more details.
- `LEDGER_INDEX`: Index (0-based) of the account on the Ledger device to use.

To export the public key and generate an entity in the `entity-$LEDGER_INDEX`
subdirectory, run:

```bash
mkdir entity-$LEDGER_INDEX
oasis-node signer export \
  --signer.dir entity-$LEDGER_INDEX \
  --signer.backend plugin \
  --signer.plugin.name ledger \
  --signer.plugin.path $LEDGER_SIGNER_PATH \
  --signer.plugin.config "wallet_id:$LEDGER_WALLET_ID,index:$LEDGER_INDEX"
```

This will create an `entity.json` file in the `entity` directory that contains
the public key for a private & public key pair generated on the Ledger device.

Account index specifies the `address_index` part of the [BIP32] path conforming
to the [BIP44] specification.

{% hint style="info" %}
You can obtain as many entities as needed for the same Ledger wallet by
specifying a different account index in `LEDGER_INDEX` environment variable
and re-running the steps above.
{% endhint %}

[Setup]: setup.md#remembering-path-to-ledger-signer-plugin
[Identifying Ledger Devices]: devices.md
[BIP32]: https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
[BIP44]: https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
