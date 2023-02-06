package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	err := parseCetusTransactionData("transactions.txt")
	if err != nil {
		panic(err)
	}
}

func parseLiquidTransactionData(path string) error {
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
	for scanner.Scan() {
		var ctr TransactionRaw
		err := json.Unmarshal(scanner.Bytes(), &ctr)
		if err != nil {
			return err
		}
		if ctr.Payload.Function != "0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12::scripts_v2::swap" {
			fmt.Println(ctr.Payload.Function)
			continue
		}
		if len(ctr.Payload.TypeArguments) < 2 || len(ctr.Payload.Arguments) < 2 {
			log.Fatal(ctr.Version)
		}
		cpIn := CoinPair{
			Coin:   ctr.Payload.TypeArguments[0],
			Amount: (ctr.Payload.Arguments[0]).(string),
		}
		cpOut := CoinPair{
			Coin:   ctr.Payload.TypeArguments[1],
			Amount: (ctr.Payload.Arguments[1]).(string),
		}
		cr := Transaction{
			Timestamp: ctr.Timestamp,
			In:        cpIn,
			Out:       cpOut,
		}
		output, err := json.Marshal(cr)
		if err != nil {
			return nil
		}
		_, err = outFile.Write(output)
		if err != nil {
			return nil
		}
		_, err = outFile.Write([]byte{'\n'})
		if err != nil {
			return nil
		}
	}
	return nil
}

func parseCetusTransactionData(path string) error {
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
	for scanner.Scan() {
		var ctr TransactionRaw
		err := json.Unmarshal(scanner.Bytes(), &ctr)
		if err != nil {
			return err
		}
		if ctr.Payload.Function != "0xa7f01413d33ba919441888637ca1607ca0ddcbfa3c0a9ddea64743aaa560e498::clmm_router::swap" {
			fmt.Println(ctr.Payload.Function)
			continue
		}
		if len(ctr.Payload.TypeArguments) < 2 || len(ctr.Payload.Arguments) < 7 {
			log.Fatal(ctr.Version)
		}
		cpIn := CoinPair{
			Coin:   ctr.Payload.TypeArguments[0],
			Amount: (ctr.Payload.Arguments[3]).(string),
		}
		cpOut := CoinPair{
			Coin:   ctr.Payload.TypeArguments[1],
			Amount: (ctr.Payload.Arguments[4]).(string),
		}
		cr := Transaction{
			Timestamp: ctr.Timestamp,
			In:        cpIn,
			Out:       cpOut,
		}
		output, err := json.Marshal(cr)
		if err != nil {
			return nil
		}
		_, err = outFile.Write(output)
		if err != nil {
			return nil
		}
		_, err = outFile.Write([]byte{'\n'})
		if err != nil {
			return nil
		}
	}
	return nil
}
