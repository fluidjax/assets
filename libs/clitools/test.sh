assetid=$(qc cid "chris" | jq -r .object.Asset.ID)
txid=$(qc cq $assetid | jq -r .hex)
asset=$(qc tq "tx.hash='$txid'")


echo $asset