package main

func main() {
	cfg := &Config{
		BindAddress: ":8081",
	}
	srv := NewServer(cfg)
	srv.Run()
}
