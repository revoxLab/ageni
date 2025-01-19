package signature

import (
	"crypto/sha256"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	priv = "0xd39612a2328892027bbafb2f40889f56ff872edd24bb755775ac4782c4bd448a"
	pubk = "0x89a04eb167a482659f1b8e0e8c597f69f5b1d4e96652feccf4b01710a45ac6c8"
)

func TestGen(t *testing.T) {
	pub, private := GenKey()
	t.Logf("pubk=%s\n", pub)
	t.Logf("priv=%s\n", private)
}

func TestSignData(t *testing.T) {
	msg := make([]byte, 0)
	msg = append(msg, []byte("any string text")...)
	hash256 := sha256.New()
	hash256.Write(msg)
	eHash := hash256.Sum(nil)
	sig, err := SignData(priv, eHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signature=%s\n", sig)
}

func TestVerify(t *testing.T) {

	signature := "0xd94b769be1490c3a29fb698f7a8115e0c48cbe61085b25082a63ab3659660289b1a2818009e7beb829a27cda3328d43d0f82ed985b716cb9fcbce8f27684d509"
	msg := make([]byte, 0)
	msg = append(msg, []byte("any string text")...)
	hash256 := sha256.New()
	hash256.Write(msg)
	eHash := hash256.Sum(nil)
	result, err := Verify(signature, eHash, pubk)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("result=%t\n", result)
}

func TestAll(t *testing.T) {
	msg := make([]byte, 0)
	msg = append(msg, []byte("any string text")...)
	hash256 := sha256.New()
	hash256.Write(msg)
	eHash := hash256.Sum(nil)
	sig, err := SignData(priv, eHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signature=%s\n", sig)
	result, err := Verify(sig, eHash, pubk)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("result=%t\n", result)
	assert.Equal(t, true, result)
}
