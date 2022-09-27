package src

//go:generate abigen --sol ../contracts/CallerTest.sol --pkg wrappers --out ../wrappers/RemoteCallerTest.go
//go:generate abigen --sol ../contracts/ReceiverTest.sol --pkg wrappers --out ../wrappers/RemoteCallRececiverTest.go
//go:generate abigen --sol ../contracts/BridgeMock.sol --pkg wrappers --out ../wrappers/BridgeMock.go
