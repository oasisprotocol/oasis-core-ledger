# Exporting Public Key to Entity

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

Before you will be able to sign anything, you will need to export the public key
of your Ledger wallet's account and use it to generate an entity.

Make sure you set the following environment variables:

- `LEDGER_SIGNER_PATH`: Location of the `ledger-signer` binary.
  See [Setup] for more details.

To export the public key and generate an entity in the `entity`
subdirectory, run:

```bash
mkdir entity
oasis-node signer export \
  --signer.dir entity \
  --signer.backend plugin \
  --signer.plugin.name ledger \
  --signer.plugin.path "$LEDGER_SIGNER_PATH"
```

This will create an `entity.json` file in the `entity` directory that contains
the public key for a private & public key pair generated on your Ledger wallet.

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

[Setup]: setup.md#remembering-path-to-ledger-signer-plugin
[Identifying Wallets]: wallets.md
