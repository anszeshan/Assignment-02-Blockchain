//Ans Zeshan Shahid
//20I-0543
//Assignment 02

package main

import (
	"crypto/sha256" // Imports the sha256 package for cryptographic hashing
	"encoding/hex"  // Imports the hex package for encoding and decoding hexadecimal numbers
	"fmt"           // Imports the fmt package for formatting and printing output

	//"https://github.com/anszeshan/assignment02bca.git"       //github
	"strconv" // Imports the strconv package for converting strings to numbers and vice versa
	"strings" // Imports the strings package for manipulating strings
	"time"    // Imports the time package for working with time and dates
)

// Transaction represents a transaction in the blockchain
type Transaction struct {
	Sender   string  // The sender's address
	Receiver string  // The receiver's address
	Amount   float64 // The amount of cryptocurrency being transferred
}

// Block represents a block in the blockchain
type Block struct {
	Index        int           // The block's index in the blockchain
	Timestamp    time.Time     // The time at which the block was mined
	Transactions []Transaction // A list of transactions in the block
	Nonce        int           // A random number used to mine the block
	PreviousHash string        // The hash of the previous block in the blockchain
	Hash         string        // The hash of the current block
}

// Blockchain represents the blockchain
type Blockchain struct {
	Blocks                 []*Block // A list of blocks in the blockchain
	NumberOfTransactions   int      // The total number of transactions in the blockchain
	BlockHashRangeMinValue int      // The minimum value of a block's hash
	BlockHashRangeMaxValue int      // The maximum value of a block's hash
}

// MerkleTree represents the Merkle Tree
type MerkleTree struct {
	Root *MerkleNode // The root node of the Merkle Tree
}

// MerkleNode represents a node in the Merkle Tree
type MerkleNode struct {
	Left  *MerkleNode // The left child node
	Right *MerkleNode // The right child node
	Data  string      // The data stored in the node
}

// NewBlock creates a new block and adds it to the blockchain
func (bc *Blockchain) NewBlock(transaction string, nonce int, previousHash string, transactionDateTime time.Time) {
	block := &Block{
		Index:        len(bc.Blocks),                  // Sets the block's index to the length of the blockchain
		Timestamp:    transactionDateTime,             // Sets the block's timestamp to the current time
		Transactions: bc.getTransactions(transaction), // Gets the transactions for the block
		Nonce:        nonce,                           // Sets the block's nonce to the given nonce
		PreviousHash: previousHash,                    // Sets the block's previous hash to the given previous hash
	}

	block.Hash = block.CreateHash()      // Calculates the block's hash
	bc.Blocks = append(bc.Blocks, block) // Appends the block to the blockchain
}

// DisplayBlocks prints all the blocks in the blockchain
func (bc *Blockchain) DisplayBlocks() {
	for _, block := range bc.Blocks { // Iterates over the blocks in the blockchain
		fmt.Printf("Block %d\n", block.Index)                 // Prints the block's index
		fmt.Printf("Transactions: %v\n", block.Transactions)  // Prints the block's transactions
		fmt.Printf("Nonce: %d\n", block.Nonce)                // Prints the block's nonce
		fmt.Printf("Previous Hash: %s\n", block.PreviousHash) // Prints the block's previous hash
		fmt.Printf("Hash: %s\n", block.Hash)                  // Prints the block's hash
		fmt.Println("-------------------------------------")  // Prints a separator
	}
}

// ChangeBlock changes the transaction of a specified block
func (bc *Blockchain) ChangeBlock(blockIndex int, newTransaction string) {
	if blockIndex >= 0 && blockIndex < len(bc.Blocks) { // Checks if the block index is valid
		bc.Blocks[blockIndex].Transactions = bc.getTransactions(newTransaction) // Replaces the block's transactions with the new transaction
		bc.Blocks[blockIndex].Hash = bc.Blocks[blockIndex].CreateHash()         // Updates the block's hash
	}
}

// VerifyChain verifies the integrity of the blockchain
func (bc *Blockchain) VerifyChain() bool {
	for i := 1; i < len(bc.Blocks); i++ { // Iterates over the blocks in the blockchain, starting from the second block
		currentBlock := bc.Blocks[i]    // Gets the current block
		previousBlock := bc.Blocks[i-1] // Gets the previous block

		if currentBlock.Hash != currentBlock.CreateHash() { // Checks if the current block's hash is equal to its calculated hash
			return false
		}

		if currentBlock.PreviousHash != previousBlock.Hash { // Checks if the current block's previous hash is equal to the previous block's hash
			return false
		}
	}

	return true
}

// CalculateHash calculates the hash of a provided string
func CalculateHash(stringToHash string) string {
	hash := sha256.Sum256([]byte(stringToHash)) //Calculates the SHA-256 hash of the string
	return hex.EncodeToString(hash[:])          // Converts the hash to a hexadecimal string and returns it
}

// setNumberOfTransactionsPerBlock sets the number of transactions to be included in a new block during block creation
func (bc *Blockchain) setNumberOfTransactionsPerBlock(numTransactions int) {
	bc.NumberOfTransactions = numTransactions // Sets the number of transactions per block
}

// setBlockHashRangeForBlockCreation sets the range for the hash values that a new block's hash must fall within for block creation
func (bc *Blockchain) setBlockHashRangeForBlockCreation(min, max int) {
	bc.BlockHashRangeMinValue = min // Sets the minimum hash value for block creation
	bc.BlockHashRangeMaxValue = max // Sets the maximum hash value for block creation
}

// CreateHash calculates the hash of the block
func (b *Block) CreateHash() string {
	hashData := fmt.Sprintf("%d%d%v%v%s", b.Index, b.Nonce, b.Timestamp, b.Transactions, b.PreviousHash)
	return CalculateHash(hashData)
}

// getTransactions converts the transaction string into a slice of Transaction structs
func (bc *Blockchain) getTransactions(transactionString string) []Transaction {
	transactions := strings.Split(transactionString, ";")
	result := make([]Transaction, 0, len(transactions))

	for _, transaction := range transactions {
		parts := strings.Split(transaction, ",")

		amount, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			continue
		}

		result = append(result, Transaction{
			Sender:   parts[0],
			Receiver: parts[1],
			Amount:   amount,
		})
	}

	return result
}

// ProofOfWork performs the Proof of Work consensus mechanism
func ProofOfWork(block *Block, targetPrefix string) int {
	nonce := 0
	for {
		block.Nonce = nonce
		block.Hash = block.CreateHash()

		if strings.HasPrefix(block.Hash, targetPrefix) {
			return nonce
		}

		nonce++
	}
}

// CreateMerkleTree creates a Merkle Tree from the transactions in a block
func CreateMerkleTree(transactions []Transaction) *MerkleTree {
	nodes := make([]*MerkleNode, len(transactions))

	for i, transaction := range transactions {
		nodes[i] = &MerkleNode{
			Data: CalculateHash(fmt.Sprintf("%s,%s,%f", transaction.Sender, transaction.Receiver, transaction.Amount)),
		}
	}

	if len(nodes) == 0 {
		return &MerkleTree{}
	}

	for len(nodes) > 1 {
		parentLevel := make([]*MerkleNode, 0, (len(nodes)+1)/2)

		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				parent := &MerkleNode{
					Left:  nodes[i],
					Right: nodes[i+1],
					Data:  CalculateHash(nodes[i].Data + nodes[i+1].Data),
				}
				parentLevel = append(parentLevel, parent)
			} else {
				parentLevel = append(parentLevel, nodes[i])
			}
		}

		nodes = parentLevel
	}

	return &MerkleTree{
		Root: nodes[0],
	}
}

func main() {
	// Create a new blockchain
	blockchain := &Blockchain{
		Blocks:                 make([]*Block, 0),
		NumberOfTransactions:   4,
		BlockHashRangeMinValue: 0,
		BlockHashRangeMaxValue: 1000,
	}

	// Add some initial transactions
	transactionString := "Alice,Bob,1.5;Charlie,Alice,2.0;Bob,Charlie,0.5;Dave,Alice,1.0"
	blockchain.NewBlock(transactionString, 0, "", time.Now())

	// Display the blocks
	blockchain.DisplayBlocks()

	// Change the transaction of a block
	blockchain.ChangeBlock(0, "Bob,Alice,1.0;Charlie,Bob,0.5;Dave,Alice,0.8")

	// Verify the blockchain's integrity
	fmt.Println("Blockchain is valid:", blockchain.VerifyChain())

	// Calculate the hash of a string
	stringToHash := "Hello, world!"
	hash := CalculateHash(stringToHash)
	fmt.Println("Hash:", hash)

	// Add more transactions and create new blocks
	transactionString = "Bob,Alice,0.8;Charlie,Alice,1.2;Dave,Charlie,0.7;Eve,Bob,0.3"
	blockchain.NewBlock(transactionString, 1, blockchain.Blocks[0].Hash, time.Now())

	transactionString = "Alice,Charlie,2.5;Dave,Bob,0.3;Eve,Alice,1.0"
	blockchain.NewBlock(transactionString, 2, blockchain.Blocks[1].Hash, time.Now())

	transactionString = "Bob,Eve,0.5;Charlie,Dave,0.9;Eve,Bob,1.2"
	blockchain.NewBlock(transactionString, 3, blockchain.Blocks[2].Hash, time.Now())

	transactionString = "Alice,Charlie,1.5;Dave,Bob,0.2;Eve,Alice,0.7"
	blockchain.NewBlock(transactionString, 4, blockchain.Blocks[3].Hash, time.Now())

	transactionString = "Bob,Eve,0.3;Charlie,Dave,1.1"
	blockchain.NewBlock(transactionString, 5, blockchain.Blocks[4].Hash, time.Now())

	// Display the updated blockchain
	blockchain.DisplayBlocks()

	// Verify the updated blockchain's integrity
	fmt.Println("Blockchain is valid:", blockchain.VerifyChain())

	fmt.Println("\n\n")

	// Perform Proof of Work and create a Merkle Tree for all blocks
	for _, block := range blockchain.Blocks {
		targetPrefix := "0000"
		nonce := ProofOfWork(block, targetPrefix)
		fmt.Printf("Proof of Work - Block %d, Nonce: %d\n", block.Index, nonce)

		transactions := block.Transactions
		merkleTree := CreateMerkleTree(transactions)
		fmt.Printf("Merkle Tree Root - Block %d: %s\n", block.Index, merkleTree.Root.Data)
	}
}
