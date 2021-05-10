package usecase

import (
	"context"
	"fmt"
)

func GetSuccess(ctx context.Context) (result string, err error) {
	result = "success"
	fmt.Println(result)
	return result, nil
}
