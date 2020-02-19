#!/bin/bash

dump (){
	echo $2
	echo -e "\033[40;31;5;82m $1  \033[0m"	
	res=$(echo $1 | sed 's/"/\\"/g')
	echo -e "\033[40;38;5;82m $res  \033[0m"
}



echo Build Binary
build

MakeIDDocs() {
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
}



MakeWallet() {
	echo Make Wallet

	createJSON='
	{
	"transferType": 0,
	"ownerseed":"'$pSeed'",
	"currency": "BTC",
	"Transfer": [{
		"TransferType": 2,
		"Expression": "t1 + t2 + t3 > 1 & p",
		"description": "Here is the transfer Type 2",
		"participants": [
		{
			"name": "p",
			"ID": "'$pAssetID'"
		},
		{
			"name": "t1",
			"ID": "'$t1AssetID'"
		},
		{
			"name": "t2",
			"ID": "'$t2AssetID'"
		},
		{
			"name": "t3",
			"ID": "'$t3AssetID'"
		}
		]
	}]
	}'

	createJSON=$(echo $createJSON | tr -d '\n')
	createJSON=$(echo $createJSON | tr -d '\r')

	dump "$createJSON" CreateWallet 
	newWallet=$(qc cw -b=true -j="$createJSON")
	serializedSignedAsset=$(echo $newWallet | jq -r .serializedSignedAsset)
	newWalletID=$(echo $newWallet | jq -r .assetid)
}

VerifyWallet() {
	echo Check Verification
	verify=$(qc vtx "$pAssetID" "$serializedSignedAsset")
}


SerializeUpdate() {
	echo Build Serialized Updated

	updateJson='{
		"ExistingWalletAssetID":"'$newWalletID'",
		"Newowner":"'$newOwnerAssetID'",
		"transferType":2,
		"currency":"BTC",
		"Transfer":[{
			"TransferType":1,
			"Expression":"different expression",
			"description":"some description goes here",
			"participants":[
				{"name":"p","ID":"'$newOwnerAssetID'"},
				{"name":"t1","ID":"'$t1AssetID'"},
				{"name":"t2","ID":"'$t2AssetID'"},
				{"name":"t4","ID":"'$t4AssetID'"}
				] }]
		}'


	updateJson=$(echo $updateJson | sed 's/\n//g')
	updateWallet=$(qc uw -j="$updateJson")
	serializedUnsignedAsset=$(echo $updateWallet | jq -r .serializedUpdate)
}
#----------------------------------------------------------------------------------------------------------------------------------------

SignForEachIDoc() {
	echo "Sign for each IDDoc"

	#sign for each IDDoc
	json='{"seed":"'$pSeed'","msg":"'$serializedUnsignedAsset'"}'
	sigP=$(qc sign -j="$json" | jq -r .signature)

	json='{"seed":"'$t1Seed'",	"msg":"'$serializedUnsignedAsset'"}'
	sigT1=$(qc sign -j="$json" | jq -r .signature)

	json='{	"seed":"'$t2Seed'",	"msg":"'$serializedUnsignedAsset'"}'
	sigT2=$(qc sign -j="$json" | jq -r .signature)

	json='{	"seed":"'$t3Seed'",	"msg":"'$serializedUnsignedAsset'"}'
	sigT3=$(qc sign -j="$json" | jq -r .signature)
}



#----------------------------------------------------------------------------------------------------------------------------------------

AggregateSign() {
	#Aggregate Sign
	echo Aggregate Sign

	json='{"Sigs":[
				{"id":"'$pAssetID'","abbreviation":"p","signature":"'$sigP'"},
				{"id":"'$t1AssetID'","abbreviation":"t1","signature":"'$sigT1'"},
				{"id":"'$t2AssetID'","abbreviation":"t2","signature":"'$sigT2'"},
				{"id":"'$t3AssetID'","abbreviation":"t3","signature":"'$sigT3'"}
				],
			"walletUpdate":'$updateJson'
		}'

	dump "$json" "Aggregate Sign"

	updateComplete=$(qc as -b=true -j="$json")

}



MakeIDDocs
MakeWallet
VerifyWallet
SerializeUpdate
SignForEachIDoc
AggregateSign