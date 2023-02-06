package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const APTOS_GRAPHQL = "https://indexer.mainnet.aptoslabs.com/v1/graphql"

//	curl 'https://indexer.mainnet.aptoslabs.com/v1/graphql' \
//	  -H 'content-type: application/json' \
//	  --data-raw '{"operationName":"AccountTransactionsCount","variables":{"address":"0x05a97986a9d031c4567e15b797be516910cfcb4156312482efc6a19c0a30c948"},"query":"query AccountTransactionsCount($address: String) {\n  move_resources_aggregate(\n    where: {address: {_eq: $address}}\n    distinct_on: transaction_version\n  ) {\n    aggregate {\n      count\n      __typename\n    }\n    __typename\n  }\n}"}' \
//	  --compressed
type TransactionNumber struct {
	Data struct {
		MoveResources struct {
			Aggregate struct {
				Count int `json:"count"`
			} `json:"aggregate"`
			TransactionVersion int `json:"transaction_version"`
		} `json:"move_resources_aggregate"`
	} `json:"data"`
}

// curl 'https://indexer.mainnet.aptoslabs.com/v1/graphql' \
//   -H 'content-type: application/json' \
//   --data-raw '{"operationName":"AccountTransactionsData","variables":{"address":"0x05a97986a9d031c4567e15b797be516910cfcb4156312482efc6a19c0a30c948","limit":100,"offset":0},"query":"query AccountTransactionsData($address: String, $limit: Int, $offset: Int) {\n  move_resources(\n    where: {address: {_eq: $address}}\n    order_by: {transaction_version: desc}\n    distinct_on: transaction_version\n    limit: $limit\n    offset: $offset\n  ) {\n    transaction_version\n    __typename\n  }\n}"}' \
//   --compressed

// curl 'https://fullnode.mainnet.aptoslabs.com/v1/transactions/by_version/79733494'

// 0x5b0b1c0d7a22fd7c59dedc00f906e9cb68033d13589a739b0d18cf25e317bdd4 - cetus

// {"data":{"move_resources":[{"transaction_version":79760981,"__typename":"move_resources"}, {"transaction_version":79760890,"__typename":"move_resources"},
type TransactionList struct {
	Data struct {
		MoveResources []struct {
			TransactionVersion int `json:"transaction_version"`
		} `json:"move_resources"`
	} `json:"data"`
}

func main() {
	address := "0x05a97986a9d031c4567e15b797be516910cfcb4156312482efc6a19c0a30c948"
	count, err := getTransactionCount(address)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Total number - ", count)
	offset := 0
	start := time.Now()
	for {
		versions, err := getTransactionList(address, offset, 100)
		if err != nil {
			log.Panic(err)
		}
		if len(versions) == 0 {
			break
		}
		offset += len(versions)
		time.Sleep(1 * time.Second)
	}
	fmt.Println(time.Since(start))
	fmt.Println(offset)
}

func getTransactionCount(address string) (int, error) {
	body := fmt.Sprintf(`{"operationName":"AccountTransactionsCount","variables":{"address":"%s"},"query":"query AccountTransactionsCount($address: String) {\n  move_resources_aggregate(\n    where: {address: {_eq: $address}}\n    distinct_on: transaction_version\n  ) {\n    aggregate {\n      count\n      __typename\n    }\n    __typename\n  }\n}"}`, address)
	resp, err := http.Post(APTOS_GRAPHQL, "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		return 0, err
	}
	val, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var tn TransactionNumber
	err = json.Unmarshal(val, &tn)
	if err != nil {
		return 0, err
	}
	if tn.Data.MoveResources.Aggregate.Count == 0 {
		fmt.Println(string(val))
	}
	return tn.Data.MoveResources.Aggregate.Count, nil
}

func getTransactionList(address string, offset int, limit int) ([]string, error) {
	body := fmt.Sprintf(`{"operationName":"AccountTransactionsData","variables":{"address":"%s","limit":%v,"offset":%v},"query":"query AccountTransactionsData($address: String, $limit: Int, $offset: Int) {\n  move_resources(\n    where: {address: {_eq: $address}}\n    order_by: {transaction_version: asc}\n    distinct_on: transaction_version\n    limit: $limit\n    offset: $offset\n  ) {\n    transaction_version\n    __typename\n  }\n}"}`, address, limit, offset)
	resp, err := http.Post(APTOS_GRAPHQL, "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}
	val, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var tl TransactionList
	err = json.Unmarshal(val, &tl)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(fmt.Sprintf("%s.txt", address), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := make([]string, len(tl.Data.MoveResources))
	for i, val := range tl.Data.MoveResources {
		result[i] = strconv.Itoa(val.TransactionVersion)
		_, err = file.WriteString(result[i] + "\n")
		if err != nil {
			return nil, err
		}
	}
	if len(result) == 0 {
		fmt.Println(string(val))
	}
	return result, nil
}
