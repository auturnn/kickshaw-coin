package utils

import "errors"

//P2P
var ErrNetworkIsNotWork = errors.New("p2p network is not working")

//DB
var ErrCreateDB = errors.New("database not created")
var ErrSaveBlock = errors.New("database failed to modify the chain")
var ErrSaveChain = errors.New("database failed to save chain")

//Blockchain
var ErrLoadDB = errors.New("blockchain is not loaded")
var ErrCreateBlockChain = errors.New("blockchain is not created")
var ErrTargetBlockChainNotFound = errors.New("target blockchain is not found")
var ErrReplaceBlockchain = errors.New("blockchain is not replaced")
var ErrTargetBlockNotFound = errors.New("target block is not found")

//server function error
var ErrPeerNotConnect = errors.New("failed connect to peer")

var ErrLogPath = errors.New("failed create to log path ")
