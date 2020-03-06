TODO:

IDDOC
  NewIDDoc
  


# QC - Qredochain Command Line Utility

## Overview
QC is based on the 'bitcoin-cli', it is a wrapper for the Asset library, which allows command line access to the QredoChain (Tendermint).

The tool is written in golang and makes use of the github.com/urfave/cl package.
https://github.com/urfave/cli/blob/master/docs/v2/manual.md



## Help

Help is available for each command

```qc help``` -  Return a list of available commands


```qc help ("command")``` - Return help on a specific command.


## Command Line Usage

Simple commands are passed as normal paramaters.
eg.

```qc createidoc "seed" "identifier"```
- Will create a new IDDoc using the supplied seed and identifier
- The command returns a pretty-printer JSON of the created iddoc transaction.


More complex commands use JSON based parameters. 

```qc createwallet '{"expression": "t1 + t2 + t3 > 1 & p","owner": "IDDOC_HASH1", "authref": "authentication refernce","participants":{"t1":"IDDOC_HASH1","t2":"IDDOC_HASH2","t2":"IDDOC_HASH3", "t2":"IDDOC_HASHP"}}' ```

Pretty printed the json looks like this:-
````
{
  "expression": "t1 + t2 + t3 > 1 & p",
  "owner": "IDDOC_HASH1",
  "authref": "authentication refernce",
  "participants": {
    "t1": "IDDOC_HASH1",
    "t2": "IDDOC_HASH2",
    "t2": "IDDOC_HASH3",
    "t2": "IDDOC_HASHP"
    }
}
````

Data returned from all command line commands is in the form of JSON.




