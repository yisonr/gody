package main

import (
	"fmt"
	"time"
)

func main() {
	result := make(map[string]int)
	for {
		exam := Memory()
		time.Sleep(1 * time.Second)
		if _, ok := result[exam]; ok {
			continue
		}
		result[exam] = 1
		fmt.Println(exam)
		if len(result) == 5 {
			break
		}
	}
	for k, _ := range result {
		fmt.Println(k)
	}
}
