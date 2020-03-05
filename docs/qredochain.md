# Qredochain

<style type='text/css'>
body {
    counter-reset: h1
}

h1 {
    counter-reset: h2
}

h2 {
    counter-reset: h3
}

h3 {
    counter-reset: h4
}

h1:before {
    counter-increment: h1;
    content: counter(h1) ". "
}

h2:before {
    counter-increment: h2;
    content: counter(h1) "." counter(h2) ". "
}

h3:before {
    counter-increment: h3;
    content: counter(h1) "." counter(h2) "." counter(h3) ". "
}

h4:before {
    counter-increment: h4;
    content: counter(h1) "." counter(h2) "." counter(h3) "." counter(h4) ". "
}</style>


<div style="width: 960px; height: 720px; margin: 10px; position: relative;"><iframe allowfullscreen frameborder="0" style="width:960px; height:720px" src="https://www.lucidchart.com/documents/embeddedchart/a5538bbd-5613-42b9-8830-7dca81439d14" id="ATsDLcAtcDZ8"></iframe></div>




note: I use the BTC chain as an example throughout this document, BTC is the first implemented external cryptocurrency, however the addition of other coins is planned once BTC has been completed.

# Introduction


Qredochain is a blockchain based on Tendermint. There is very little linkage between Qredochain and Tendermint. All consensus rule logic and processes for the Transactions stored in the Tendermint chain are handled by the 'Assets Library' a Golang library which incorporates the Protobuffer definitions and functionality for creating, pasrsing, transferring and validating all transactions.  Persistent Consensus Rule data such as current Wallet balances is stored in a Badger Key/Value database, the internal Tendermint KV database is only used minimally.

Qredochain  is effectively a side chain which temporarily and safely captures Cryptocurrency from other blockchains and facilitates rapid and cheap transfers. Additionally, it enables the attachment of a range of conditions which requires pre-specified parties to authorise a transfer, by way of a BLS Signature before it is actioned and accept on chain.

The Peg-in and Peg-out functionality is handled by a Multi Party Computation Cluster. These MPC nodes together with the Watcher use the Qredochain as their trigger and source of truth to create Addresses in external blockchains, and to sign Settlement transactions to transfer Cryptrocurrency out of the previously created addresses.

The Watcher is a standlance service which mediates communication between teh Qredochain, External Blockchains and the MPC Cluster. The watcher service, if compromised must not allow the theft of funds, however a denial of service is permissable. 


As funds are deposited into the generated Addresses and reach the pre-requites number of confirmations, a watcher service creates a funding transaction into the Qredochain, which adds the incoming funds to the Wallet Asset associated with that External Address. This Peg-in transaction contains an SPV proof which is checked by each node in the Qredochain to ensure the external transaction is genuine.

Similarly a Peg-out transaction releases the locked-up funds from a Peg-In back to a new address on the  underlying Blockchain. 



---
# Contents
   
1. Header
1. Introduction
1. Contents
1. Assets
    1. Signed Asset 
    1. Asset
    1. Consensus immutable/mutable
1. Payloads    
    1. IDDoc
    1. Group
    1. Wallet
    1. Peg-In (Underlying Transaction)
    1. MPC - Mapping
    1. Peg-Out
    1. KVAsset
1. Crystalisation
1. Example Usage

---

# Assets

## Signed Asset - Outer wrapper to hold the signature

Every asset on the Qredochain is contained within an outer wrapper (PBSignedAsset) defined in protobuffers below.

```
message PBSignedAsset{
    bytes Signature            = 1; # BLS (aggregate) Signature
    map<string, bytes> Signers = 3; # "abbreviation":IDDoc ID
    PBAsset Asset              = 4; # Asset (see #2)
}
```

The PBSignedAsset holds the signature of the serialized Asset (see Section #2 below)


Each of the Signers(Field #3) creates a BLS Signature of the serialised PBAsset(Field #4), these signatures are aggregated into a single Signature(Field #1).


The Signers are represented by the AssetID of their Identity Document (IDDoc), which is a type of Asset. Signers can be associated (int the map<string, bytes> Signers ) with an Abbreviation String, this aids readability when the Identity is used in an expression. 

Consensus Rule: Abbreviations must be unique and alphanumeric:  ^[a-zA-Z0-9]+$

A signature can be verified by obtaining the IDDoc for each Signer(#), which contained their BLS Public Key, each of their keys is BLS added together to create an aggregate Public Key (Apk) 
	
   ```bls_verify( serialized_asset, Apk, signature)```


---

## Asset - Wrapper to contain all Types of Qredochain transactions

A PBAsset is a further container which wraps every type of Asset. It contains information appropriate to All Assets, including its ID (a hash of the Payload)

```
message PBAsset {
    PBAssetType Type                      = 1;   //Enum representing the type of payload
    bytes ID                              = 2;   //Hash of the Payload - used as a unique identifier
    bytes Owner                           = 3;   //The current ID doc of the Owner 
    int64 Index                           = 4;   //Starting at 1, incremeents with each update of the asset                
    PBTransferType TransferType           = 5;   //If an update to an existing asset (ie. Index >1), this holds the type of update
    map<string, PBTransfer> Transferlist  = 6;   //A list of Transfer Rules
    map<string,bytes> Tags                = 7;   //Key/Values associted (unused?)
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
    

1. A mutable Asset, such as a Wallet, each update has the same AssetID(#2), but an incrementing Index(#4). If a mutable Asset is create without a TransferList (ie. No instructions on how to do an update) it is impossible to update it.


The Asset container essentially facilitates the updating of existing Transactions. eg. A Wallet Asset which transfers a Cryptocurrency to another Wallet Asset, will be updated with a new version detailing the transfer of currency together with sufficient signatures to make the Update valid. 
All updated Transactions contain a further set of rules which define the required signatures for the next update, they can be the same as both or completely different. 

---

## Consensus Rules - Different types of Assets

### Immutable Asset
When a immutable asset is posted to the chain, it is checked to ensure it adheres to a number of rules
Mandatory fields are not empty
The aggregated Signature verifies against an aggregated Public key of  Signer(s)


### Mutable Asset
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
A concrete example of a PBTransfer: 

```
message PBTransfer {               
    PBTransferType 			= Settlement_Enum
    Expression              	= ‘(T1+T2+T3)>1 & P’
    Participants 			= 
				abbrev:’T1’, iddoc:’ASSETID_IDDOC_FOR_USER_T1’
				abbrev:’T2’,’iddoc:'ASSETID_IDDOC_FOR_USER_T2’
				abbrev:’T3’,’iddoc:’ASSETID_IDDOC_FOR_USER_T3’
                abbrev:’P’ ,’iddoc:’ASSETID_IDDOC_FOR_USER_P’
    Description              	= “some textual description”; 
}
```

A User wants to create an Update to the exsiting asset using the PBTransfer specified above, the Update Transaction must:

1. Have the same AssetID
1. Index = Previous Index + 1
1. TransferType = Settlement
1. The signers of the Transfer must be sufficient to make the expression  '(T1+T2+T3)>1 & P' return true. 
    The PBSignedAsset (outer wrapper) will could contain the following field.

      map<string, bytes> Signers = {
           	    abbrev:’T1’, iddoc:’Signature_of_Asset_Created_by_User_T1’
				abbrev:’T3’,’iddoc:’Signature_of_Asset_Created_by_User_T3’
                abbrev:’P’ ,’iddoc:’Signature_of_Asset_Created_by_User_P’
      }
	

As the expression is parsed When the expression is parser - If a signature is available for a supplied Abbreviation the abbreviation is replaced by the Number 1, if it is not available it is replaced with 0.
If we have Signatures for T1, T3 & P, but not for T2. This is substituted as
		‘(1+0+1)>1 & 1’
        
This evaluates to 1 (true), so the signature is valid, and the Transaction passes BLS Signature Verification.


# Payloads
Within each PBAsset is the Payload, which is data specific to the Asset Type


---
## IDDoc - Identity Document

An immutable Identity, the AssetID is the hash of the serialised Asset
This Transaction contains public keys for individual users. All the Keys are derived from a single Seed, the seed is stored by the User, using this seed, they can generate Private Keys for the Public Keys in the Identity Doucment as required.

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
## Trustee Group (Group)


```
message PBGroup {
    PBGroupType Type                = 1;
    map<string, bytes> GroupFields  = 2;
    map<string, bytes> Participants = 3;
    string Description              = 4;
}
```

There is currently only 1 type of PBGroupType, which is a Trustee Group.
A Trustee Group represents a number of Participants together with an expression.

When used in a Transfer Rule which determines whether or not an Asset can be updated, a trustee group can be used interchangably with a IDDoc.
It enables a Group Asset ID to be used in place of a IDDoc Asset ID. This enables mutable Transfer Conditions.
So, in the case of the previous example

	(T1+T2+T3)>1 & P

if we create the following Trustee Group

```
message PBGroup {
    PBGroupType 		= TrusteeGroup;
    GroupFields  		= 	{‘expression’:’(t1+t2+t3)>1’}
    Participants 		=
		abbrev:’T1’, iddoc:’ASSETID_IDDOC_FOR_USER_T1’
		abbrev:’T2’,’iddoc:'ASSETID_IDDOC_FOR_USER_T2’
		abbrev:’T3’,’iddoc:’ASSETID_IDDOC_FOR_USER_T3’
    string Description   = ‘A description of the trustee group’;
}
```
We can now replace the expression with 

	TG1 & P

The primary benefit of this indirect approach is that the Trustee Group can be updated without effecting any Wallets where it has been used, say, for example where an employee who is a trustee of many wallets leaves their employment. A new user is able to be assigned as a trustee of many wallets by updating a single Trustee Group.

---

## Wallet

A Wallet updates funds in the Qredochain. 
A BTC Address produced by the MPC are associated with Wallets, via MPC transactions.
External Funds deposited into that BTC Address is added to the Wallet's balance as store by a Qredochain node via the consensus rules.

WalletTransfers is a set of transaction which move funds from this wallet to other wallets.
When the wallet transaction is committed to the chain,   the transferred funds are deducted from this wallet's balance, and added to other wallets balances.

A settlement update locks the Wallet and prevents any further updates, until the settlement confirmation has been produced by the MPC/Watchers. This confirmation is made when an underlying BTC transaction has been broadcast, the wallet is then unlocked 



```
message PBWallet {
    PBCryptoCurrency Currency                   = 1;
    int64 SpentBalance                          = 2; 
    repeated PBWalletTransfer WalletTransfers   = 3;
    string Description                          = 11;
    bytes Principal                             = 12;
	bytes Creditor                              = 13;
	bytes Initiator                             = 14;
	bytes Address                               = 15;
    bytes Counterparty                          = 16;
    bytes TransactionHash                       = 17;
}
```


---

## Peg-In (Underlying)

A Peg-In (alias: Underlying) transaction is made by the watcher to indicate the creation of a new UTXO whose address has a mapping on the Qredochain. 
As the Peg-In transaction is commited to the chain, the consensus rules, (after determining that the UTXO hasn't been previous processed), adds the incoming amount to the Wallet's balance.

Any Qredochain node needs to be sure that the funds specified in the transaction are genuine, the SPV Proof can be verified against a standard Bitcoin Node.

```
message PBUnderlying {
    PBUnderlyingType Type               = 1; // Type
    PBCryptoCurrency CryptoCurrencyCode = 2; // Which cryptocurrency chain this relates to.
    bytes Proof                         = 3; // Merkle proof of the transaction - allows nodes to easily determine validity
    int64 Amount                        = 4; // Value of the transaction  
    bytes Address                       = 5; // Address money was sent to.
    bytes TxID                          = 6; // Credits this Asset ID
}
```

---

## MPC 

An MPC transaction is generated by the MPC Cluster, it creates a mapping between underlying BTC Addresses and their beneficiary Qredochain Wallet (AssetID).

Bitcoin UTXO >> MPC >> Wallet


```
message PBMPC {
    PBMPCType   Type        = 1;
    bytes       Address     = 2;
    bytes       Signature   = 3;
    bytes       AssetID     = 4;
}
```

---

## Peg-Out

Peg-Out is a MPC generated transaction upon the broadcast of the settlement transaction. It confirms that a Bitcoin Node has accept the underlying BTC Transaction, it reduces the balance of the Qredochain Wallet, and unlocks it to enable further transfers

---

## KVAsset


A KVAsset is simply a wrapper around a set of Keys & Values.
The immutable field contains a set of keys which can't be changed in any updates. This rule is enforced as a consensus rule. 


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

When a user wishes to settle out of the Qredo system by using a settlement transaction, the underlying chain eg. Bitcoin Chain, requires a fee to be paid. As the underlying chain's fee changes the amount a user needs to pay will vary. This fee can be communicated to users within the Qredo system be looking up a value in a KVAsset. The value would represent a fee as Satoshis/Byte, and parties within the system can determine the appropriate fee.

The Public Keys for the MPC nodes, the Watcher service, and any other permissioned services on the Qredo Network will have entries in a KVAsset, attached to theses Assets are conditions which will require a number of Signatures before they can be updated.


---

# Crystalisation

The values of each  UTXOs once added to the Qredochain do not remain tied to the Wallet they originally fund.  At any point a crystalisation process can be performed, mapping UTXOs on an underlying chain to Balances stored in a Qredochain Wallet.
The sum of all UTXOs will match the sum of all Qredochain Wallet balances.
The mapping process is deterministic, so any Qredochain given a specific blockheight will generate the same set of relationships.
The goal of such a process is to generate a set of unsigned transactions to allow a user to prove the existence and whereabouts of their underlying assets. As Qredo generates an MPC mapping transaction, it also signs the BTC addres to provide a 'proof of funds'
When a settlement of the Assets is required, the transactions generated by the crystalisation process are sent to the MPC Cluster for signing, (any change is returned to an MPC issued BTC Address, and used in any future crystalization processes.

---

# Example Usage

Here we walk though a typical workflow purely from the Qredochain Transaction point of view. A majority of the communiction between parties, such as obtaining signatures from other trustees is handled by the Matrix Communication Protocol.


## Funding - Addings fund to your account
1. Alice, Bob & Charlie & Dave create their IDDocument, each IDDoc is added to the Qredochain and assigned an AssetID
1. Alice creates a Wallet Asset of type Bitcoin (BTC), its added to the Qredochain and assigned AssetID (Wallet_Alice). The Wallet has a transfer rule which requires either Charle or Dave to additional sign any updates to the Wallet (expression: `Charlie | Dave & Alice')
1. The Watcher detects new Wallet_Alice and requests the MPC to issue a new BTC Address based on the Asset_ID.
1. MPC creates a BTC_Address_A and a Qredochain transaction mapping the BTC Address to Wallet_Alice (MPC_A)
1. BTC_Address_A is added as a watch only address on the watchers BTC Node. From this point forward the BTC Node will monitor the address for changes. 
1. Externally Anne deposits BTC into BTC_Address_A
1. The watcher periodically (every minute) requests wallet changes from its BTC Node. Where new funds have been added with sufficient confirmations, it creates a PEG-IN transaction on the Qredochain. 
1. As this PEG-IN (underlying) transaction is commited to the chain, Qredochain consensus rules uses MPC_A to provide the mapping from BTC_Address_A to Wallet_Alice, and the amount of BTC deposited into the address is added to the Wallet_Alice BTC balance, which is a value in the consensus database with key  "Wallet_Alice.balance"
1. Funding this address can be repeated at will, and simply increases the balance in the database key.


## Spending - Sending funds to other parties
1. Alice arranges to send Bob some BTC
1. Alice creates an update Wallet_A transaction, which transfers some of the BTC to another Wallet_Bob, the transaction is signed by Alice and either Charlie or Dave, and the signatures aggregated into a single signature. The update is broadcast to the QredoChain
1. As the Wallet Update is commited to the Chain, the Qredo consensus rules deduct the transferred amount from Wallet_Alice.Balance key and add it to Wallet_Bob.Balance.
1. Assets can be transferred between parties at will, time to transact is dependant on the commitment of a single transaction to the Qredochain, which is currently around 1 second.



## Settlement - Getting funds out of Qredo
1. Alice constructs an update to Wallet_Alice, the type is settlement, and the required signatures from either Dave or Charlie are obtained, to make it valid.
1. The Update is committed to the Qredochain, locking the Wallet from any further updates
1. The watcher detects a settlement, crystalizes the chain, and sends the resultant transaction to the MPC for signing.
1. The MPC queries any Qredochain node to confirm the existence, destination and amounts of the Wallet_Alice settlement.
1. The MPC creates a new Address to accept the unspent funds/change from the transaction, this is added to the BTC_Node for monitoring, and a MPC transaction is add to the Qredochain, where the Wallet mapping is empty.
1. The MPC signs the transactions and broadcasts them to the underlying blockchain 
1. The MPC creates a PEG-OUT Qredochain transaction finalizing the settlement, which updates the wallet balance and unlocks the wallet, allowing further updates.







