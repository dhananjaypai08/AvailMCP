package avail

import (
	"fmt"
	"log"
	"os"

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

	// Submit data to Avail
	appId := uint32(AppId)
	tx := sdk.Tx.DataAvailability.SubmitData([]byte(message))
	res, err := tx.ExecuteAndWatchInclusion(acc, SDK.NewTransactionOptions().WithAppId(appId))
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
		return "", err
	}

	// Check if the transaction was successful
	if !res.IsSuccessful().UnsafeUnwrap() {
		log.Fatal("Transaction was not successful")
		return "", fmt.Errorf("Transaction was not successful")
	}
	response := fmt.Sprintf(`Block Hash: https://avail-turing.subscan.io/block/%v, Block Index: %v, Tx Hash: https://avail-turing.subscan.io/extrinsic/%v, Tx Index: %v`, res.BlockHash.ToHexWith0x(), res.BlockNumber, res.TxHash.ToHexWith0x(), res.TxIndex)
	// Printing out all the values of the transaction
	fmt.Println(response)
	fmt.Println("Data submission completed successfully")
	return response, nil
}
