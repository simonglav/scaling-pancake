package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type TransactionAnalytics struct {
	Price       string `json:"price"`
	Pair        string `json:"pair"`
	TotalAmount string `json:"total_amount"`
	Timestamp   string `json:"timestamp"`
}

type Transaction struct {
	Timestamp string   `json:"timestamp"`
	In        CoinPair `json:"in"`
	Out       CoinPair `json:"out"`
}

type CoinPair struct {
	Coin   string `json:"coin"`
	Amount string `json:"amount"`
}

// TODO price, pair, totalamount, time(minute)

func main() {
	err := calculateTotal("transactions.txt")
	if err != nil {
		panic(err)
	}
}

func calculateTotal(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	outFile, err := os.OpenFile("data.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer outFile.Close()
	scanner := bufio.NewScanner(file)
	// currentMinute := ""
	// minuteAmount := 0
	// timeMin := ""
	// price := ""

	for scanner.Scan() {
		var tr Transaction
		err := json.Unmarshal(scanner.Bytes(), &tr)
		if err != nil {
			return err
		}

	}
	return nil
}
