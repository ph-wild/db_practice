package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"db_practice/internal/models"
)

func ParseOrdersFromFile(ctx context.Context, filePath string, ch chan<- models.Order) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var order models.Order
		if err := decoder.Decode(&order); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		default:
			ch <- order
		}

	}
	return nil
}
