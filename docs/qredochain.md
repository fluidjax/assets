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

# Qredochain


<!-- table made https://ecotrust-canada.github.io/markdown-toc/ -->


- [Qredochain](#qredochain)
- [Introduction](#introduction)
- [Assets](#assets)
  * [Signed Asset - Outer wrapper to hold the signature](#signed-asset---outer-wrapper-to-hold-the-signature)
  * [Asset - Wrapper to contain all Types of Qredochain transactions](#asset---wrapper-to-contain-all-types-of-qredochain-transactions)
  * [Consensus Rules - Different types of Assets](#consensus-rules---different-types-of-assets)
    + [New Asset (Immutable & Mutable Assets)](#new-asset--immutable---mutable-assets-)
    + [Updated Asset (Mutable Asset)](#updated-asset--mutable-asset-)
    + [Update - TransferRules](#update---transferrules)
    + [Transfer example](#transfer-example)
- [Payloads](#payloads)
  * [IDDoc - Identity Document](#iddoc---identity-document)
  * [Trustee Group (Group)](#trustee-group--group-)
  * [Wallet](#wallet)
    + [Setup](#setup)
    + [Adding Funds](#adding-funds)
    + [Reducing Funds](#reducing-funds)
    + [Fees](#fees)
  * [Peg-In (Underlying)](#peg-in--underlying-)
  * [MPC](#mpc)
  * [Peg-Out](#peg-out)
  * [KVAsset](#kvasset)
    + [Example usage:](#example-usage-)
      - [Settlement Fees](#settlement-fees)
      - [Pre-Defined Members of the Qredo Network](#pre-defined-members-of-the-qredo-network)
    + [Consensus Rules](#consensus-rules)
- [Crystalization](#crystalization)
- [Example Usage](#example-usage)
  * [Funding - Adding fund to your account](#funding---adding-fund-to-your-account)
  * [Spending - Sending funds to other parties](#spending---sending-funds-to-other-parties)
  * [Settlement - Getting funds out of Qredo](#settlement---getting-funds-out-of-qredo)
- [Watcher](#watcher)
  * [Wallet Creation](#wallet-creation)
  * [Underlying Funding Transaction](#underlying-funding-transaction)
  * [Settlement](#settlement)


---


# Introduction

<div style="width: 960px; height: 720px; margin: 10px; position: relative;"><iframe allowfullscreen frameborder="0" style="width:960px; height:720px" src="https://www.lucidchart.com/documents/embeddedchart/a5538bbd-5613-42b9-8830-7dca81439d14" id="ATsDLcAtcDZ8"></iframe></div>

 - Document Source: https://github.com/qredo/assets/blob/master/docs/qredochain.md
 - Assets Repository: https://github.com/qredo/assets
 - Author: Chris Morris chris@qredo.com


An overview of the Qredo System, this document covers the "Tendermint Node - Qredochain" and the "Watcher"


note: I use the BTC chain as an external CryptoCurrency throughout this document, BTC is the first implemented external cryptocurrency, however the addition of other coins is planned during later phases of development.

## What is Qredochain:

The name 'Qredochain' is the name of the Application built on top of Tendermint. It encompasses both the Transactions on the Tendermint Blockchain,  together with the Consensus rules for including new blocks into the chain. Qredochain also acts as a side-chain to the Cryptocurrencies it interfaces with,  It temporarily and safely captures Cryptocurrency from other blockchains and facilitates rapid and cheap transfers between other Qredochain users. Additionally, it enables the attachment of a range of conditions which requires pre-specified parties to authorise a transfer before it is accepted into the Qredochain. Authorisation to perform transfer is given by individual users signing the transaction with a BLS Signature, these signatures are then aggregated into a single signature and attached to each transfer transaction.

There is very little linkage between Qredochain and Tendermint. All consensus rule logic and processes for the Transactions stored in the Tendermint chain are handled by the 'Assets Library', a Golang library which incorporates the Protobuffer definitions and functionality for creating, parsing, transferring and validating all transactions.  Persistent Consensus Rule data such as current Wallet balances are stored in a Badger Key/Value database, the internal Tendermint KV database is only used minimally.

Qredochain works on the principal of a ‘pegged sidechain’, Assets are transferred into the Qredochain using Peg-in, and similarly transferred out of Qredochain using the Peg-Out mechanism. All assets on the Qredochain are created as assets on an another Cryptocurrency chain are pegged-in and locked up. Assets transferred out of the Qredochain during the settlement process are destroyed on the Qredochain and, are released back on the original Cryptocurrency chain. The Peg-in and Peg-out functionality is managed by the Multi Party Computation Cluster.

 These MPC nodes together with the Watcher use the Qredochain as their trigger and source of truth to create Addresses in external blockchains. When Cryptocurrency is to be transferred back out of the system to an external address, a Settlement Update Transactions is added to a wallet, and the MPC nodes again sign a Peg-out transaction to release the locked up funds.


## The Watcher
The Watcher is a standalone service which mediates communication between the Qredochain, External Blockchains and the MPC Cluster. The watcher service, if compromised must not allow the theft of funds, however a denial of service is permissible. 


As funds are deposited into the generated Addresses and reach the pre-requites number of confirmations, a watcher service creates a funding transaction into the Qredochain, which adds the incoming funds to the Wallet Asset associated with that External Address. This Peg-in transaction contains an SPV proof which is checked by each node in the Qredochain to ensure the external transaction is genuine.

Similarly a Peg-out transaction releases the locked-up funds from a Peg-In back to a new address on the  underlying Blockchain. 

All Qredochain Assets are encoded using Google's Protobuffer. 
The full Protobuffer definition is within the assets repo @  https://github.com/qredo/assets/blob/master/libs/protobuffer/assets.proto



# Assets

## Signed Asset - Outer wrapper to hold the signature

Every asset on the Qredochain is contained within an outer wrapper (PBSignedAsset) defined in Protobuffer definition below.

```
message PBSignedAsset{
    bytes Signature            = 1; # BLS (aggregate) Signature
    map<string, bytes> Signers = 3; # "abbreviation":IDDoc ID
    PBAsset Asset              = 4; # Asset (see #2)
}
```

The PBSignedAsset holds the signature of the serialized Asset (see Section #2 below)


Each of the Signers(Field #3) creates a BLS Signature of the serialized PBAsset(Field #4), these signatures are aggregated into a single Signature(Field #1).


The Signers are represented by the AssetID of their Identity Document (IDDoc), which is also type of Asset on the Qredochain. Each signers can be associated (int the map<string, bytes> Signers ) with an Abbreviation String, this aids readability when the Identity is used in an expression. For example, instead of a long hex string, trustees could be labelled as T1, T2 and Principals as P,.

Consensus Rule: Abbreviations must be unique and alphanumeric:  ^[a-zA-Z0-9]+$

A signature can be verified by obtaining the IDDoc for each Signer(Field 3.bytes), which contains their BLS Public Key, each of their keys BLS added together to create a single aggregate Public Key 
    
 The signature (#Field 1) is an aggregation of all the signatures from each of the signers.

---

## Asset - Wrapper to contain all Types of Qredochain transactions

A PBAsset is a container within a PBSignedAsset. It contains generic information common to each Asset, including its Asset ID, type, index and Transfer information which allows the Asset to be updated to new versions.

```
message PBAsset {
    PBAssetType Type                      = 1;   //Enum representing the type of payload
    bytes ID                              = 2;   //Random ID
    bytes Owner                           = 3;   //The current ID doc of the Owner 
    int64 Index                           = 4;   //Starting at 1, increments with each update of the asset                
    PBTransferType TransferType           = 5;   //If an update to an existing asset (ie. Index >1), this holds the type of update
    map<string, PBTransfer> Transferlist  = 6;   //A list of Transfer Rules
    map<string,bytes> Tags                = 7;   //Key/Values associated (unused?)
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


The Asset container  facilitates the updating of existing Transactions. eg. A Wallet Asset which transfers a Cryptocurrency to another Wallet Asset, will be updated with a new version detailing the transfer of currency together with sufficient signatures to make the Update valid. 
All updated Assets transactions contain a further set of rules which define the required signatures for the next update, they can be the same as the previous Transfer Rules or completely different. 

The ID or AssetID, is a cryptographically random 256bit number, generated by the initiator.

---

## Consensus Rules - Different types of Assets

Any transaction is either a New Asset, where the AssetID doesn't already exist, or, where permitted (ie. mutable), it is an Updated Asset, where it updates a previous Asset with new data. An 'Updated Asset' can update both a previous 'Updated Asset' or could be the first update in the chain, updating the first 'New Asset'

 ```"New_Asset" > "Updated_Asset" > "Updated_Asset" >  "Updated_Asset" ....```

### New Asset (Immutable & Mutable Assets)
Both mutable and immutable assets have a 'New Asset' transaction type. A 'New Asset' transaction is the first transaction for any given AssetID.
A number of consensus rules are applied to every 'New Asset'
The primary objective of the Consensus Ruleset is to ensure security.
The rules which are checked before each Transaction is added to the Chain is the core mechanism for the security of the system. 

Mandatory fields are not empty
The aggregated Signature verifies against an aggregated Public key of  Signer(s)

1. Transaction does not already exist.
1. AssetID **does not** already exist.
1. Asset Index equals 1
1. Assets are correctly formatted - contain all mandatory fields
1. Based on the self declared signers, the signature verifies - The composition of declare signers is not a consensus rule, it could be simply the Principal (Owner of the Wallet) or an aggregated signature of every participant in the Expressions. This is left to the User facing application to enforce. Additional signers at this stage do not offer any additional security, but may be required as part of an audit trail.

"New Assets" have additional checks specific to their types (eg. A wallet has checks relating to their balances) these are detailed in their own sections below.



### Updated Asset (Mutable Asset)
A mutable asset has a number of addition consensus rules
To check an "Updated Asset" the previous version of the transaction is obtained and rebuilt into an Asset object, it is then used together with the new "Updated Asset" to perform a number of checks.

1. Transaction does not already exist.
1. AssetID **does** already exist.
1. Asset Index equals the previous + 1
1. Immutable fields in a Mutable Asset
    1. PBAssetType
    1. AssetID
    1. TransferRules - except the rule used in the CurrentUpdate
1. Assets are correctly formatted - contain all mandatory fields
1. Based on the self declared signers, the signature verifies 
1. The declared signers are sufficient to return true when their signatures are used in the expression from the previous version of the Asset.


### Update - TransferRules

A mutable asset, either a "New Asset" or "Updated Asset" requires TransferRules which defines who is required to Sign any future "Updated Asset" to make it valid. TransferRules take the form:


```
message PBTransfer {               
    PBTransferType Type             = 1;    
    string Expression               = 2; 
    map<string, bytes> Participants = 3;
    string Description              = 4; 
}
```

Any number of transfer rules can be attached to 


### Transfer example

Here is a  concrete example with the Transfer fields completed.

```
message PBTransfer {               
    PBTransferType      = Settlement_Enum
    Expression          = ‘(T1+T2+T3)>1 & P’
    Participants        = 
        abbrev:’T1’, iddoc:’ASSETID_IDDOC_FOR_USER_T1’
        abbrev:’T2’,’iddoc:'ASSETID_IDDOC_FOR_USER_T2’
        abbrev:’T3’,’iddoc:’ASSETID_IDDOC_FOR_USER_T3’
        abbrev:’P’ ,’iddoc:’ASSETID_IDDOC_FOR_USER_P’
    Description         = “some textual description”; 
}
```

The Participants field is represented by an Abbreviation and a IDDoc Asset ID.
A concrete example of a PBTransfer: 

A User who wants to create an Update to the existing asset using the PBTransfer specified above, the Update Transaction must obey all the consensus rule but importantly must fulfil the TransferRule :-

The signers of the Transfer must be sufficient to make the expression  '(T1+T2+T3)>1 & P' return true. 
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
Each Payload of course has all the fields of its PBSignedAsset and PBAsset wrappers, together with a number of fields specific to each type.
Detailed below is each of the Assets types:

---
## IDDoc - Identity Document

This transaction type contains public keys for individual users. All the Keys are derived from a single Seed, the seed is stored by the User, using this seed, they can generate Private Keys for the Public Keys in the Identity Document as required.

A special case for an IDDoc is that the AssetID is the hash of the serialised Asset .(There is sufficient entropy in the public keys)

```
message PBIDDoc {
    string AuthenticationReference = 1;
    bytes BeneficiaryECPublicKey   = 2;
    bytes SikePublicKey            = 3;
    bytes BLSPublicKey             = 4;
    int64 Timestamp                = 5;
}
```

There are no addition consensus rules beyond those specified in the standard "New Asset" Consensus Rules above.

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

When used in a Transfer Rule which determines whether or not an Asset can be updated, a trustee group can be used interchangeably with a IDDoc.
It enables a Group Asset ID to be used in place of a IDDoc Asset ID. This enables mutable Transfer Conditions.
So, in the case of the previous example

    (T1+T2+T3)>1 & P

if we create the following Trustee Group

```
message PBGroup {
    PBGroupType         = TrusteeGroup;
    GroupFields          =     {‘expression’:’(t1+t2+t3)>1’}
    Participants         =
        abbrev:’T1’, iddoc:’ASSETID_IDDOC_FOR_USER_T1’
        abbrev:’T2’,’iddoc:'ASSETID_IDDOC_FOR_USER_T2’
        abbrev:’T3’,’iddoc:’ASSETID_IDDOC_FOR_USER_T3’
    string Description   = ‘A description of the trustee group’;
}
```
We can now replace the expression with 

    TG1 & P

The primary benefit of this indirect approach is that the Trustee Group can be updated without effecting any Wallets where it has been used (or **will** be used), say, for example an employee who is a trustee of many wallets leaves their employment. A new user is able to be assigned as a trustee of many wallets by updating a single Trustee Group.

There are no addition consensus rules beyond those specified in the standard "Updated Asset" Consensus Rules above.

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


Wallets have additional Rules to enable them to retain balances.
In addition to either the Asset Update or Asset Creation rules above, a wallet must also ensure that the balance it stores doesn't lose or create money either by error or malice.

### Setup
1. Upon creation a wallet is assigned a Zero balance, this is stored in the consensus KV store using the key
    [AssetID].balance

### Adding Funds
There are 2 valid ways to add funds to the total stored in a wallet
1. An external funding transaction, when a Peg-In (Underlying) transaction is made by the watcher which notifies the Qredochain that a previous issued address has new funds. The MPC transaction is used to find which Wallet (AssetID) is associated with this external transaction. As the Peg-In transaction is committed to the chain. The AssetID balance is update

    ```[AssetID].balance = [AssetID].balance + incoming_external_funds```

1. A Wallet Update transfers funds to other AssetID using the WalletTransfer field as this transaction is committed the balance of a wallet is incremented by the amount transferred from the other Wallet.

    ```[AssetID].balance = [AssetID].balance + fund_from_other_wallet```

### Reducing Funds
There is one way to reduce the balance of a Wallet, this mechanism can be used to transfer funds to other assetsIDs within the Qredo system, or as part of a Peg-Out settlement transaction, where the funds are removed from the system and underlying funds in the BTC chain are unlocked and transferred to an address of the owners choosing.
1. Wallet transfers, where a wallet update includes a WalletTransfer (specifying the recipient AssetID and the amount), for **each** transfer the balance of the Wallet is reduce by the amount transferred. Before the transaction is declared valid, a check is made to ensure that the sum of add Reducing transaction doesn't leave a balance of less than zero.

    ```[AssetID].balance = [AssetID].balance - funds_transferred_to_other_wallet ```


1. Where a Wallet transfer is of the type 'Settlement', the balance of the wallet is reduced by the funds settled, again a check is made to ensure it doesn't result in a less than zero wallet.

    ```[AssetID].balance = [AssetID].balance - settled_funds```


### Fees
Fees can be added to any Wallet Update by simply adding additional WalletTransfer entries. At present this is automatically implemented by the external programs generating transactions,  However, if mandatory fees are required a KVAsset could be created with 
    1. The AssetID of the Qredo Fee Wallet,
    1. The fee amount.
A new consensus rule added, to ensure that all Wallet Update transactions include a WalletTransfer to Qredo as specified in the KVAsset






---

## Peg-In (Underlying)

A Peg-In (alias: Underlying) transaction is made by the watcher to indicate the creation of a new UTXO whose address has a mapping on the Qredochain. 
As the Peg-In transaction is committed to the chain, the consensus rules, (after determining that the UTXO hasn't been previous processed), adds the incoming amount to the Wallet's balance.

Any Qredochain node needs to be sure that the funds specified in the transaction are genuine, the SPV Proof can be verified against a standard Bitcoin Node.

```
message PBUnderlying {
    PBUnderlyingType Type               = 1; // Type
    PBCryptoCurrency CryptoCurrencyCode = 2; // Which cryptocurrency chain this relates to.
    bytes Proof                         = 3; // Merkle proof of the transaction - allows nodes to easily determine validity
    int64 Amount                        = 4; // Value of the transaction  
    bytes Address                       = 5; // Address money was sent to.
    bytes TxID                          = 6; // The underlying TxID of the external transaction
}
```

Peg-In are immutable Assets which adhere to the standard "New Asset" ruleset above,

1. In addition as the asset is committed to the chain it obtains the MPC transaction (detailed below) which maps the BTC Address (Field #5) in the PBUnderlying message to the AssetID (of the related wallet), and updates the balance of the AssetID (wallet) incrementing it by the amount transferred.
1. A check is made to ensure that the TxID doesn't already exist in the Qredochain. An a new KV field is created where the Key is the TxID, to ensure no further copies of the same external UTXO can be added subsequently.
1. Because the PBUnderlying Transaction can't be trusted, the Qredochain node checks a BTC node to confirm that the details in the Underlying Transaction actually exist at sufficient depth in the BTC Chain. (Mechanism not finalized, but possibly by SPV)


---

## MPC 

An MPC transaction is generated by the MPC Cluster, it creates a mapping between underlying BTC Addresses and their beneficiary Qredochain Wallet (AssetID). It also notifies the Owner of a Wallet of the address where they can deposit external BTC funds, to fund their Qredochain wallet.

The MPC is used by the Underlying (Peg-in) transaction to map incoming funds to Qredochain Wallets.

Bitcoin UTXO >> MPC >> Wallet


```
message PBMPC {
    PBMPCType   Type        = 1;
    bytes       Address     = 2;
    bytes       Signature   = 3;
    bytes       PublicKey   = 4; 
    bytes       Proof       = 5;
 
```


The Proof is a signature generated by the MPC Cluster, using a standard BLS verification function:
```
    bool = verify(msg, pubKey, signature)

    Where:
    - pubKey    = PublicKey(#Field 4)
    - msg       = PublicKey(#Field 4) 
    - signature = Signature (#Field 3)
```

The function returns true, proving that the MPC cluster is able to make valid signature for the supplied private, as if it possesses the private key. (Which by virtue of the properties of the MPC network, it doesn't)

    


The MPC Transaction needs to be signed by members of the MPC network.
These MPC members public keys are held in a pre-defined KVAsset.


---

## Peg-Out

Peg-Out is a MPC generated transaction upon the broadcast of the settlement transaction. It confirms that a Bitcoin Node has accepted the underlying BTC Transaction, it reduces the balance of the Qredochain Wallet, and unlocks it to enable further transfers

---

## KVAsset


A KVAsset is simply a wrapper around a set of Keys & Values.
KVAsset are mutable and can be updated. However there are fields within each KVAsset which can be made immutable, their keys are listed in the "Immutable" field which is  set of keys which can't be changed in any updates. This rule is enforced as a consensus rule. 


```
message  PBKVAsset {
    PBKVAssetType Type              = 1;
    map<string, bytes> AssetFields  = 2; 
    repeated string Immutable       = 3;
    string Description              = 4; 
}
```

### Example usage:

#### Settlement Fees 

When a user wishes to settle out of the Qredo system by using a settlement transaction, the underlying chain eg. Bitcoin Chain, requires a fee to be paid. As the underlying chain's fee changes the amount a user needs to pay will vary. This fee can be communicated to users within the Qredo system by looking up a value in a KVAsset. The value would represent a fee as Satoshis/Byte, and parties within the system can determine the appropriate fee based on this field. In this case the KVAsset will have TransferRules where a set of trusted users authorise any updates to the Fee amounts


#### Pre-Defined Members of the Qredo Network

The Public Keys for the MPC nodes, the Watcher service, and any other permissioned services on the Qredo Network will have entries in a KVAsset, attached to theses Assets are conditions which will require a number of Signatures before they can be updated.

### Consensus Rules
1. An updated KVAsset can't change/delete any field which is present in the Immutable set (Field #3)
1. An updated KVAsset must contain all the Immutable fields specified in the previous version (new ones can be added)


---

# Crystalization

The values of each  UTXOs once added to the Qredochain do not remain tied to the Wallet they originally fund.  At any point a crystalization process can be performed, mapping UTXOs on an underlying chain to Balances stored in a Qredochain Wallet.
The mapping process is deterministic, so any Qredochain Node given a specific block-height will generate the same set of relationships.


The goal of this process is two fold

1.) The sum of all UTXOs will match the sum of all Qredochain Wallet balances.
All addresses used within the system (specified in MPC transactions) have Proofs that guarantee that the MPC network is able to generate signed transactions to Peg-out all locked up funds. Together these prove that all coins locked into the Qredochain system are spendable by the Qredo System, and therefore the entire system is solvent.

1. A user at any point in time, can ask for a Crystalization, where the funds they own in their wallet will be mapped to a real UTXO on the underlying chain, again when taken together with the MPC proof, and the solvency proof, they can be confident their funds in the Qredochain are backed by real funds on the underlying chain.


When a settlement of the Assets is required, the transactions generated by the crystalization process are sent to the MPC Cluster for signing, (any change is returned to an MPC issued BTC Address, and used in any future crystalization processes.

---

# Example Usage

Here we walk though a typical workflow purely from the Qredochain Transaction point of view. A majority of the communication between parties, such as obtaining signatures from other trustees is handled by the Matrix Communication Protocol, and is not mentioned.


## Funding - Adding fund to your account
1. Alice, Bob & Charlie & Dave create their IDDocuments, each IDDoc is added to the Qredochain.
1. Alice creates a Wallet Asset (Wallet_Alice) of type Bitcoin (BTC), its added to the Qredochain. The Wallet has a transfer rule which requires either Charles or Dave to additional sign any updates to the Wallet (expression: `Charlie | Dave & Alice')
1. The Watcher detects the new Wallet_Alice and requests the MPC to issue a new BTC Address based, the user supplied element of the algorithm to produced the address is the Asset_ID.
1. MPC creates a BTC_Address_A and a Qredochain transaction mapping the BTC Address to Wallet_Alice (MPC_A)
1. BTC_Address_A is added as a watch only address on the watchers BTC Node. From this point forward the BTC Node will monitor the address for changes. 
1. Externally Anne deposits BTC into BTC_Address_A
1. The watcher periodically (every minute) requests wallet changes from its BTC Node. Where new funds have been added with sufficient confirmations, it creates a PEG-IN transaction on the Qredochain. 
1. As this PEG-IN (underlying) transaction is committed to the chain, Qredochain consensus rules uses MPC_A to provide the mapping from BTC_Address_A to Wallet_Alice, and the amount of BTC deposited into the address is added to the Wallet_Alice BTC balance, which is a value in the consensus database with key  "Wallet_Alice.balance"
1. Funding this address can be repeated at will, and simply increases the balance in the database key.


## Spending - Sending funds to other parties
1. Alice arranges to send Bob some BTC
1. Alice creates an update Wallet_A transaction, which transfers some of the BTC to another wallet (Wallet_Bob), the transaction is signed by Alice and either Charlie or Dave, and the signatures aggregated into a single signature. The update is broadcast to the QredoChain.
1. As the Wallet Update is committed to the Chain, the Qredo consensus rules deduct the transferred amount from Wallet_Alice.Balance key and add it to Wallet_Bob.Balance.

(Assets can be transferred between parties at will, time to transact is dependant on the commitment of a single transaction to the Qredochain, which is currently around 1 second. - an additional fee transaction can be added here which will Pay a pre-determined Qredo wallet a minimal fee)


## Settlement - Getting funds out of Qredo
1. Alice constructs an update to Wallet_Alice, the type is settlement, and the required signatures from either Dave or Charlie are obtained, to make it valid. (Note: these signer could be different for different types of transaction)
1. The Update is committed to the Qredochain, locking the Wallet from any further updates
1. The watcher detects a settlement, crystalizes the chain, and sends the resultant transaction to the MPC for signing.
1. The MPC queries any Qredochain node to confirm the existence, destination and amounts of the Wallet_Alice settlement.
1. The MPC creates a new Address to accept the unspent funds/change from the transaction, this is added to the BTC_Node for monitoring, and a MPC transaction is add to the Qredochain, where the Wallet mapping is empty.
1. The MPC signs the transactions and broadcasts them to the underlying blockchain 
1. The MPC creates a PEG-OUT Qredochain transaction finalizing the settlement, which updates the wallet balance and unlocks the wallet, allowing further updates.




# Watcher

The watcher is an untrusted mediator, compromise will result in a denial of service but not any ability to steal funds. The only trust in the system should lie with the Qredochain, and the MPC Cluster.

The watcher service is responsible for ferrying data to and from the Qredochain, external blockchains and the MPC Cluster.
Functions can be broken down into the following 3 processes.

## Wallet Creation
When a user creates a new wallet, the Watcher which is constant monitoring the Qredochain for new transactions picks up the transaction, and send instructs the MPC Node to generate a new Address.
The MPC node needs to verify the validity of the request against either a proof supplied along with the transaction, or by querying a Qredochain node.
Any resultant MPC transaction can be ferried back to the Qredochain via the Watcher, its validity can be check by verifying the signatures against the public keys of the MPC Cluster
The Watcher takes this newly created address and adds it to a watch wallet on the BTC Node.


## Underlying Funding Transaction
When a user funds a BTC Address previous provided by the MPC cluster, the transaction is noted by the BTC Node and when the node is queried periodically by the watcher, it can provide the full UTXO details, together with a Proof.
The watcher can generate a Peg-In (Underlying) transaction to reflect this funding transaction in the Qredochain.
Each Qredochain Node can check the existence of this BTC Transaction by queuing a BTC Node, (with SPV proof supplied?)
Compromise of the Watcher would prevent the funding transaction from being made available in the Qredochain, but would not put funds at risk.


## Settlement
When a user requests a settlement, it is necessary to ensure that the process is stepped so they are unable perform the withdrawal more than once.
A user must commit a valid Wallet update to the chain, which requests the sending of the underlying funds to specfied Bitcoin Address.
The watcher forwards the request to the MPC node.


This is a work in progress....

1. User retrieves the UTXO(s) for their specific wallet's undelying funds using crystalization.
1. The Transaction is commited into the Qredochain, where:
    1. Check if the UTXOs are still unspent.
    1. The UTXOs involved are locked, and can't be re-used.
    1. The Wallet has its balance reduced (it can remain unlocked)
1. The Watcher sends the request to the MPC Cluster
1. The MPC Cluster checks the validty of the request.
1. If valid it signs the transaction, broadcasts it.
1. MPC, via the watcher posts a Peg-Out (Settlement Completion) TX, containing details of this underlying transaction back to the Qredochain.
1. MPC Peg-Out is committed to the Qredochain
1. An Peg-In Underlying Transaction is sent to the Qredochain based on the Change address, we don't need to wait for confirmation as its us that made it.



    







