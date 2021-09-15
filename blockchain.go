package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"
const lastBlockHash = "l"

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

type BlockChainInterator struct {
	currBlockHash []byte
	db            *bolt.DB
}

func (bc *BlockChain) AddBlock(data string) {

	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte(lastBlockHash))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte(lastBlockHash), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash
		return nil
	})

}

func NewBlockChain() *BlockChain {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	var tip []byte
	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockBucket))

		if b == nil {
			b, err := tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)
			}
			genesis := NewGenesisBlock()
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte(lastBlockHash), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
			log.Println("Create genesis block success.")
		} else {
			tip = b.Get([]byte(lastBlockHash))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := &BlockChain{tip, db}
	return bc
}

func (bc *BlockChain) Iterator() *BlockChainInterator {
	return &BlockChainInterator{bc.tip, bc.db}
}

func (i *BlockChainInterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encode_block := b.Get(i.currBlockHash)
		block = DeserializeBlock(encode_block)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currBlockHash = block.PreBlockHash
	return block
}
