package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/fs"
	"math/big"
	"os"

	"github.com/auturnn/kickshaw-coin/utils"
)

const (
	walletName string = "kickshaw.wallet"
)

type fileLayer interface {
	hasWalletFile() bool
	writeFile(name string, data []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}

type layer struct{}

func (layer) hasWalletFile() bool {
	_, err := os.Stat(walletName)
	return !os.IsNotExist(err)
}

func (layer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (layer) readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

var files fileLayer = layer{}

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func decodeString(payload string) []byte {
	bytes, err := hex.DecodeString(payload)
	utils.HandleError(err)
	return bytes
}

func Sign(payload string, w *wallet) string {
	payloadBytes := decodeString(payload)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadBytes)
	utils.HandleError(err)

	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}

	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]

	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)

	return &bigA, &bigB, nil
}

func Verify(sign, payload, addr string) bool {
	r, s, err := restoreBigInts(sign)
	utils.HandleError(err)

	x, y, err := restoreBigInts(addr)
	utils.HandleError(err)

	//not used privateKey
	puK := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	payloadBytes := decodeString(payload)
	return ecdsa.Verify(&puK, payloadBytes, r, s)
}

func createPrivKey() *ecdsa.PrivateKey {
	prk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return prk
}

func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := files.readFile(walletName)
	utils.HandleError(err)

	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleError(err)

	return key
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)
	utils.HandleError(files.writeFile(walletName, bytes, 0644))
}

func addrFromKey(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if files.hasWalletFile() {
			// true = restore from file
			w.privateKey = restoreKey()
		} else {
			//has a wallet already
			// false = create prk , save to file
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = addrFromKey(w.privateKey)
	}
	return w
}

// func Start() {
// 	prkBytes, err := hex.DecodeString(prk) // 16진수 인코딩이 맞는지 확인
// 	utils.HandleError(err)

// 	restoredKey, err := x509.ParseECPrivateKey(prkBytes)
// 	utils.HandleError(err)

// 	fmt.Println(restoredKey)

// 	signBytes, err := hex.DecodeString(sign)
// 	rBytes := signBytes[:len(signBytes)/2]
// 	sBytes := signBytes[len(signBytes)/2:]

// 	var bigR, bigS = big.Int{}, big.Int{}
// 	bigR.SetBytes(rBytes)
// 	bigS.SetBytes(sBytes)

// 	hmsgBytes, err := hex.DecodeString(hmsg)
// 	utils.HandleError(err)

// 	ok := ecdsa.Verify(&restoredKey.PublicKey, hmsgBytes, &bigR, &bigS)
// 	fmt.Println(ok)
// }

// const msg string = "ploy_kickshaw"

// func Start() {
// 	//ecdsa => Elliptic Curve Digital Signature Algorithem
// 	prk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	utils.HandleError(err)

// 	hmsg := utils.Hash(msg)
// 	fmt.Println(hmsg)

// 	keyAsBytes, err := x509.MarshalECPrivateKey(prk) //parse the key
// 	utils.HandleError(err)
// 	fmt.Printf("prk : %x\n", keyAsBytes)

// 	hAsBytes, err := hex.DecodeString(hmsg)
// 	utils.HandleError(err)

// 	r, s, err := ecdsa.Sign(rand.Reader, prk, hAsBytes)
// 	utils.HandleError(err)

// 	sign := append(r.Bytes(), s.Bytes()...)
// 	fmt.Printf("sign: %x\n", sign)

// 	ok := ecdsa.Verify(&prk.PublicKey, hAsBytes, r, s)
// 	fmt.Println(ok)
// }
