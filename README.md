# Distributed-Banking-System

Simple distributed banking system using the [Raft](https://raft.github.io/) consensus algorithm.

This implement referenced [Creating Distributed KV Database by Implementing Raft Consensus Using Golang](https://yusufs.medium.com/creating-distributed-kv-database-by-implementing-raft-consensus-using-golang-d0884eef2e28) and used  [hashicorp/raft](https://github.com/hashicorp/raft) Go library for Raft.

## Source tree

Distributed-Banking-System  
└ go.mod  
└ go.sum  
└ README.md  
└ LICENSE  
└ node_1_data - data of node_1, including snapshots and dataRepo  
└ node_1_data - data of node_2, including snapshots and dataRepo  
└ node_1_data - data of node_3, including snapshots and dataRepo  
└ test.sh - testing scripts of Distributed-Banking-System  
└ reset.sh - reset Raft data script  
└ node_1_data - data of node_3, including snapshots and dataRepo  
└ node  
  └─ main.go - The main go file for running the raft node  
└ fsm - For making use of the replicated log [fsm](https://github.com/hashicorp/raft/blob/main/fsm.go)  
└ server  
  └─ raft_handler - The handler for raft cluster management  
  └─ banking_handler - The handler for basic features  
  └─ server.go - The HTTP server using Echo framework  

## How to run?

### Run servers

Set three or more server with different program in different terminal tab:

```shell
mkdir node_1_data node_2_data node_3_data

SERVER_PORT=2221 RAFT_NODE_ID=node1 RAFT_PORT=1111 RAFT_VOL_DIR=node_1_data go run dbs/node
SERVER_PORT=2222 RAFT_NODE_ID=node2 RAFT_PORT=1112 RAFT_VOL_DIR=node_2_data go run dbs/node
SERVER_PORT=2223 RAFT_NODE_ID=node3 RAFT_PORT=1113 RAFT_VOL_DIR=node_3_data go run dbs/node
```

### Creating clusters

After running each server, register nodes as followers to node_1 as a leader with POST to`/raft/join`.

```shell
curl --location --request POST 'localhost:2221/raft/join' \
--header 'Content-Type: application/json' \
--data-raw '{
	"node_id": "node_2", 
	"raft_address": "127.0.0.1:1112"
}'

>>>
{"data":{"applied_index":"3","commit_index":"3","fsm_pending":"0","last_contact":"0","last_log_index":"3","last_log_term":"2","last_snapshot_index":"0","last_snapshot_term":"0","latest_configuration":"[{Suffrage:Voter ID:node1 Address:127.0.0.1:1111} {Suffrage:Voter ID:node_2 Address:127.0.0.1:1112}]","latest_configuration_index":"0","num_peers":"1","protocol_version":"3","protocol_version_max":"3","protocol_version_min":"0","snapshot_version_max":"1","snapshot_version_min":"0","state":"Leader","term":"2"},"message":"node node_2 at 127.0.0.1:1112 joined successfully"}
```

```shell
curl --location --request POST 'localhost:2221/raft/join' \
--header 'Content-Type: application/json' \
--data-raw '{
	"node_id": "node_3", 
	"raft_address": "127.0.0.1:1113"
}'

>>>
{"data":{"applied_index":"4","commit_index":"4","fsm_pending":"0","last_contact":"0","last_log_index":"4","last_log_term":"2","last_snapshot_index":"0","last_snapshot_term":"0","latest_configuration":"[{Suffrage:Voter ID:node1 Address:127.0.0.1:1111} {Suffrage:Voter ID:node_2 Address:127.0.0.1:1112} {Suffrage:Voter ID:node_3 Address:127.0.0.1:1113}]","latest_configuration_index":"0","num_peers":"2","protocol_version":"3","protocol_version_max":"3","protocol_version_min":"0","snapshot_version_max":"1","snapshot_version_min":"0","state":"Leader","term":"2"},"message":"node node_3 at 127.0.0.1:1113 joined successfully"}
```

### Testing

Run test.sh for testing distributed banking system.
```shell
./test.sh
```

You can reset data of banking system.
```shell
./reset.sh
# After running the script, you should restart raft servers to reset data.
```

## Implement features

### Raft cluster management

#### POST /raft/join

You can add a new node to the cluster.

| Key          | Description                    | ValueType | Example          |
| ------------ | ------------------------------ | --------- | ---------------- |
| node_id      | The ID of the node             | string    | "node_1"         |
| raft_address | The address of the raft server | string    | "127.0.0.1:1111" |

Example:

```shell
curl --location --request POST 'localhost:2221/raft/join' \
--header 'Content-Type: application/json' \
--data-raw '{
	"node_id": "node_2", 
	"raft_address": "127.0.0.1:1112"
}'

>>>
{"data":{"applied_index":"3","commit_index":"3","fsm_pending":"0","last_contact":"0","last_log_index":"3","last_log_term":"2","last_snapshot_index":"0","last_snapshot_term":"0","latest_configuration":"[{Suffrage:Voter ID:node1 Address:127.0.0.1:1111} {Suffrage:Voter ID:node_2 Address:127.0.0.1:1112}]","latest_configuration_index":"0","num_peers":"1","protocol_version":"3","protocol_version_max":"3","protocol_version_min":"0","snapshot_version_max":"1","snapshot_version_min":"0","state":"Leader","term":"2"},"message":"node node_2 at 127.0.0.1:1112 joined successfully"}
```

#### DELETE /raft/remove/:node_id

You can delete one of the existing nodes in the cluster.

| Key     | Description                                | ValueType | Example  |
| ------- | ------------------------------------------ | --------- | -------- |
| node_id | The ID of the node that you want to remove | string    | "node_1" |

Example:

```sh
curl --location --request DELETE 'localhost:2221/raft/remove/node_2'

>>> {"message":"node node_2 removed successfully"}
```

#### GET /raft/leaderstats

You can get the node ID and IP address of a leader node.

Example:

```shell
curl --location --request GET 'localhost:2221/raft/leaderstats'

>>> {"data":{"Address":"127.0.0.1:1111","Id":"node1"},"message":"Here is the raft status"}
```

#### GET /raft/nodesstats

You can get the node IDs and IP addresses of all nodes in the cluster.

Example:

```shell
curl --location --request GET 'localhost:2221/raft/nodesstats'

>>> {"data":"[{Suffrage:Voter ID:node1 Address:127.0.0.1:1111} {Suffrage:Voter ID:node_3 Address:127.0.0.1:1113} {Suffrage:Voter ID:node_2 Address:127.0.0.1:1112}]","message":"Here is the nodes status in the cluster"}
```

### Basic Features

All states are persisted, even after the application is restarted.

#### POST /deposit

You can deposit a certain amount of tokens into one's account.

If the account does not exists, create a new one.

| Key     | Description                               | ValueType | Example    |
| ------- | ----------------------------------------- | --------- | ---------- |
| account | The account that you want to deposit      | string    | "account1" |
| amount  | Amount of tokens that you want to deposit | uint64    | 10         |

Example:

```shell
url --location --request POST 'localhost:2221/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
        "account": "account1", 
        "amount": 10
}'

>>> {"data":"Deposited 10 tokens to account1.","message":"success persisting data"}
```

#### POST /transfer

You can transfer a certain amount of tokens from one's account to another.

If the receiver's account does not exist, a new one is created.

If the sender's account does not exist, the system returns an error.

| Key      | Description                                | ValueType | Example    |
| -------- | ------------------------------------------ | --------- | ---------- |
| sender   | The sender's account                       | string    | "account1" |
| receiver | The receiver's account                     | string    | "account2" |
| amount   | Amount of tokens that you want to transfer | uint64    | 10         |

Example:

```shell
curl --location --request POST 'localhost:2221/transfer' --header 'Content-Type: application/json' --data-raw '{
        "sender": "account1", 
        "receiver": "account2",
        "amount": 10
}'

>>> {"data":"Transfered 10 tokens from account1 to account2.","message":"success persisting data"}
```

#### GET /get/:account

You can query a balance of an account.

If the account does not exist, the system returns an error.

| Key     | Description                                  | ValueType | Example    |
| ------- | -------------------------------------------- | --------- | ---------- |
| account | The account that you want to query a balance | string    | "account1" |

Example:

```shell
curl --location --request GET 'localhost:2221/get/account2'

>>> {"data":"The account 'account2' has 50 tokens.","message":"success persisting data"}
```

### Replication

All stats are replicated between a Raft cluster that consists of nodes.

## Improvement Plan

1. In this implementation, any clients can deposit and transfer tokens and get a balance, even not their accounts. The authentication scheme should be considered for the real-world banking systems.
2. An error message passing in the transfer feature was not implemented. This may be because I am not familiar with Go. I need to improve the Go skills.
3. Resonse message is not formatted.
