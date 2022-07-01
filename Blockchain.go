package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

func (chain *Block) Print() {
	fmt.Printf("P ID HASH: %x\n", chain.PrevHash)
	fmt.Printf("MENSAJE: %s\n", chain.Data)
	fmt.Printf("ID HASH: %x\n", chain.Hash)
}

func Genesis() *Block {
	return CreateBlock("AQUI INICIO TODO", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	chain := InitBlockChain()
	base := "http://localhost:3000/api/v1/blocks"
	cont := 0
	chain.blocks[cont].Print()
	for {
		cont++
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		chain.AddBlock(text)
		toSend := make(map[string]string)
		toSend["id"] = strconv.Itoa(cont)
		toSend["data"] = text
		toSend["prevHash"] = fmt.Sprintf("%x", chain.blocks[cont-1].Hash)
		toSend["hash"] = fmt.Sprintf("%x", chain.blocks[cont].Hash)
		json_data, err := json.Marshal(toSend)
		if err != nil {
			log.Fatal(err)
		}
		postContent := bytes.NewBuffer([]byte(json_data))
		resp, err := http.Post(base, "application/json", postContent)
		if err != nil {
			panic(err)
			fmt.Println(resp.Body)
		}
		chain.blocks[cont].Print()
	}
}

// Language: go
