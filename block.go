package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp    int64
	PreBlockHash []byte
	Transactions []*Transaction
	Hash         []byte
	Nonce        int
}

func NewBlock(transactions []*Transaction, preBlockHash []byte) *Block {
	new_block := &Block{time.Now().Unix(), preBlockHash, transactions, []byte{}, 0}
	pow := NewProofOfWork(new_block)
	nonce, hash := pow.Run()
	new_block.Hash = hash[:]
	new_block.Nonce = nonce
	return new_block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) HashTransactions() []byte {

	var txHashes [][]byte
	var hash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	hash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return hash[:]
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
