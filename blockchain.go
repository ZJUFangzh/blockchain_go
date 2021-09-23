package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"
const lastBlockHash = "l"
const genesisBlockData = "This is the first trans for Fang"

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

type BlockChainInterator struct {
	currBlockHash []byte
	db            *bolt.DB
}

func (bc *BlockChain) MineBlock(transactions []*Transaction) {

	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte(lastBlockHash))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(transactions, lastHash)

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

func NewBlockChain(address string) *BlockChain {
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
			coinbase := NewCoinbaseTx(address, genesisBlockData)
			genesis := NewGenesisBlock(coinbase)
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

func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspendTxs []Transaction
	spendTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)

		Outputs:
			for id, out := range tx.Vout {
				if spendTXOs[txId] != nil {
					for _, spendOutId := range spendTXOs[txId] {
						if spendOutId == id {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspendTxs = append(unspendTxs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith((address)) {
						inId := hex.EncodeToString(in.Txid)
						spendTXOs[inId] = append(spendTXOs[inId], in.Vout)
					}
				}
			}
		}
		if len(block.PreBlockHash) == 0 {
			break
		}
	}
	return unspendTxs

}

func (bc *BlockChain) FindUTXO(address string) []TxOutput {
	var unspendUTXOs []TxOutput
	unspendTxs := bc.FindUnspentTransactions(address)
	for _, tx := range unspendTxs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				unspendUTXOs = append(unspendUTXOs, out)
			}
		}
	}
	return unspendUTXOs
}

func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspendOutputs := make(map[string][]int)
	unspendTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspendTxs {

		txId := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				unspendOutputs[txId] = append(unspendOutputs[txId], outIdx)
				accumulated += out.Value
			}
			if accumulated >= amount {
				break Work
			}
		}
	}
	fmt.Printf("Find enough amount: %d in %d tx.\n", accumulated, len(unspendOutputs))
	return accumulated, unspendOutputs
}
