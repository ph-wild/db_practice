package services

import (
	"encoding/json"
	"os"
	"db_practice/internal/models"
)

func ParseOrdersFromFile(filePath string, ch chan<- models.Order) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var order models.Order
		if err := decoder.Decode(&order); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		ch <- order
	}

	close(ch)
	return nil
}
