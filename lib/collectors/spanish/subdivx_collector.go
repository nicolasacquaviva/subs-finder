package spanish

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/manifoldco/promptui"
	"github.com/nicolasacquaviva/subs-finder/lib/utils"
	"net/url"
	"strings"
)

type Subtitle struct {
	Author      string
	Description string
	Downloads   string
	Label       string
	Link        string
}

func SubdivxCollect(c *colly.Collector, name string, termSize *utils.TermSize) {
	var subs []Subtitle
	baseUrl := fmt.Sprintf(
		"https://www.subdivx.com/index.php?accion=5&masdesc=&buscar=%s&oxdown=1",
		strings.Join(strings.Split(name, " "), "+"),
	)
	renderOptions := true

	c.OnHTML("#contenedor_izq", func(e *colly.HTMLElement) {
		var subtitle Subtitle
		var index = len(subs)

		e.ForEach("#menu_detalle_buscador", func(_ int, e *colly.HTMLElement) {
			subtitle.Link = e.ChildAttr(".titulo_menu_izq", "href")
			subs = append(subs, subtitle)
		})

		e.ForEach("#buscador_detalle", func(_ int, e *colly.HTMLElement) {
			subtitle.Author = e.ChildText(".link1")
			// long descriptions break the list render when moving through the options
			// https://github.com/manifoldco/promptui/issues/143
			subs[index].Description = e.ChildText("#buscador_detalle_sub")
			subs[index].Label = "Search term " + name
			subs[index].Downloads = strings.Split(e.ChildText("#buscador_detalle_sub_datos"), " ")[1]

			// the description is usually the longest field, this is to avoid the issue stated above
			// shortening the description so the row content doesn't break into a new line
			const X_OFFSET = 40
			var textLength = (len(subs[index].Author) +
				len(subs[index].Description) +
				len(subs[index].Downloads) + X_OFFSET)
			var lengthForDescription = (termSize.X -
				len(subs[index].Author) -
				len(subs[index].Downloads) -
				X_OFFSET)

			if textLength > termSize.X {
				subs[index].Description = subs[index].Description[:lengthForDescription] + "..."
			}
			index++
		})
	})

	c.OnScraped(func(r *colly.Response) {
		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "- {{ .Description }} By: {{ .Author | cyan }} (Downloads: {{ .Downloads | cyan }})",
			Inactive: "  {{ .Description }} By: {{ .Author }} (Downloads: {{ .Downloads }})",
			Selected: "{{ .Description | red }}",
		}

		label := fmt.Sprintf("Subtitles for '%s':", name)
		prompt := promptui.Select{
			Items:     subs,
			Size:      10,
			Label:     label,
			Templates: templates,
		}

		if renderOptions {
			i, _, err := prompt.Run()
			utils.HandleError(err)
			renderOptions = false
			c.Visit(subs[i].Link)
		}
	})

	c.OnHTML("#detalle_datos", func(e *colly.HTMLElement) {
		// download => https://www.subdivx.com/sub8/482353.rar
		u, err := url.Parse(e.ChildAttr(".link1", "href"))
		utils.HandleError(err)
		queryParams, _ := url.ParseQuery(u.RawQuery)
		id := queryParams["id"][0]
		subdir := queryParams["u"][0]

		if subdir == "1" {
			subdir = ""
		}
		downloadLink := fmt.Sprintf(
			"https://www.subdivx.com/sub%s/%s", subdir, id,
		)

		err = utils.DownloadFile("./"+id, downloadLink, ".rar")

		utils.HandleError(err)
	})

	c.OnError(func(r *colly.Response, e error) {
		utils.HandleError(e)
	})

	c.Visit(baseUrl)
}
