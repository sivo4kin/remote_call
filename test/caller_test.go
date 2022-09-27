package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sivo4kin/remote_call/wrappers"
	"github.com/stretchr/testify/require"
	"math/big"
	"sync"
	"testing"
	"time"
)

var (
	tx *types.Transaction
)

func Test_RemoteCaller(t *testing.T) {
	mx := new(sync.Mutex)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	eventChan := make(chan *wrappers.BridgeReceiveRequestWithReturnValue)

	eventsMap := make(map[*wrappers.BridgeReceiveRequestWithReturnValue]struct{}, 0)


	go getEventChanFromBridge(ctx, eventChan)



	t.Run("Run_sendReceiveRequestWithVerifyCallResultTx", func(t *testing.T) {
		tx, err := sendFailToSetTestValueCalltTx(t)
		if err == nil {
			mx.Lock()
			t.Log(tx.Hash())
			mx.Unlock()
		}

		backend.Commit()
	})

	t.Run("Run_sendFailToSetTestValueCalltTx", func(t *testing.T) {
		tx, err := sendFailToSetTestValueCalltTx(t)
		if err == nil {
			mx.Lock()
			t.Log(tx.Hash())
			mx.Unlock()
		}

		backend.Commit()
	})

	t.Run("Run_sendRemoteCalltTx", func(t *testing.T) {
		tx, err := sendRemoteCalltTx(t)
		if err == nil {
			mx.Lock()
			t.Log(tx.Hash())
			mx.Unlock()
		}
		backend.Commit()
	})

	t.Run("Run_sendAllwaysSuccessTestTx", func(t *testing.T) {
		tx, err := sendAllwaysSuccessTestTx(t)
		if err == nil {
			mx.Lock()
			t.Log(tx.Hash())
			mx.Unlock()
		}
		backend.Commit()
	})
	t.Run("Run_sendTestTx", func(t *testing.T) {
		tx, err := sendFailingTx(t)
		if err == nil {
			mx.Lock()
			t.Log(tx.Hash())
			mx.Unlock()
		}
		backend.Commit()
	})

loop:
	for {
		select {
		case event, ok := <-eventChan:
			if ok {

				//fmt.Printf("SENT DATA: %s\n", common.Bytes2Hex(event.SentData))
				//fmt.Printf("Returndata DATA: %s\n", common.Bytes2Hex(event.Returndata))
				fmt.Printf("Returndata DATA String:%s\n", event.Returndata)
				//fmt.Printf("Returndata DATA Bytes: %s\n", common.Bytes2Hex(event.RetBytes))
				fmt.Printf("SUCCESS: %v\n", event.Success)
				//fmt.Printf("RAW: %v\n", event.Raw)

			/*	err_msg, err := unpackError(event.Returndata)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Contract returned message: %v\n", err_msg)*/

				eventsMap[event] = struct{}{}
				if len(eventsMap) == 5 {
					break loop
				} else {
				}
			}


		case <-ctx.Done():
			break loop
		}

	}
	require.NoError(t, ctx.Err())
	for event2 := range eventsMap {
		if !event2.Success {
			t.Logf("FAILED Tx  [%v]", event2.Raw.TxHash)
		} else {
			t.Logf("Succeded Tx %v  !!!", event2.Raw.TxHash)
		}
	}

}

func sendFailingTx(t *testing.T) (*types.Transaction, error) {
	testdata := []byte("WRONGDATASENT")
	return bridge.ReceiveRequestWithReturnValue(owner,
		testdata,
		recieverAddress)
}

func sendAllwaysSuccessTestTx(t *testing.T) (*types.Transaction, error) {
	testdata := []byte("SUCCESSMESSAGE")
	return bridge.ReceiveRequestAllwaysSuccess(owner,
		testdata,
		recieverAddress)
}

func sendRemoteCalltTx(t *testing.T) (*types.Transaction, error) {
	createNodeABIPacked, err := recieverABI.Pack("setTestValue", big.NewInt(123))
	if err != nil {
		fmt.Printf("sendRemoteCalltTx %v", err)
		return nil, err
	}

	return bridge.ReceiveRequestWithReturnValue(owner,
		createNodeABIPacked,
		recieverAddress)
}



func sendReceiveRequestWithVerifyCallResultTx(t *testing.T) (*types.Transaction, error) {
	createNodeABIPacked, err := recieverABI.Pack("failToSetTestValue", big.NewInt(123))
	if err != nil {

		fmt.Printf("failToSetTestValue %v\n", err)
		return nil, err
	}

	return bridge.ReceiveRequestWithVerifyCallResult(owner,
		createNodeABIPacked,
		recieverAddress)
}


func sendFailToSetTestValueCalltTx(t *testing.T) (*types.Transaction, error) {
	createNodeABIPacked, err := recieverABI.Pack("failToSetTestValue", big.NewInt(123))
	if err != nil {
		fmt.Printf("sendRemoteCalltTx %v", err)
		return nil, err
	}

	return bridge.ReceiveRequestWithReturnValue2(owner,
		createNodeABIPacked,
		recieverAddress)
}

func getEventChanFromBridge(ctx context.Context, eventChan chan<- *wrappers.BridgeReceiveRequestWithReturnValue) {
	recieveRequestEvent = make(chan *wrappers.BridgeReceiveRequestWithReturnValue)
	defer close(recieveRequestEvent)

	bridgeInstance, err = wrappers.NewBridge(bridgeAddress, backend)
	if err != nil {
		fmt.Printf("NewBridge %v", err)
		return
	}
	sub, err := bridgeInstance.WatchReceiveRequestWithReturnValue(&bind.WatchOpts{Context: ctx}, recieveRequestEvent)
	if err != nil {
		fmt.Printf("WatchOracleRequest %v", err)
		return
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			fmt.Printf("sub %v", err)
			return
		case received, ok := <-recieveRequestEvent:
			if ok {
				eventChan <- received

			}
		case <-ctx.Done():
			return
		}
	}

}


func unpackError(result []byte) (string, error) {
	var (
		errorSig            = []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
		abiString, _        = abi.NewType("string", "", nil)
	)
	if !bytes.Equal(result[:4], errorSig) {
		return "<tx result not Error(string)>", errors.New("TX result not of type Error(string)")
	}
	vs, err := abi.Arguments{{Type: abiString}}.UnpackValues(result[4:])
	if err != nil {
		return "<invalid tx result>", err
	}
	return vs[0].(string), nil
}
