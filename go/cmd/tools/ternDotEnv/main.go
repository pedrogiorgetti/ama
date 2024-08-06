package main

import (
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	command := exec.Command("tern", "migrate", "--migrations", "./internal/db/postgres/migrations", "--config", "./internal/db/postgres/migrations/tern.conf");

	if err := command.Run(); err != nil {
		panic(err)
	}
}