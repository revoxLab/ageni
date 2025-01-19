package signature

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestSign(t *testing.T) {
	fmt.Printf("prvk:%s\n", prvtk)
	privateKey, err := crypto.HexToECDSA(prvtk)
	if err != nil {
		t.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	data := []byte(source)
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("signature0:%s\n", hex.EncodeToString(signature)) // 0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301
	fmt.Println("signature1:" + hexutil.Encode(signature))       // 0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		t.Fatal(err)
	}

	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println(matches) // true

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	fmt.Printf("recover addr=%s\n", address.String())
	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	matches = bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println(matches) // true

	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println(verified) // true
}

const (
	prvtk = "18ab35ec2701be95b824011790b83e70b7c087b5666e2a6fcd423d4b364edda4"
	//source       = "readon.me"
	source       = "abc123"
	addr         = "0x2e19aA55985f2E43212eb644e3d04793e90a77E7"
	badSignature = "0x3D9A74E7F5685F7170BF10E1C9966DF1D1D16835F9DAE95BD82206BF8A8884587508E16D8D971C717485450E8FA0343431646DC56CCF31933EED28DCCE74AB881B"
	signature    = "0x3d9a74e7f5685f7170bf10e1c9966df1d1d16835f9dae95bd82206bf8a8884587508e16d8d971c717485450e8fa0343431646dc56ccf31933eed28dcce74ab8800"
)

func TestPersonalSign(t *testing.T) {
	sig, err := PersonalSign(prvtk, source)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signature=%s\n", sig)
	re, er := ValidateSignature(source, sig, addr)
	if er != nil {
		t.Fatal(er)
	}
	t.Logf("re=%t\n", re)
}

func TestValidateSignature3(t *testing.T) {

	_, err := ValidateSignature(source, badSignature, addr)

	assert.NotNil(t, err)

	b, err := ValidateSignature(source, signature, addr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, b)
}

func TestValidateSignature(t *testing.T) {

	_, err := ValidateSignature(source, badSignature, addr)

	assert.NotNil(t, err)

	b, err := ValidateSignature(source, signature, addr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, b)
}

func TestValidateSignature2(t *testing.T) {
	what, err := VerifyMessage("", "0x64581cce5414f924fd26ae1eaaf96f65b3723dea2e2cb8b7c4e586617a7e20dc3d4e7818848b9267136407b9979cbdefd793a4b0742115e93ed98f2e1f20cd4501", "0xc694C87314E84e9a6FeAe5932196A4C53187e7D6")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("what=%t\n", what)
}

func TestEthereumSign(t *testing.T) {

	sig, err := VoteContentSign(prvtk, source)
	//sig, err := PersonalSign(prvtk, source)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("address=%s\n", addr)
	t.Logf("source=%s\n", source)
	t.Logf("signature=%s\n", sig)
	//re, er := ValidateSignature(source, sig, addr)
	//if er != nil {
	//	t.Fatal(er)
	//}
	//t.Logf("re=%t\n", re)

}

func getRSVFromSignature(signature []byte) (*big.Int, *big.Int, byte, error) {
	r := new(big.Int)
	s := new(big.Int)
	v := byte(signature[64]) + 27

	r.SetBytes(signature[:32])
	s.SetBytes(signature[32:64])

	return r, s, v, nil
}

func TestRSV(t *testing.T) {
	sigBytes, err := hexutil.Decode("0x9dc0b7aab209d7d7d82f03686e9ae3bf33c68f60220013037dbd6cf44a10fcaf69cf04cbecb200b791218495ec1688181321f28235d55960e12ea1f89b3861be00")
	if err != nil {
		t.Fatal(err)
	}
	r, s, v, err := getRSVFromSignature(sigBytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("r=%v,s=%v,v=%d\n", r.Bytes(), s.Bytes(), v)
}
