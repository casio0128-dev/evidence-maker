package main

import (
	"encoding/json"
	"evidence-maker/conf"
	"evidence-maker/utils"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"sync"
)

func main() {
	var c *conf.Config

	confJSON, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(confJSON, &c); err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	if err := utils.OutputExcelFile(wg, c); err != nil {
		panic(err)
	}
	wg.Wait()
}
