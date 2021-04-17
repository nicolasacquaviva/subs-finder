package collectors

import (
	"github.com/gocolly/colly/v2"
	"github.com/nicolasacquaviva/subs-finder/lib/collectors/spanish"
	"github.com/nicolasacquaviva/subs-finder/lib/utils"
)

func ExecuteCollector(lang string, name string) {
	c := colly.NewCollector()
	termSize, err := utils.GetTerminalSize()
	utils.HandleError(err)
	utils.ClearConsole()

	if lang == "español" {
		spanish.SubdivxCollect(c, name, termSize)
	}
}
