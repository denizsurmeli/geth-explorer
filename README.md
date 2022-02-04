# geth-explorer
Just some dirty-quick monday morning hack for playing with geth API

### Usage
```
go run explore.go [--network] [--opeartion|--txhash|--blockhash]
```
- `--network`: Specify the network. You can select `mainnet`,`rinkeby`,`ropsten`,`goerli`. By default it's `mainnet`.
- `--operation`: There are currently five options. Maybe more in the future ?
    - `listen_blocks`: Simple block view. Starts fetching from the latest found block, listens the network for current updates. Shows `blockNumber`,`blockSize`,`blockHash`, `numberOfTransactions` and a simple view of transactions(`from`,`to`,`value`,`wei`,`totalValue`).
    - `listen_headers`:Simple header view. Starts fetching from the latest found block, listens the network for current updates. Shows `parentHash`,`unclesHash`,`miner`,`stateRoot`,`transactionsRoot`,`receiptRoot`,`difficulty`,`GasLimit`,`Gas Used`,`Timestamp`.
    - `lens_transaction`: Simple tx view. Don't forget to pass `txhash`.
    - `lens_block`:Simple block view. Don't forget to pass `blocknumber`
    - `lens_txpool`: Only shows the pending `tx`'s in a flow. 
### NOTE
Only partially tested with `Infura` API. You only need to pass `projectId` in your `.env` file. 