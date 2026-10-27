package main

import (
	"fmt"
	"net/http"
)

func getUrl(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка Get запроса")
		return 0
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	return statusCode
}

func main() {

	result := getUrl("https://vk.ru/feed")
	fmt.Println("Статус:", result)
}
