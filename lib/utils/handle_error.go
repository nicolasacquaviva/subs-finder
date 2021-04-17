package utils

import (
	"fmt"
	"os"
)

func HandleError(e error) {
	if e != nil {
		fmt.Println("Error:", e.Error())
		os.Exit(1)
	}
}
