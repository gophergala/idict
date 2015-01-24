package main

import (
	"fmt"
	"testing"
	"os"
)

func TestReadConfig(t *testing.T) {
	config := readConfig("settings.txt")
	fmt.Printf("%v \n", config)
}

func TestAnkiImport(t *testing.T) {
	words := []Word{Word{0, "hello", "world", nil, "", ""}}
	file, _ := os.OpenFile("test.txt", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666); defer file.Close()
	channel := make(chan string, 2)
	ankiImport(words, file, channel)
	
	<- channel; <- channel
}
