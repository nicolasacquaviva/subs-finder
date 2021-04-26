package collectors

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/nicolasacquaviva/subs-finder/lib/utils"
)

func ExecuteCollector(lang string, name string) {
	collector := &Collector{
		c: colly.NewCollector(),
	}
	termSize, err := utils.GetTerminalSize()
	utils.HandleError(err)
	utils.ClearConsole()

	switch lang {
	case "español":
		collector.subdivxCollect(name, termSize)
	default:
		fmt.Printf("Not yet implemented\n")
	}
}
