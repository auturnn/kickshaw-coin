package wallet

import (
	"crypto/x509"
	"io/fs"
)

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) readFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}

func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}
