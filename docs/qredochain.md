# Qredochain



<div style="width: 960px; height: 720px; margin: 10px; position: relative;"><iframe allowfullscreen frameborder="0" style="width:960px; height:720px" src="https://www.lucidchart.com/documents/embeddedchart/a5538bbd-5613-42b9-8830-7dca81439d14" id="ATsDLcAtcDZ8"></iframe></div>


Qredochain is a blockchain based on Tendermint. 
It is effectively a side chain which temporarily captures Coins from other blockchains and facilitates rapid and cheap transfers. Additionally, it enables the attachment of a range of conditions which requires pre-specified parties to authorise a transfer, by way of a BLS Signature before it is actioned and accept on chain.

The Peg-in and Peg-out functionality is handled by a Multi Party Computation Cluster. The MPC nodes together with ther Watcher use the Qredochain as their trigger to create Addresses in external blockchains (Peg-in), and to sign Settlement transactions to transfer Cryptrocurrency out of the Qredochain (Peg-out).

As funds are deposited into the Peg-In addresses and reach the pre-requites number of confirmations, a watcher service creates a funding transaction into the Qredochain, which adds the incoming funds to the Wallet Asset associated with that External Address. This Peg-in transaction contains an SPV proof which is checked by each node in the Qredochain to ensure the external transaction is genuine.

Similarly a Peg-out transaction releases the locked-up funds from a Peg-In back to a new address on the  underlying Blockchain. 

---


1. Signed Asset
1. Asset
1. IDDoc
1. Group
1. Wallet
1. Peg-In (Underlying Transaction)
1. MPC - Mapping
1. Peg-Out


# Signed Asset

Every asset on the Qredochain is contained within an outer wrapper (PBSignedAsset. ) defined in protobuffers below.

```
message PBSignedAsset{
    bytes Signature            = 1;
    map<string, bytes> Signers = 3;
    PBAsset Asset              = 4;
}
```

The PBSignedAsset facilitates the signing of every transaction on the chain.

Each of the Signers(#3) creates a BLS Signature of the serialised PBAsset(4), these signatures are aggregated into a single Signature(1).

The Signers are represented by the AssetID of their Identity Document (IDDoc), which is a type of Asset. 
A signature can be verified by obtaining the IDDoc for each Signer(#), which contained their BLS Public Key, each of their keys is BLS added together to create an aggregate Public Key (Apk) 
	
   ```bls_verify( serialized_asset, Apk, signature)```


---

# Asset

An PBAsset(#4) is a further container which wraps every type of Asset. 

```
message PBAsset {
    PBAssetType Type                      = 1;   
    bytes ID                              = 2;  
    bytes Owner                           = 3;  
    int64 Index                           = 4;                  
    PBTransferType TransferType           = 5;
    map<string, PBTransfer> Transferlist  = 6;
    map<string,bytes> Tags                = 7; 
    oneof Payload {                          
        PBWallet Wallet                   = 15;
        PBGroup Group                     = 16;
        PBIDDoc Iddoc                     = 17;
        PBUnderlying Underlying           = 18;
        PBKVAsset KVAsset                 = 19;
        PBMPC MPC                         = 20;
    }
}
```


There are two types of Assets

1. An immutable Asset, such as an IDDoc, which can’t be updated and only has a single version, Index(#4) is always equal to 1

1. A mutable Asset, such as a Wallet, each update has the same AssetID(#2), but an incrementing Index(#4).


The Asset container essentially facilitates the updating of existing Transactions. eg. A Wallet Asset which transfers a Cryptocurrency to another Wallet Asset, will be updated with a new version detailing the transfer of currency together with sufficient signatures to make the Update valid. 
All updated Transactions contain a further set of rules which define the required signatures for the next update, they can be the same as both or completely different. 

---

# Consensus Rules

## Mutable Asset
When a immutable asset is posted to the chain, it is checked to ensure it adheres to a number of rules
Mandatory fields are not empty
The aggregated Signature verifies against an aggregated Public key of  Signer(s)


## Immutable Asset
A mutable asset has a number of addition consensus rules
1) It must include a set of TransferRules, which defines who is required to Sign an Updated version of the Asset before it is accepted into the chain.

```
message PBTransfer {               
    PBTransferType Type             = 1;	
    string Expression               = 2; 
    map<string, bytes> Participants = 3;
    string Description              = 4; 
}
```

The Participants fields is represented by an Abbreviation and a IDDoc Asset ID.
Example of a PBTransfer :

```
message PBTransfer {               
    PBTransferType 			= Settlement_Enum
    Expression              	= ‘(t1+t2+t3)>1 & p’
    Participants 			= 
				abbrev:’t1’, iddoc:’ASSETID_IDDOC_FOR_USER_T1’
				abbrev:’t2’,’iddoc:'ASSETID_IDDOC_FOR_USER_T2’
				abbrev:’t3’,’iddoc:’ASSETID_IDDOC_FOR_USER_T3’
    Description              	= “some textual description”; 
}
```

An updated version of an Immutable Asset must: (using the Transfer above)
1. Have the same AssetID
1. Index = Previous Index + 1
1. TransferType = Settlement
1. The signers of the Transfer must be sufficient to make the expression return true. 
	
    	(t1+t2+t3)>1 & p 

When the expression is parse - If a signature is available for a supplied Abbreviation the abbreviation is replaced by the Number 1, if it is not available it is replaced with 0.
If we have Signatures for t1, t2 & p, but not for t3. This is substituted as
		‘(1+1+0)>1 & 1’ 

which evaluates to 1 (true), so the Update transaction has the required signatures to make it valid.


Asset Payload
Within each PBAsset is the Payload, which is data specific to the Asset Type

The following assets are available







---
## IDDoc

An immutable Identity, the AssetID is the hash of the serialised Asset

```
message PBIDDoc {
    string AuthenticationReference = 1;
    bytes BeneficiaryECPublicKey   = 2;
    bytes SikePublicKey            = 3;
    bytes BLSPublicKey             = 4;
    int64 Timestamp                = 5;
}
```




---
## Trustee Group

(There is currently only 1 type of PBGroupType, which is a Trustee Group.)

```
message PBGroup {
    PBGroupType Type                = 1;
    map<string, bytes> GroupFields  = 2;
    map<string, bytes> Participants = 3;
    string Description              = 4;
}
```

A Group represents a number of Participants together with an expression


A Trustee Group Transaction enables a Group Asset ID to be used in place of a IDDoc Asset ID. This enables mutable Transfer Conditions.
So, in the case of the previous example

	‘(T1+T2+T3)>1 & P’


We can represent ‘(T1+T2+T3)>1’ in a Trustee Group (TG1)

```
message PBGroup {
    PBGroupType 		= TrusteeGroup;
    GroupFields  		= 	{‘expression’:’(t1+t2+t3)>1’
    Participants 		=
		abbrev:’T1’, iddoc:’ASSETID_IDDOC_FOR_USER_T1’
		abbrev:’T2’,’iddoc:'ASSETID_IDDOC_FOR_USER_T2’
		abbrev:’T3’,’iddoc:’ASSETID_IDDOC_FOR_USER_T3’

    string Description   = ‘A description of the trustee group’;
}
```

			
Where TG1 is the asset ID of the Trustee Group, we can now write the expression 
 
    (T1+T2+T3)>1 & P
    
 as 

	TG1 & p

The Trustee Group can be updated using Transfer Rules implemented in the same was as a Wallet

---

## KVAsset

A mutable set of Keys & Values for general purpose use.

```
message  PBKVAsset {
    PBKVAssetType Type              = 1;
    map<string, bytes> AssetFields  = 2; 
    repeated string Immutable       = 3;
    string Description              = 4; 
}
```

Example usage:

Settlement Underlying Fees.

When a user wishes to removed money out of the Qredo system by using a settlement transaction, the underlying chain eg. Bitcoin Chain, requires a fee to be paid. This fee needs to be updated periodically as the underlying chain its fee requirements.

A Fee KVAsset is used 


## Crystalisation

The values of each  UTXOs once added to the Qredochain do not remain tied to the Wallet they originally fund.  At any point a crystalisation process can be performed, mapping UTXOs on an underlying chain to Balances stored in a Qredochain Wallet.
The sum of all UTXOs will match the sum of all Qredochain Wallet balances.
The mapping process is deterministic, so any Qredochain given a specific blockheight will generate the same set of relationships.
The goal of such a process is to generate a set of unsigned transactions to allow a user to prove the existence and whereabouts of their underlying assets. As Qredo generates an MPC mapping transaction, it also signs the BTC addres to provide a 'proof of funds'
When a settlement of the Assets is required, the transactions generated by the crystalisation process are sent to the MPC Cluster for signing, (any change is returned to an MPC issued BTC Address, and used in any future crystalization processes.



## Example Usage

### Funding - Addings fund to your account
1. Alice, Bob & Charlie & Dave create their IDDocument, each IDDoc is added to the Qredochain and assigned an AssetID
1. Alice creates a Wallet Asset of type Bitcoin (BTC), its added to the Qredochain and assigned AssetID (Wallet_Alice). The Wallet has a transfer rule which requires either Charle or Dave to additional sign any updates to the Wallet (expression: `Charlie | Dave & Alice')
1. The Watcher detects new Wallet_Alice and requests the MPC to issue a new BTC Address based on the Asset_ID.
1. MPC creates a BTC_Address_A and a Qredochain transaction mapping the BTC Address to Wallet_Alice (MPC_A)
1. BTC_Address_A is added as a watch only address on the watchers BTC Node. From this point forward the BTC Node will monitor the address for changes. 
1. Externally Anne deposits BTC into BTC_Address_A
1. The watcher periodically (every minute) requests wallet changes from its BTC Node. Where new funds have been added with sufficient confirmations, it creates a PEG-IN transaction on the Qredochain. 
1. As this PEG-IN (underlying) transaction is commited to the chain, Qredochain consensus rules uses MPC_A to provide the mapping from BTC_Address_A to Wallet_Alice, and the amount of BTC deposited into the address is added to the Wallet_Alice BTC balance, which is a value in the consensus database with key  "Wallet_Alice.balance"
1. Funding this address can be repeated at will, and simply increases the balance in the database key.


### Spending - Sending funds to other parties
1. Alice arranges to send Bob some BTC
1. Alice creates an update Wallet_A transaction, which transfers some of the BTC to another Wallet_Bob, the transaction is signed by Alice and either Charlie or Dave, and the signatures aggregated into a single signature. The update is broadcast to the QredoChain
1. As the Wallet Update is commited to the Chain, the Qredo consensus rules deduct the transferred amount from Wallet_Alice.Balance key and add it to Wallet_Bob.Balance.
1. Assets can be transferred between parties at will, time to transact is dependant on the commitment of a single transaction to the Qredochain, which is currently around 1 second.





### Settlement - Getting funds out of Qredo
1. Alice constructs an update to Wallet_Alice, the type is settlement, and the required signatures from either Dave or Charlie are obtained, to make it valid.
1. The Update is committed to the Qredochain, locking the Wallet from any further updates
1. The watcher detects a settlement, crystalizes the chain, and sends the resultant transaction to the MPC for signing.
1. The MPC queries any Qredochain node to confirm the existence, destination and amounts of the Wallet_Alice settlement.
1. The MPC creates a new Address to accept the unspent funds/change from the transaction, this is added to the BTC_Node for monitoring, and a MPC transaction is add to the Qredochain, where the Wallet mapping is empty.
1. The MPC signs the transactions and broadcasts them to the underlying blockchain 
1. The MPC creates a PEG-OUT Qredochain transaction finalizing the settlement, which updates the wallet balance and unlocks the wallet, allowing further updates.







