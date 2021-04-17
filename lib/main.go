package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/nicolasacquaviva/subs-finder/lib/collectors"
	"github.com/nicolasacquaviva/subs-finder/lib/utils"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: subs-finder 'finding dory'")
		os.Exit(1)
	}

	name := os.Args[1]
	langPrompt := promptui.Select{
		Label: "Language",
		Items: []string{"english", "español"},
	}
	_, lang, err := langPrompt.Run()
	utils.HandleError(err)

	collectors.ExecuteCollector(lang, name)
}
