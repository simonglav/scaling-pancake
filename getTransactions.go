package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type TransactionRaw struct {
	Version string `json:"version"`
	Payload struct {
		Function      string        `json:"function"`
		TypeArguments []string      `json:"type_arguments"`
		Arguments     []interface{} `json:"arguments"`
	} `json:"payload"`
	Timestamp string `json:"timestamp"`
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

func main() {
	filePath := "0x05a97986a9d031c4567e15b797be516910cfcb4156312482efc6a19c0a30c948.txt"
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	outFile, err := os.OpenFile("transactions.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	counter := 0
	defer outFile.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := fmt.Sprintf("https://fullnode.mainnet.aptoslabs.com/v1/transactions/by_version/%s", scanner.Text())
		body, err := responseBody(url)
		if err != nil || len(body) < 4000 {
			fmt.Println(string(body))
			break
		}
		_, err = outFile.Write(body)
		if err != nil {
			log.Println(err)
			break
		}
		_, err = outFile.Write([]byte{'\n'})
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(200 * time.Millisecond)
		counter++
	}
	fmt.Println("Done - ", counter)
}

func responseBody(requestURL string) ([]byte, error) {
	res, err := http.Get(requestURL)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	return resBody, nil
}
