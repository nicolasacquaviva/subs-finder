package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/manifoldco/promptui"
	"io"
	"net/http"
	"os"
	"strings"
)

type Subtitle struct {
	Author      string
	Description string
	Downloads   string
	Link        string
}

func handleError(e error) {
	if e != nil {
		fmt.Println("Error:", e.Error())
		os.Exit(1)
	}
}

func downloadFile(filePath string, fileUrl string) error {
	res, err := http.Get(fileUrl)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	out, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, res.Body)

	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: subs-finder 'finding dory'")
		os.Exit(1)
	}

	name := os.Args[1]
	url := fmt.Sprintf("https://www.subdivx.com/index.php?accion=5&masdesc=&buscar=%s&oxdown=1", strings.Join(strings.Split(name, " "), "+"))
	c := colly.NewCollector()
	renderOptions := true

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

		if renderOptions {
			i, _, err := prompt.Run()

			handleError(err)

			renderOptions = false

			c.Visit(subs[i].Link)
		}

	})

	c.OnHTML("#detalle_datos", func(e *colly.HTMLElement) {
		// download => https://www.subdivx.com/sub8/482353.rar
		id := strings.Replace(strings.Split(e.ChildAttr(".link1", "href"), "=")[1], "&u", "", 1)
		downloadLink := fmt.Sprintf("https://www.subdivx.com/sub8/%s.rar", id)

		err := downloadFile("./"+id+".rar", downloadLink)

		handleError(err)
	})

	c.OnError(func(r *colly.Response, e error) {
		handleError(e)
	})

	c.Visit(url)
}
