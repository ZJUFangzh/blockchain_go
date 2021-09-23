package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	// bc *BlockChain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get the balance of address")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and gen the reward to address")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

// func (cli *CLI) addBlock(data string) {
// 	cli.bc.MineBlock(data)
// 	fmt.Println("AddBlock Success!")
// }
func (cli *CLI) createBlockChain(address string) {
	bc := NewBlockChain(address)
	bc.db.Close()
	fmt.Println("CreateBlockChain Done!")
}

func (cli *CLI) printChain() {
	// bci := cli.bc.Iterator()
	bc := NewBlockChain("")
	defer bc.db.Close()
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PreBlockHash)
		// fmt.Printf("Data: %s\n", block.transactions)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PreBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockChain(from)
	defer bc.db.Close()
	tx := NewUTXOTransactions(from, to, amount, bc)

	bc.MineBlock([]*Transaction{tx})
	fmt.Printf("Send %d coins from %s to  %s success.\n", amount, from, to)

}

func (cli *CLI) getBalance(address string) {
	bc := NewBlockChain(address)
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// addBlockData := addBlockCmd.String("data", "", "Block data")
	createBlockChainAddr := createBlockChainCmd.String("address", "", "The address to send genesis block reward to")
	getBalanceAddr := getBalanceCmd.String("address", "", "The balance of address")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	switch os.Args[1] {
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddr == "" {
			createBlockChainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockChain(*createBlockChainAddr)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddr == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddr)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
