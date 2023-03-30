package store

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func dbDSN() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", errors.New("cannot load .env file")
	}
	var hold = make(map[string]string, 3)
	for _, v := range [3]string{"HOSTNAME", "PASSWORD", "DBNAME"} {
		if val, ok := os.LookupEnv(v); ok && val != "" {
			hold[v] = val
		} else {
			return "", errors.New("incomplete database credentials")
		}
	}
	dbLogin := fmt.Sprintf("postgres://%v:%v@localhost:5432/%v", hold["HOSTNAME"], hold["PASSWORD"], hold["DBNAME"])
	return dbLogin, nil
}
