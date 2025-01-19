package signature

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/readonme/open-studio/common/log"
)

func GenKey() (string, string) {
	pub, private, _ := ed25519.GenerateKey(rand.Reader)
	return hexutil.Encode(pub), hexutil.Encode(private.Seed())
}

func Sum(e interface{}) []byte {
	buf, err := json.Marshal(e)
	if err != nil {
		log.Info("json marshal error=%v\n", err)
		return nil
	}

	hash256 := sha256.New()
	hash256.Write(buf)
	return hash256.Sum(nil)
}

func SignData(privateKeyHex string, dataHash []byte) (string, error) {
	privateKeyBytes, err := hexutil.Decode(privateKeyHex)
	if err != nil {
		return "", err
	}
	pk := ed25519.NewKeyFromSeed(privateKeyBytes)
	log.Info("sig src=%x\n", dataHash)
	signature := ed25519.Sign(pk, dataHash)
	return hexutil.Encode(signature), nil
}

func Verify(signatureHex string, srcHash []byte, pubHex string) (bool, error) {
	pubk, err := hexutil.Decode(pubHex)
	if err != nil {
		return false, err
	}
	msg, err := hexutil.Decode(signatureHex)
	if err != nil {
		return false, err
	}
	log.Info("src hash=%x\n", srcHash)
	return ed25519.Verify(pubk, srcHash, msg), nil
}
