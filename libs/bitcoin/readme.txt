Bitcoin Underlying Chain Notes:


Testnet
/Applications/Electrum.app/Contents/MacOS/Electrum --testnet

Explorer:
https://live.blockcypher.com/btc-testnet/

Using 
alias b='bitcoin-cli'



ssh btctest
Locally  Electrum 

Test Addresses
n3tnmv9QxrPPJ68Q5gm4V4XCSHtQoiC5Yb
3.3mbtc sent in block   1665007


Block Hash 1665007
000000000000005d55bdf31eb16890fa15db71a6145068bde2c6688b930ee06d

Previous Hash 1665006
00000000000000986bcbdaef170a3c6ec3ba952ad603e9d80116b0068b4aeac2


Watch Address
b importaddress n3tnmv9QxrPPJ68Q5gm4V4XCSHtQoiC5Yb "" false false


b listreceivedbyaddress 6 true true   //confirmations, include empty, include watch only



    //Get all the transactions since this blockhash (ie.set this as the last time we did this check)
    //b listsinceblock blockhash confirmations watchonly removed
    //
b listsinceblock 00000000000000986bcbdaef170a3c6ec3ba952ad603e9d80116b0068b4aeac2 1 true true





//Get all the UTXOs between 6-999999 confirmations
b listunspent 6 9999999 "[\"n3tnmv9QxrPPJ68Q5gm4V4XCSHtQoiC5Yb\"]"



//Get current height
b getblockcount

//Get the Block Hash for height
b getblockhash 1665007


Solution
--------





1) Issue
When we issue a new address add it to the 'track list' with 
b importaddress n3tnmv9QxrPPJ68Q5gm4V4XCSHtQoiC5Yb "" false false

2) Watcher periodically query the chain
Scan chain fro last scanned height to current height - 6

eg. 10 Blocks ago to now (ignore)


Take UTXOs between last scan Height  & 6 blocks old




lastScanBlockHeight=0
hash = b getblockhash $lastScanBlockHeight
b listsinceblock 00000000000000986bcbdaef170a3c6ec3ba952ad603e9d80116b0068b4aeac2 1 true true


Every 60 seconds
----------------
height='b getblockcount'

1665010 - 6 

ubuntu@ip-172-31-11-118:~$ b getblockhash 1665004
0000000000000122b953129fe74aec8d0657e429be9304baf111f5504144d8c2
ubuntu@ip-172-31-11-118:~$ b getblockhash 1665005
000000000000042bec31aeb537de9d68b00557e8c2660e4f39be4499c9def15a
ubuntu@ip-172-31-11-118:~$ b getblockhash 1665006
00000000000000986bcbdaef170a3c6ec3ba952ad603e9d80116b0068b4aeac2
ubuntu@ip-172-31-11-118:~$ b getblockhash 1665007
000000000000005d55bdf31eb16890fa15db71a6145068bde2c6688b930ee06d
ubuntu@ip-172-31-11-118:~$ b getblockhash 1665008
000000000000011669702f82d51d6f8a8abf4f9146a4ee737381f67dd917b437
ubuntu@ip-172-31-11-118:~$ b getblockhash 1665009
00000000000000a0b293d81598cec7c17cfed742183637f71b02e5d238c7c8c6
ubuntu@ip-172-31-11-118:~$ b getblockhash 1665010
0000000000000049020ab2359aecfe780e85d083038c79031470c224980c7431

b listsinceblock 000000000000042bec31aeb537de9d68b00557e8c2660e4f39be4499c9def15a 1 true true

bitcoin-cli rescanblockchain 100000 120000










