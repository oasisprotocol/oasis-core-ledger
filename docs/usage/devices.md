# Identifying Ledger devices

{% hint style="info" %}
To identify a Ledger device for use with the Oasis Core and Oasis Core Ledger
CLI tools, unlock the device and make sure the Oasis App is open.
{% endhint %}

To list all Ledger devices connected to a system, run:

```bash
oasis-core-ledger list_devices
```

If your Ledger device is properly connected, you should see an output similar to
the one below:

```text
- Wallet ID: 431fc6
  App version: 1.7.2
```

From now on, use the `Wallet ID` to identify the Ledger device you want to use.

For convenience, set the `LEDGER_WALLET_ID` environment value to its value,
e.g.:

```bash
LEDGER_WALLET_ID=431fc6
```
