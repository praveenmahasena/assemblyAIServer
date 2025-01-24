package internal

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/praveenmahasena/server/internal/server"
)

func Start() error {
	if envErr := godotenv.Load(".env"); envErr != nil {
		return fmt.Errorf("error during loading up env file %v", envErr)
	}
	nw := server.New("tcp", ":42069")
	return nw.Run()
}
