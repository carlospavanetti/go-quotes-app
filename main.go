package main

func main() {
	client := CreateRedisClient()
	srv := StartHttpServer(client)
	WaitForShutdown(srv)
}
