# geth-explorer
Just some dirty-quick monday morning hack for playing with geth API

### Usage
```
go run explore.go [--network] [--subscribe]
```
- `--network`: Specify the network. You can select `mainnet`,`rinkeby`,`ropsten`,`goerli`
- `--subscribe`: There are currently two options. Maybe more in the future ?
    - `blocks`: Simple block view. Starts fetching from the latest found block, listens the network for current updates. Shows `blockNumber`,`blockSize`,`blockHash`, `numberOfTransactions` and a simple view of transactions(`from`,`to`,`value`,`wei`,`totalValue`).
    - `headers`:Simple header view. Starts fetching from the latest found block, listens the network for current updates. Shows `parentHash`,`unclesHash`,`miner`,`stateRoot`,`transactionsRoot`,`receiptRoot`,`difficulty`,`GasLimit`,
    `Gas Used`,`Timestamp`.

### NOTE
Only partially tested with `Infura` nodes. Only pass `projectId` in your `.env` file. 