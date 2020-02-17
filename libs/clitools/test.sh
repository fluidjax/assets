seed=$(qc cid "chris" | jq -r .seed)
txid=$(qc cw  $seed)
#asset=$(qc tq "tx.hash='$txid'")


#echo $asset
