package avail

import (
	"fmt"
	"log"
	"os"

	"github.com/availproject/avail-go-sdk/primitives"
	SDK "github.com/availproject/avail-go-sdk/sdk"
	"github.com/joho/godotenv"
)

func SendDataToDA(AppId uint32, message string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	seed := os.Getenv("SEED")
	if seed == "" {
		log.Fatal("SEED environment variable is not set")
		return "", fmt.Errorf("SEED environment variable is not set")
	}

	acc, err := SDK.Account.NewKeyPair(seed)
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
		return "", err
	}
	fmt.Println("Your account Address: " + acc.SS58Address(42))

	sdk, err := SDK.NewSDK("https://turing-rpc.avail.so/rpc")
	if err != nil {
		log.Fatalf("Failed to initialize SDK: %v", err)
		return "", err
	}

	appId := uint32(AppId)
	tx := sdk.Tx.DataAvailability.SubmitData([]byte(message))
	res, err := tx.ExecuteAndWatchInclusion(acc, SDK.NewTransactionOptions().WithAppId(appId))
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
		return "", err
	}

	if !res.IsSuccessful().UnsafeUnwrap() {
		log.Fatal("Transaction was not successful")
		return "", fmt.Errorf("Transaction was not successful")
	}
	response := fmt.Sprintf(`Block Hash: https://avail-turing.subscan.io/block/%v, Block Index: %v, Tx Hash: https://avail-turing.subscan.io/extrinsic/%v, Tx Index: %v`, res.BlockHash.ToHexWith0x(), res.BlockNumber, res.TxHash.ToHexWith0x(), res.TxIndex)
	fmt.Println(response)
	fmt.Println("Data submission completed successfully")
	return response, nil
}

func GetDataFromDA(txn_hash string, block_hash string) (string, error) {
	sdk, err := SDK.NewSDK("https://turing-rpc.avail.so/rpc")
	if err != nil {
		log.Fatalf("Failed to initialize SDK: %v", err)
	}

	blockHash, err := primitives.NewBlockHashFromHexString(block_hash)
	if err != nil {
		log.Fatalf("Failed to create block hash: %v", err)
	}

	block, err := SDK.NewBlock(sdk.Client, blockHash)
	if err != nil {
		log.Fatalf("Failed to create block: %v", err)
	}

	txHash, err := primitives.NewH256FromHexString(txn_hash)
	if err != nil {
		log.Fatalf("Failed to create transaction hash: %v", err)
	}

	blobs := block.DataSubmissions(SDK.Filter{}.WTxHash(txHash))
	if err != nil {
		log.Fatalf("Failed to create data submissions: %v", err)
	}

	blob := &blobs[0]

	accountId, err := primitives.NewAccountIdFromMultiAddress(blob.TxSigner)
	if err != nil {
		log.Fatalf("Failed to create account ID: %v", err)
		return "", fmt.Errorf("Transaction was not successful")
	}

	response := fmt.Sprintf(`Tx Hash: %v, Tx Index: %v, Data: %v, App Id: %v, Signer: %v,`, blob.TxHash, blob.TxIndex, string(blob.Data), blob.AppId, accountId.ToHuman())
	fmt.Println(response)

	fmt.Println("Data retrieval completed successfully")

	return response, nil
}
