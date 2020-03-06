/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package bitcoin

import (
	"testing"

	"github.com/qredo/assets/libs/qredochain"
	"github.com/stretchr/testify/assert"
)

//Tricky bitcoin.conf setup, this works:

/*
testnet=1
maxconnections=12
muxuploadtarget=20
daemon=1
server=1

[test]
rpcclienttimeout=5000
rpcbind=172.31.11.118
rpcallowip=::/0
rpcuser=qredo
rpcpassword=KcsUHi4Hn89ELZJo66vygGtGA
*/

func Test_rpc(t *testing.T) {

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	client, err := BTCTestNetConnector()
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, client, "Client should not be nil")

}

func Test_AddAddressToWallet(t *testing.T) {
	/*
		Testing details - Use uncompressed
		Address = mpfcmiBS1LocFoxubnSdhPZBSttZNHXVXR (uncompressed)
		PrivKey = cP7H9KX8MY63tKYftZHo5fTiWTPgJaLzAfYiEWigVvCyb243jS6Y
		Pubkey = 04AD76B96E17964C15C4FB43AEA8E59110023286EBE4058EE73128B18EB88AB90C81D9C423F3687B6A4CD3D18FC655D5B8AEC9EA87AB097D1C4B3A06FAC4C8AD4B

		//Public Key is Uncompressed 130chars

		>>b importpubkey 04AD76B96E17964C15C4FB43AEA8E59110023286EBE4058EE73128B18EB88AB90C81D9C423F3687B6A4CD3D18FC655D5B8AEC9EA87AB097D1C4B3A06FAC4C8AD4B "Qredo1" false
		>>b listreceivedbyaddress 0 true true mpfcmiBS1LocFoxubnSdhPZBSttZNHXVXR

		//At this BlockHeight we have no Transactions
		InitiHash	000000000014dc9751e473e5559fc870e2b8a99965c015c265a8fe4a08f68c3f


		TX1			20a966e9aeec1c212a9d87975cb6c1b80aae468c3d803e5b7241b3a95d755844
		VAL			0.003
		Block		1665132
		BlockHash	000000000000006de69aa36edcaf2ed0b80c0bfd53204ccd21dae80962a8cae7

		TX2			f77fee60787e9b5a742dbcca50cacd0a3d2edaa488930e5facec2eb1e6960253
		VAL			0.002
		Block		1665133
		BlockHash	00000000f95f5703387b1a799d4075cad325a88ad375bc2fed1f1a28a1216601


		>>b listsinceblock 0000000000180e12c6d8adce69ed14b68016acb5ad884921fcefd94260a908d7 6 true false
		returns the "lastblock" fields which is used for the next scan (in this case 6 back from tip)
		shows stuff in the mempool!


	*/
	upto := int64(1665123)

	client, err := BTCTestNetConnector()
	assert.Nil(t, err, "Error should be nil")
	blockhash, err := client.GetBlockHash(upto)
	assert.NotNil(t, blockhash, "blockhash should not be nil")

	res, err := client.ListSinceBlock(blockhash, 6, true)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, res, "blockhash should not be nil")

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, res.LastBlock, "lastBlockHash should be nil")

}

func Test_GetBlockCount(t *testing.T) {
	client, err := BTCTestNetConnector()
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, client, "Client should not be nil")

	blockCount, err := client.GetBlockCount()
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, blockCount > 1600000, "Block Count invalid (low)")
	assert.True(t, blockCount < 9900000, "Block Count invalid (high)")

}

func Test_GetBlockHash(t *testing.T) {
	client, err := BTCTestNetConnector()
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, client, "Client should not be nil")
	blockhash, err := client.GetBlockHash(10)
	assert.Nil(t, err, "Error should be nil")
	blockhashString := blockhash.String()
	assert.Equal(t, blockhashString, "00000000700e92a916b46b8b91a14d1303d5d91ef0b09eecc3151fb958fd9a2e", "Invalid blockhash for height 10")
}

// func Test_ListSinceBlock(t *testing.T) {
// 	client, err := BTCTestNetConnector()
// 	assert.Nil(t, err, "Error should be nil")

// 	blockHash, err := chainhash.NewHashFromStr("000000000000011f443d9b796854152ff874867580786f7271a1f9855377e13c")

// 	nextBlockHash, newTXCount, err := client.ProcessRecentTransactions(blockHash, 6)
// 	assert.Nil(t, err, "Error should be nil")
// 	assert.Nil(t, err, "Error should be nil")
// 	fmt.Println("New Mature Transactions Count ", newTXCount)
// 	fmt.Println("Next BlockHash ", nextBlockHash)

// }

func BTCTestNetConnector() (*UnderlyingConnector, error) {
	nc, err := qredochain.NewNodeConnector("127.0.0.1:26657", "NODEID", nil, nil)
	if err != nil {
		return nil, err
	}
	client, err := NewUnderlyingConnector(TestnetHost, TestnetUser, TestnetPass)

	connector := UnderlyingConnector{
		client,
		nc,
	}
	return &connector, nil

}
