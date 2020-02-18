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




createJSON='{"ownerseed":"3772e3fa880e1912498d2fc48a367a2058c69ea4bf6ec3cf41fbbb6d8089f8868f3c46e31d8e9ab251ea5e4c6f5ded53",
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
qc cw -b=true -j="$createJSON"

