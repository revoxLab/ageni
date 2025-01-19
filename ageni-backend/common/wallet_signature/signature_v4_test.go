package signature

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
)

var data = `0x6ddd28cd8b8d02aab7d6ec0f464f5f72508bddef1d2dbbd3ab41e5cda700b8ca3966fbd50ee000575493e20a9b23304993724e97eb75a454305fc0bfe784eb0f1c`

//var typedData = `{"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"collectWithSig":[{"name":"collector","type":"address"},{"name":"profileId","type":"uint256"},{"name":"essenceId","type":"uint256"},{"name":"data","type":"bytes"},{"name":"postDatas","type":"bytes"},{"name":"nonce","type":"uint256"},{"name":"deadline","type":"uint256"}]},"primaryType":"collectWithSig","domain":{"name":"Link3","version":"1","chainId":"0x61","verifyingContract":"0x57e12b7a5F38A7F9c23eBD0400e6E53F2a45F271"},"message":{"collector":"0x0FA438Df14B40AEBECb57ca5DAd998385D05F8c5","data":"0x","deadline":"0x1b09da880","essenceId":109,"nonce":"0x0","postDatas":"0x","profileId":109}}`
var typedData = `{"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"registerEssenceWithSig":[{"name":"profileId","type":"uint256"},{"name":"name","type":"string"},{"name":"symbol","type":"string"},{"name":"essenceTokenURI","type":"string"},{"name":"essenceMw","type":"address"},{"name":"transferable","type":"bool"},{"name":"initData","type":"bytes"},{"name":"nonce","type":"uint256"},{"name":"deadline","type":"uint256"}]},"primaryType":"registerEssenceWithSig","domain":{"name":"Link3","version":"1","chainId":"0x61","verifyingContract":"0x57e12b7a5F38A7F9c23eBD0400e6E53F2a45F271"},"message":{"deadline":"0x1b09da880","essenceMw":"0x7FD80D2c47eD1f204851f2809f54f5A31E4d55a3","essenceTokenURI":"https://readon-api.readon.me/v1/ctt/96007","initData":"0x0000000000000000000000000000000000000000000000000000000001406f400000000000000000000000000000000000000000000000001bc16d674ec80000000000000000000000000000b5f7f430a66a845af7d02668d5708e3f603dcc34000000000000000000000000b9ef4d19a095eeef61c1bfb56ef6f76490cadbe20000000000000000000000000000000000000000000000000000000000000000","name":"ReadON Super Bullet NFT","nonce":"0x6","profileId":109,"symbol":"RSB","transferable":true}}`

func TestSignV4(t *testing.T) {
	sig, err := SignV4(typedData, `b1dcfeaf7f55d8dde28db453c47863e304a6f6a79680bbe198f9a6354555608c`)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("sig=%s\n", sig)
}

func TestVerifyV4(t *testing.T) {
	b64 := base64.StdEncoding.EncodeToString([]byte(typedData))
	at := AuthToken{
		TypedData: b64,
		//Signature: `0x9e6e0daa668689bf03546f9c1db2030e6bf5ac8f721de7925c0abf9b0640df643d89b310fd91ff07ebad4df6fcf10eca0fbe7fb6458e0c36d04470a2cf8b492c1c`,
		Signature: `0xdb84a10cc9182b1424c2d2c17dace3a32aa9f2429ce3498dba942898e4d2719f54954caaf6164f4e109dbb23d122a6c7aa10eae762a4df23f0b2c12d2296347b00`,
		//Address:   "0x0FA438Df14B40AEBECb57ca5DAd998385D05F8c5",
		Address: "0x5F5Df3C0AD7579D91127BDC2d112a86273d28fC3",
	}
	data, err := json.Marshal(at)
	if err != nil {
		t.Fatal(err)
	}
	address, err := VerifyAuthTokenAddress(string(data))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Verified address:", address)
}

func TestValidateSignature4(t *testing.T) {

	srcAddr, err := VerifyAuthTokenAddress(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("address=%s\n", srcAddr)
}
