package test

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sivo4kin/remote_call/wrappers"
)

var (
	backend                                                                              *backends.SimulatedBackend
	owner                                                                                *bind.TransactOpts
	bridge, bridgeInstance                                                               *wrappers.Bridge
	caller                                                                               *wrappers.CallerTest
	reciever                                                                             *wrappers.ReceiverTest
	err                                                                                  error
	ownerKey                                                                             *ecdsa.PrivateKey
	ownerAddress, blsSignatureTestAddress, callerAddress, recieverAddress, bridgeAddress common.Address
	recieveRequestEvent                                                                  chan *wrappers.BridgeReceiveRequestWithReturnValue
	bridgeABI                                                                            abi.ABI
	recieverABI                                                                          abi.ABI
)

func MustGetABI(json string) abi.ABI {
	abi, err := abi.JSON(strings.NewReader(json))
	if err != nil {
		panic("could not parse ABI: " + err.Error())
	}
	return abi
}

func init() {

	bridgeABI = MustGetABI(wrappers.BridgeABI)
	recieverABI =MustGetABI(wrappers.ReceiverTestABI)
	ownerKey, _ = crypto.GenerateKey()

	ownerAddress = crypto.PubkeyToAddress(ownerKey.PublicKey)

	genesis := core.GenesisAlloc{
		ownerAddress: {Balance: new(big.Int).SetInt64(math.MaxInt64)},
	}
	backend = backends.NewSimulatedBackend(genesis, math.MaxInt64)

	owner, err = bind.NewKeyedTransactorWithChainID(ownerKey, big.NewInt(1337))
	if err != nil {
		panic(err)
	}


	callerAddress, _, caller, err = wrappers.DeployCallerTest(owner, backend)
	if err != nil {
		panic(err)
	}

	recieverAddress, _, reciever, err = wrappers.DeployReceiverTest(owner, backend)
	if err != nil {
		panic(err)
	}

	bridgeAddress, _, bridge, err = wrappers.DeployBridge(owner, backend)
	if err != nil {
		panic(err)
	}

	backend.Commit()
}
