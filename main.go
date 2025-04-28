package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type ScriptPubKeyResult struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex,omitempty"`
	ReqSigs   int32    `json:"reqSigs,omitempty"` // Deprecated: removed in Bitcoin Core
	Type      string   `json:"type"`
	Address   string   `json:"address,omitempty"`
	Addresses []string `json:"addresses,omitempty"` // Deprecated: removed in Bitcoin Core
}

type Utxo struct {
	Hash         string             `json:"hash"`
	Idx          uint64             `json:"idx"`
	BlockNumber  uint64             `json:"bn"`
	PubKeyScript ScriptPubKeyResult `json:"pkey"`
	Value        uint64             `json:"val"`
}

type UtxoCSV struct {
	Hash        string `csv:"hash"`
	Idx         uint64 `csv:"idx"`
	BlockNumber uint64 `csv:"block_number"`
	Address     string `csv:"address"`
	Value       uint64 `csv:"value"`
}

func main() {
	// Prepare JSON file reader
	jsonFile, err := os.Open("utxodump.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)

	// Consume the opening bracket of the JSON array
	if _, err := decoder.Token(); err != nil {
		fmt.Println("Error reading JSON array:", err)
		return
	}

	// Prepare CSV file writer
	cvsFile, err := os.Create("utxodump.csv")
	if err != nil {
		fmt.Println("could not create file:", err)
		return
	}
	defer cvsFile.Close()

	writer := csv.NewWriter(cvsFile)
	defer writer.Flush()

	// Write the header
	if err := writer.Write([]string{"hash", "idx", "block_number", "address", "value"}); err != nil {
		fmt.Println("could not write header: ", err)
		return
	}

	// Loop through the array
	i := 0
	for decoder.More() {
		i++
		if i%1000 == 0 {
			fmt.Printf("\rProgress: %d%%", i)
		}

		var line Utxo
		if err := decoder.Decode(&line); err != nil {
			fmt.Println("Error decoding JSON object:", err)
			return
		}

		record := []string{line.Hash, strconv.FormatUint(line.Idx, 10), strconv.FormatUint(line.BlockNumber, 10), line.PubKeyScript.Addresses[0], strconv.FormatUint(line.Value, 10)}
		if err := writer.Write(record); err != nil {
			fmt.Println("Could not write record:", err)
			return
		}
	}

	// Read the closing bracket of the JSON array
	if _, err := decoder.Token(); err != nil {
		fmt.Println("Error reading closing bracket of JSON array:", err)
		return
	}
}
