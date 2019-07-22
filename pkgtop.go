package main

import (
	"fmt"
	"log"
	"strings"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var diskUsage map[string]int

// Convert boolean value to integer.
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Find maximum key-value pair in map.
func maxValMap(m map[string]int) string {
	var max int = 0
	var key string = ""
	for k, v := range m {
        if max < v {
			max = v
			key = k
        }
    }
    return key
}

// Find the maximum item size in a string array.
func maxArrItemSize(s []string) int {
	var arr []int
	var max int = 0
	for _, p := range s {
		arr = append(arr, len(p))
	}
	for _, l := range arr {
        if max < l {
			max = l
        }
    }
	return max
}


func getDfText(diskUsage map[string]int, width int) string {
	var diskUsageText string
	width = int(float64(width)/2.5)
	for k, v := range diskUsage {
		diskUsageText += fmt.Sprintf(" %s%s[%s %s%d%%] \n", k, 
			strings.Repeat(" ", len(maxValMap(diskUsage)) + 1 - len(k)), 
			strings.Repeat("|", (width*v)/100), 
			strings.Repeat(" ", width-(width*v)/100 + btoi(v < 10)), v)
	}
	return diskUsageText
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	diskUsage = map[string]int{
		"dev": 0,
		"run": 1,
		"/dev/sda1": 75,
		"tmpfs": 4,
	}
	termWidth, termHeight := ui.TerminalDimensions()


	dfText := widgets.NewParagraph()
	dfText.Text = getDfText(diskUsage, termWidth)
	//dfText.Border = false

	pkgText := widgets.NewParagraph()
	pkgText.Text = "~"
	//pkgText.Border = false

	pkgs := []string{
		"apache~2.4.39-1~6.25MiB~'Fri 11 Jan 2019 03:34:39'",
		"autoconf~2.69-5~2.06MiB~'Fri 11 Jan 2019 03:34:39'",
		"automake~1.16.1-1~1598.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"bind-tools~9.14.2-1~5.85MiB~'Fri 11 Jan 2019 03:34:39'",
		"bison~3.3.2-1~2013.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"brook~20190401-1~13.98MiB~'Fri 11 Jan 2019 03:34:39'",
		"chafa~1.0.1-1~327.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"cmatrix~2.0-1~95.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"compton~6.2-2~306.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"docker~1:18.09.6-1~170.98MiB~'Fri 11 Jan 2019 03:34:39'",
	}

	pd := (termWidth - maxArrItemSize(pkgs))/len(strings.Split(pkgs[0], "~"))
	for i, p := range pkgs {
		pkg := strings.Split(p, "~")
		pkgs[i] = fmt.Sprintf("%s %s %s %s %s %s %s %d", 
			pkg[0], 
			strings.Repeat(" ", pd-len(pkg[0])),
			pkg[1], 
			strings.Repeat(" ", pd-len(pkg[1])),
			pkg[2], 
			strings.Repeat(" ", pd-len(pkg[2])),
			pkg[3],
			pd)
	}

	pkgList := widgets.NewList()
	pkgList.Rows = pkgs
	
	termGrid := ui.NewGrid()
	termGrid.SetRect(0, 0, termWidth, termHeight)
	termGrid.Set(
		ui.NewRow(0.25,
			ui.NewCol(0.5, dfText),
			ui.NewCol(0.5, pkgText),
		),
		ui.NewRow(0.65,
			ui.NewCol(1.0, pkgList),
		),
		ui.NewRow(0.10,
			ui.NewCol(1.0, pkgText),
		),
	)
	ui.Render(termGrid)
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>", "<C-d>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				termGrid.SetRect(0, 0, payload.Width, payload.Height)
				dfText.Text = getDfText(diskUsage, payload.Width)
				ui.Clear()
				ui.Render(termGrid)
			}
		}
	}

}