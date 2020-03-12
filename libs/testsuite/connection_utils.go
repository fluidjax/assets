package testsuite

//Utilities to help connect to an instance of the Qredochain - tendermint server and run testss

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/ethereum/go-ethereum/node"
	"github.com/qredo/assets/libs/prettyjson"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var done chan bool
var ready chan bool
var app *qredochain.QredoChain
var tnode *node.Node
var nc *qredochain.NodeConnector
var out <-chan ctypes.ResultEvent
var subClient *tmclient.HTTP

func ShutDown() {
	print("Shutdown")

	if tnode != nil {
		tnode.Stop()
		tnode.Wait()
	}
	if subClient != nil {
		subClient.Unsubscribe(context.Background(), "test", "tx.height>0")

	}
}

func StartWait(count int, wg *sync.WaitGroup) {
	incomingCount := 0
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case _ = <-out:
				incomingCount++
				if incomingCount == count {
					time.Sleep(1000 * time.Millisecond)
					wg.Done()
				}
			}
		}
	}()
}
func SubscriptionClient() {
	var err error
	subClient, err = tmclient.NewHTTP("tcp://localhost:26657", "/websocket")
	if err != nil {
		os.Exit(1)
	}
	err = subClient.Start()
	if err != nil {
		os.Exit(1)
	}

	out, err = subClient.Subscribe(context.Background(), "test", "tx.height>0", 1000)
	if err != nil {
		print("error")
	}
}

func StartTestChain() {
	//Check if there is a Node already running, and use it
	var err error
	nc, err = qredochain.NewNodeConnector("127.0.0.1:26657", "NODEID", nil, nil)
	if err == nil {
		SubscriptionClient()

		defer func() {
		}()
		return
	}
	//If no running Node start one
	go InitiateChain()
	ready = make(chan bool, 1)
	<-ready //wait for server to come up
}

func InitiateChain() {
	//Make Database
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	app = qredochain.NewQredoChain(db)

	flag.Parse()

	//Start Qredochain
	tnode, err := qredochain.NewTendermint(app, "/tmp/example/config/config.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	tnode.Start()

	//Make Connector
	nc, err = qredochain.NewNodeConnector("127.0.0.1:26657", "NODEID", nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
	defer func() {
		nc.Stop()
		db.Close()
		tnode.Stop()
		tnode.Wait()
	}()
	done = make(chan bool, 1)
	ready <- true //notify the server is up
	<-done        //wait

}

func prettyStringFromSignedAsset(signedAsset *protobuffer.PBSignedAsset) string {
	original := reflect.ValueOf(signedAsset)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)

	pp, _ := prettyjson.Marshal(copy.Interface())
	return string(pp)
}
func TranslateRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		TranslateRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		TranslateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			TranslateRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:

		if original.Type() == reflect.TypeOf([]byte("")) {

			b := original.Bytes()
			h := hex.EncodeToString(b)
			hb := []byte(h)
			copy.Set(reflect.ValueOf(hb))

		} else {
			copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
			for i := 0; i < original.Len(); i += 1 {
				TranslateRecursive(copy.Index(i), original.Index(i))
			}
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			TranslateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		b64 := base64.StdEncoding.EncodeToString([]byte(original.String()))
		copy.Set(reflect.ValueOf(b64))
		//translatedString := dict[original.Interface().(string)]
		//		copy.SetString(original.String())

	// And everything else will simply be taken from the original
	default:

		copy.Set(original)
	}

}

//Dump - Retrieve from the chain and pretty printzlib
func dump(key []byte) {

	//Pause to wait for tx into block
	time.Sleep(2 * time.Second)
	sa, err := nc.GetAsset(hex.EncodeToString(key))
	if err != nil {
		panic("Error dumping")
	}

	fmt.Printf("%v", prettyStringFromSignedAsset(sa))
}
