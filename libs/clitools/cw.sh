#!/bin/bash

echo Build Binary
build

echo Make IDDocs

principal=$(qc cid -b=true "p")
pAssetID=$(echo $principal | jq -r .assetid)
pSeed=$(echo $principal | jq -r .seed)

trustee1=$(qc cid -b=true "t1")
t1AssetID=$(echo $trustee1 | jq -r .assetid)
t1Seed=$(echo $trustee1 | jq -r .seed)

trustee2=$(qc cid -b=true "t2")
t2AssetID=$(echo $trustee2 | jq -r .assetid)
t2Seed=$(echo $trustee2 | jq -r .seed)

trustee3=$(qc cid -b=true "t3")
t3AssetID=$(echo $trustee3 | jq -r .assetid)
t3Seed=$(echo $trustee3 | jq -r .seed)

trustee4=$(qc cid -b=true "t4")
t4AssetID=$(echo $trustee4 | jq -r .assetid)
t4Seed=$(echo $trustee4 | jq -r .seed)

newOwner=$(qc cid -b=true "t3")
newOwnerAssetID=$(echo $newOwner | jq -r .assetid)
newOwnerSeed=$(echo $newOwner | jq -r .seed)


createJSON='{"ownerseed":"'$pSeed'",
"expression":"t1 + t2 + t3 > 1 & p",
"transferType":1,
"participants":[
	{"name":"p","ID":"'$pAssetID'"},
	{"name":"t1","ID":"'$t1AssetID'"},
	{"name":"t2","ID":"'$t2AssetID'"},
	{"name":"t3","ID":"'$t3AssetID'"}
],
"currency":"BTC"
}'


createJSON=$(echo $createJSON | tr -d '\n')
createJSON=$(echo $createJSON | tr -d '\r')

echo Make Wallet

newWallet=$(qc cw -b=true -j="$createJSON")
serializedSignedAsset=$(echo $newWallet | jq -r .serializedSignedAsset)
newWalletID=$(echo $newWallet | jq -r .assetid)

echo Check Verification
verify=$(qc vtx "$pAssetID" "$serializedSignedAsset")


echo Update Wallet

updateJson='{
	"ExistingWalletAssetID":"'$newWalletID'",
    "Newowner":"'$newOwnerAssetID'",
	"transferType":1,
	"currency":"BTC",
	"Transfer":[{
		"TransferType":1,
		"Expression":"exp",
		"description":"some description goes here",
	    "participants":[
			{"name":"p","ID":"'$newOwnerAssetID'"},
			{"name":"t4","ID":"'$t4'"},
			{"name":"t2","ID":"'$t2'"},
			{"name":"t3","ID":"'$t3'"}
			] }]
	}'
# echo updateJson
# updateJson=$(echo $updateJson | tr -d '\n')
# updateJson=$(echo $updateJson | tr -d '\r')
# updateJson=$(echo $updateJson | tr -s '"' '\"')


updateWallet=$(qc uw -b=false -j="$updateJson")
serializedSignedAsset=$(echo $updateWallet | jq -r .serializedSignedAsset)
newWalletID=$(echo $updateWallet | jq -r .assetid)

echo $updateWallet