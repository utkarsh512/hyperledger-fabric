# Hyper-ledger Fabric 

Guide for CS61065 Assignment 4

## Installing

* Install prereqs listed [here](https://hyperledger-fabric.readthedocs.io/en/release-2.2/prereqs.html).

* Run the docker instance as 

    ```shell
    $ sudo dockerd
    ```

* Install hyperledger fabric binaries as described [here](https://hyperledger-fabric.readthedocs.io/en/release-2.2/install.html).


## Initializing test network

* Go to the `test-network` directory as

    ```shell
    $ cd fabric/fabric-samples/test-network
    ```

* Run the following command

    ```shell
    $ ./network.sh up
    ```

This will create:

* Two organizations each having one peer

* One orderer

* One fabric tools instance


## Creating channel

* Run the command

    ```shell
    $ ./network.sh createChannel -c <channel_name>
    ```

## Creating chaincode

* Run the command

    ```shell
    $ ./network.sh deployCC -ccn <chain_code_name> -ccp <chain_code_src_path> -ccl <language> -c <channel_name>
    ```

    For example, there are several chaincodes in `fabric-samples` directory. We can run following command:

    ```shell
    $ ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
    ```
## Invoking chaincode

First, we need to configure peer for query/invoke operations. 

###  Configuring peer

Run the following set of commands:

```shell
$ export PATH=${PWD}/../bin:$PATH
$ export FABRIC_CFG_PATH=$PWD/../config/
$ export CORE_PEER_TLS_ENABLED=true
$ export CORE_PEER_LOCALMSPID="Org1MSP"
$ export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
$ export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
$ export CORE_PEER_ADDRESS=localhost:7051
```

### Chaincode invokation

Run the following command

### Query transaction example

```shell
$ peer chaincode query -C mychannel -n basic -c '{"Args":["Inorder"]}'
```

#### Invoke transaction example

```shell
$ peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'
```

**Note**: Kindly change the arguments to the following flag:

* `-C`: channel name
* `-n`: chaincode name
* `-c`: function name (with the arguments) to invoke

## Writing smart contracts

* Create a directory as `fabric/fabric-samples/cs61065-chaincode`.

* Create sub-directory:

    * Part A: `student-register`
    * Part B: `bst`

* Build the file as:

    ```shell
    $ go build
    ```

* If some import error occurs, run the following command and build again

    ```shell
    $ go get
    ```

* Deploy the chaincode on test-network as

```shell
$ cd ../test-network
$ ./network.sh deployCC -ccn testchaincode -ccp ../test-local -ccl go
```

## Creating DApps

* Stop the network and start it again as

    ```shell
    $ ./network.sh down
    $ ./network.sh up -ca
    $ ./network.sh createChannel
    ```

* Deploy the chaincode as described in previous sections.

* Create a directory `cs61065-application` inside `test-network` and create an empty `main.js` file.

* Install `fabric-ca-client`, `fabric-network` and `prompt-sync` using `npm` under `testapp` directory as

    ```shell
    $ cd testapp
    $ npm install fabric-ca-client
    $ npm install fabric-network
    $ npm install prompt-sync
    ```