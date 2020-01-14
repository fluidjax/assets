
# Assets

An asset object is a multi function representation of a piece of data.

Each asset has a unique  Asset ID, which remains the same as the asset is updated.

Only valid Assets and their Asset Transfers are written to the blockchain, so the most recent version on chain can be deemed correct and current.

TransferRules detailed in each Asset define if and how a new version of the Asset is to be permitted.


# Mutable/Immutable Assets

## Immutable Asset
An immutable Asset such as an IDDocument, once the declaration of the Asset has been made.
It is fixed and can't be changed. The immutability is determined by the Assets lack of a TransferRule.
No rule to transfer to a new version means it is not possible to create a valid updated version, rendering the Asset immutable.


## Mutable Asset
A mutable Asset, such as a Wallet or TrusteeGroup
Once created the Payload of such as Asset can be changed.
The change is only valid if the new version adheres to the transfer rules previous declared.
As an example a wallets owner can be transferred to a new owner if the original owner plus 2 of 3
trustees declared in the original Asset sign the new version of the Asset.

The transfer expression in this case may look like this 

```
"(t1 + t2 + t3 > 1) & P"
```

Indicating that if 2 of the three trustees (t1, t2, t3)  together with the Principal (p) sign a transferred Wallet Asset
it will be valid.

    
An asset can be analysed using the follow functions: 
1) TruthTable - this function analyses the Asset and returns a list of all the valid combinations of signatures which will permit the transfer
in the example above this would be :-

```
        "(t1 + t2 + 0 > 1) & P"
        "(t1 + 0 + t3 > 1) & P"
        "(0 + t2 + t3 > 1) & P"
```

2) ParseChange - this function depends on the Asset Type
* IDDoc  - not valid, this is an immutable type
* Wallet - unimplemented
* Group  - a breakdown of how the asset has change between versions, detailing new, delete and unchanged fields (IDs)


# Transfer Overview

Transfers from old to new version are controlled using TransferRules.

Each Asset can have any number of TransferRules, but each Rule type must only exist once.
This permits different rules depending on the type of transfer (PBTransferType).

For Example a wallet could have 2 different TransferRules, one to enable the transfer to 
other partties, but more stringent rules when it comes to Settling (withdrawing the contents)
the wallet on the underlying chain.

# API

An [asset] is currently either an IDDoc, Wallet or Group 


````
    [asset].Payload              - retrieve the Payload object 
    [asset].Verify               - Verify the signed PBAsset (including Payload)
    [asset].Sign                 - Sign the PBAsset (including Payload)
    [asset].Save                 - save asset to store
    [asset].Dump                 - pretty print asset
    [asset].Key()                - get key of asset
    [asset].SerializeAsset       - Serialize the Asset for transmission/signing
    [asset].SerializeSignedAsset - Serialize the Signed Asset for storage
     

    Sign(data, IDDoc)            - Sign the complete Asset  with IDDoc's keys
    Verify(data, sig, IDDoc)     - Verify the complete Asset signature with IDDoc's keys
    Load(store, key)             - retrieve SignedAsset from store
    New[asset]                   - create a new Object, pre-populate essential items
    ReBuild[asset]               - Recreate a New[asset] based on a SignedAssett from the Store
    NewUpdate[asset]             - Recreate a New[asset], based on and pointing too a previous Asset.
````

````
    Transfers
    [asset].AddTransfer          - Add a new transfer rule to an asset
    [asset].IsValidTransfer      - Calculate using sigs if the expression has been fullfilled
    [asset].TruthTable           - a truth table to all sucessful permutation of signatures
    [asset].AggregatedSign       - Aggregate the supplied Signatures/PubKeys and add to Asset
    [asset].FullVerify           - Fully validate if a transfer is correct and allowed on chain
    [asset].ParseChanges         - breakdown the changes between old & new version of Asset
````