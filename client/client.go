package main

import (
	"context"
	"fmt"
	util "goexpert/error"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	util.ErrorHandler(err)

	response, err := http.DefaultClient.Do(request)
	util.ErrorHandler(err)

	defer response.Body.Close()

	value, err := io.ReadAll(response.Body)
	util.ErrorHandler(err)

	file, err := os.OpenFile("cotacao.txt", os.O_CREATE|os.O_WRONLY, 0644)
	util.ErrorHandler(err)

	_, err = file.WriteString(fmt.Sprintf("Dolar: %s", string(value)))
	util.ErrorHandler(err)
}
