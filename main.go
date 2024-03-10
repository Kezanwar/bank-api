package main

func main() {
	server := NewApiServer(":8000")
	server.Run()
}
