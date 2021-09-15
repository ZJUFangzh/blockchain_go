package main

func main() {

	blockChain := NewBlockChain()
	defer blockChain.db.Close()

	cli := CLI{blockChain}
	cli.Run()
}
