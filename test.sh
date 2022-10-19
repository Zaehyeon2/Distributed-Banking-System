echo '###### Run test.sh after running Raft servers'
echo '########## Remove existing nodes'
curl --location --request DELETE 'localhost:2221/raft/remove/node_2'
curl --location --request DELETE 'localhost:2221/raft/remove/node_3'

echo ''
echo '###### (Get nodestats) There is only one node in the cluster "localhost:2221"'
curl --location --request GET 'localhost:2221/raft/nodesstats'

echo ''
echo '###### (Join node) Registering localhost:2222 as a follower to localhost:2221 as a leader'
curl --location --request POST 'localhost:2221/raft/join' \
--header 'Content-Type: application/json' \
--data-raw '{
	"node_id": "node_2", 
	"raft_address": "127.0.0.1:1112"
}'
echo ''
echo '###### (Get nodestats) You can find there are two nodes in the cluster'
curl --location --request GET 'localhost:2221/raft/nodesstats'

echo ''
echo '###### (Join node) Registering localhost:2223 as a follower to localhost:2221 as a leader'
curl --location --request POST 'localhost:2221/raft/join' \
--header 'Content-Type: application/json' \
--data-raw '{
	"node_id": "node_3", 
	"raft_address": "127.0.0.1:1113"
}'

echo ''
echo '###### (Get nodestats) You can find there are three nodes in the cluster'
curl --location --request GET 'localhost:2221/raft/nodesstats'

echo ''
echo '###### (Get leaderstats) You can find NodeId and IP address of the leader node'
curl --location --request GET 'localhost:2221/raft/leaderstats'

echo ''
echo '###### (Delete node) Deleting localhost:2222'
curl --location --request DELETE 'localhost:2221/raft/remove/node_2'

echo ''
echo '###### (Get nodestats) You can find there are two nodes in the cluster'
curl --location --request GET 'localhost:2221/raft/nodesstats'

echo ''
echo '###### (Join node) Registering localhost:2222 as a follower to localhost:2221 as a leader'
curl --location --request POST 'localhost:2221/raft/join' \
--header 'Content-Type: application/json' \
--data-raw '{
	"node_id": "node_2", 
	"raft_address": "127.0.0.1:1112"
}'

echo ''
echo '###### (Get nodestats) You can find there are three nodes in the cluster'
curl --location --request GET 'localhost:2221/raft/nodesstats'

echo ''
echo '###### (Deposit) Depositing 100 tokens into account1'
curl --location --request POST 'localhost:2221/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
        "account": "account1", 
        "amount": 100
}'

echo ''
echo '###### (Deposit) Depositing 0 tokens into account2 (creating account2)'
curl --location --request POST 'localhost:2221/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
        "account": "account2", 
        "amount": 0
}'

echo ''
echo '###### (Get) You can find account1 has 100 tokens'
curl --location --request GET 'localhost:2221/get/account1'

echo ''
echo '###### (Get) You can find a balance of each account'
echo '########## (Get) A balance of account1'
curl --location --request GET 'localhost:2221/get/account1'
echo ''
echo '########## (Get) A balance of account2'
curl --location --request GET 'localhost:2221/get/account2'

echo ''
echo '###### (Transfer) Transferring 10 tokens from account1 into account2'
curl --location --request POST 'localhost:2221/transfer' --header 'Content-Type: application/json' --data-raw '{
        "sender": "account1", 
        "receiver": "account2",
        "amount": 10
}'

echo ''
echo '###### (Get) You can find a balance of each account'
echo '########## (Get) A balance of account1'
curl --location --request GET 'localhost:2221/get/account1'
echo ''
echo '########## (Get) A balance of account2'
curl --location --request GET 'localhost:2221/get/account2'

echo ''
echo '###### (Transfer) If sender`s account does not have enough tokens, return error.'
curl --location --request POST 'localhost:2221/transfer' --header 'Content-Type: application/json' --data-raw '{
        "sender": "account1", 
        "receiver": "account2",
        "amount": 100000000000
}'curl --location --request POST 'localhost:2221/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
	"account": "account1", 
	"amount": 10
}'

echo ''
echo '###### (Transfer) If sender`s account does not exist (account123), return error.'
curl --location --request POST 'localhost:2221/transfer' --header 'Content-Type: application/json' --data-raw '{
        "sender": "account123", 
        "receiver": "account2",
        "amount": 100000000000
}'

echo ''
echo '###### (Get) If account does not exist (account123), return error.'
curl --location --request GET 'localhost:2221/get/account123'
