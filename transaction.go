package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 1000

type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// encode and set tx id
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{subsidy, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()
	return &tx
}

func (in *TxInput) CanUnlockOutputWith(unlockData string) bool {
	return in.ScriptSig == unlockData
}

func (out *TxOutput) CanBeUnlockedWith(unlockData string) bool {
	return out.ScriptPubKey == unlockData
}

func NewUTXOTransactions(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("Not enough funds")
	}

	for txId, outs := range validOutputs {
		txIdByte, err := hex.DecodeString(txId)
		if err != nil {
			log.Panicln(err)
		}
		for _, out := range outs {
			input := TxInput{txIdByte, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		fmt.Printf("acc:%d is bigger than amount:%d, create a change.\n", acc, amount)
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}
