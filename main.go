package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
)

type Subtitle struct {
	Author      string
	Description string
	Downloads   string
	Link        string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Enter something to look for the subtitle")
		os.Exit(1)
	}

	name := os.Args[1]
	url := fmt.Sprintf("https://www.subdivx.com/index.php?accion=5&masdesc=&buscar=%s&oxdown=1", strings.Join(strings.Split(name, " "), "+"))
	c := colly.NewCollector()

	var subs []Subtitle

	c.OnHTML("#contenedor_izq", func(e *colly.HTMLElement) {
		var subtitle Subtitle

		e.ForEach("#menu_detalle_buscador", func(_ int, e *colly.HTMLElement) {
			subtitle.Link = e.ChildAttr(".titulo_menu_izq", "href")
		})

		e.ForEach("#buscador_detalle", func(_ int, e *colly.HTMLElement) {
			subtitle.Author = e.ChildText(".link1")
			// long descriptions break the list render when moving through the options
			// https://github.com/manifoldco/promptui/issues/143
			// subtitle.Description = e.ChildText("#buscador_detalle_sub")
			subtitle.Downloads = strings.Split(e.ChildText("#buscador_detalle_sub_datos"), " ")[1]

			subs = append(subs, subtitle)
		})
	})

	c.OnScraped(func(r *colly.Response) {
		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "- By: {{ .Author | cyan }} (Downloads: {{ .Downloads | cyan }})",
			Inactive: "  By: {{ .Author }} (Downloads: {{ .Downloads }})",
			Selected: "{{ .Description | red | cyan }}",
		}

		prompt := promptui.Select{
			Items:     subs,
			Label:     "Subtitles:",
			Templates: templates,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Println("Error running prompt", err.Error())
			os.Exit(1)
		}

		fmt.Println(result)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("on error", e.Error(), r.Body)
		os.Exit(1)
	})

	c.Visit(url)
}
