package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp    int64
	PreBlockHash []byte
	Data         []byte
	Hash         []byte
	Nonce        int
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PreBlockHash, timestamp, b.Data}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, preBlockHash []byte) *Block {
	new_block := &Block{time.Now().Unix(), preBlockHash, []byte(data), []byte{}, 0}
	pow := NewProofOfWork(new_block)
	nonce, hash := pow.Run()
	new_block.Hash = hash[:]
	new_block.Nonce = nonce
	return new_block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
