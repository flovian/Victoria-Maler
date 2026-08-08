# Bitcoin Integration

Bitcoin Core is the trust and verification layer of Victoria Maler. Cleanup
reports and evidence files are hashed with SHA-256, and that hash is anchored to
Bitcoin in an OP_RETURN output. Once a hash is on-chain it cannot be altered, so
tampering with a record becomes immediately detectable.

## Anchoring flow

1. A cleanup report or evidence file is submitted.
2. A canonical string (or the file bytes) is hashed with SHA-256.
3. `bitcoin.AnchoringService.AnchorHash` frames the hash with the `ECOVIC|`
   prefix into a standard OP_RETURN script.
4. When a bitcoind node is available the transaction is created, funded,
   signed and broadcast. The resulting `txid` is stored on the record.
5. Without a node, a clearly labelled `sim-<hash-prefix>` txid is stored and
   the record is marked `simulated`.

## On-chain verification

`bitcoin.VerifyHash(hash, txid)`:

* For `sim-` txids it reports a `simulated` verification against the stored hash.
* For real txids it fetches the transaction with `getrawtransaction` and checks
  that one of the `vout` scriptPubKeys carries the `ECOVIC|<hash>` payload.

## Configuration (`.env`)

| Variable | Default | Description |
|---|---|---|
| `BITCOIN_ENABLED` | `false` | When `true`, attempts real on-chain anchoring. |
| `BITCOIN_RPC_URL` | `http://127.0.0.1:18443` | bitcoind RPC endpoint. |
| `BITCOIN_RPC_USER` | `bitcoin` | RPC username. |
| `BITCOIN_RPC_PASS` | `bitcoin` | RPC password. |
| `BITCOIN_NETWORK` | `regtest` | Network label shown on the status page. |

## Running with a real node

1. Start bitcoind with a wallet: `bitcoind -regtest -server -rpcuser=bitcoin -rpcpassword=bitcoin -fallbackfee=0.0002`.
2. Mine a block and fund the wallet: `bitcoin-cli -regtest generatetoaddress 101 $(bitcoin-cli -regtest getnewaddress)`.
3. Set `BITCOIN_ENABLED=true` in `.env` and restart the server.

New submissions are then anchored on-chain and verifiable as `anchored`.
Existing `simulated` records remain clearly labelled as non-on-chain.

## Tamper detection

Recomputing the hash of a record and comparing it to the stored (anchored) hash
reveals any modification. `GET /api/cleanups/verify?id=N` performs this check
and reports a `verified: false` result when a record has been changed.
