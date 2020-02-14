package qc

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/qredo/assets/libs/clitools/lib/prettyjson"
)

var (
	res = make(map[string]interface{})
)

func getEnv(name, defaultValue string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return v
}

//Use - helper to remove warnings
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func hex2base64(h string) string {
	bin, _ := hex.DecodeString(h)
	return base64.StdEncoding.EncodeToString(bin)
}

func addResultItem(key string, item interface{}) {
	res[key] = item
}

func ppResult() {
	pp, _ := prettyjson.Marshal(res)
	fmt.Println(string(pp))
}

//PadRight - right pad a string
func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}
