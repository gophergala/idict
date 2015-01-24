package main

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {

}

func testReadConfig(t *testing.T) {
	config := readConfig("settings.txt")
	fmt.Printf("%v", config)
}
