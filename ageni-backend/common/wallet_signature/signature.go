package signature

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func ValidateSignature(source, signatureHex, addr string) (bool, error) {
	sourceHash := crypto.Keccak256Hash([]byte(source))
	fmt.Printf("source hash=%s\n", sourceHash.String())
	sigHash, err := hexutil.Decode(signatureHex)
	if err != nil {
		return false, err
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(sourceHash.Bytes(), sigHash)
	if err != nil {
		return false, err
	}
	address := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	return strings.EqualFold(strings.ToLower(addr), strings.ToLower(address.String())), nil
}

func VerifyMessage(source, signatureHex, addr string) (bool, error) {

	address, err := verifyMessage(source, signatureHex)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(strings.ToLower(addr), strings.ToLower(address)), nil
}

func verifyMessage(message string, signedMessage string) (string, error) {
	// Hash the unsigned message using EIP-191
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
	hash := crypto.Keccak256Hash(hashedMessage)

	// Get the bytes of the signed message
	decodedMessage := hexutil.MustDecode(signedMessage)

	// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}

	// Recover a public key from the signed message
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
	if sigPublicKeyECDSA == nil {
		err = errors.New("could not get a public key from the message signature")
	}
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
}

func PersonalSign(private, source string) (string, error) {
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return "", err
	}
	//address := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	//data := []byte(source)
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(source)) + source)

	hash := crypto.Keccak256Hash(hashedMessage)
	fmt.Printf("real hash=%s\n", hash.String())

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signature), nil
}

func VoteContentSign(private, source string) (string, error) {
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return "", err
	}

	sourceBytes := []byte(source)
	srcHash := crypto.Keccak256Hash(sourceBytes)

	prefix := "\x19Ethereum Signed Message:\n32"
	data := append([]byte(prefix), srcHash.Bytes()...)
	hashBytes := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hashBytes.Bytes(), privateKey)
	if err != nil {
		return "", err
	}
	signature[64] += 27

	return hexutil.Encode(signature), nil
}
