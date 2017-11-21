package main

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/wsxiaoys/terminal/color"

	"github.com/gizak/termui"
	yaml "gopkg.in/yaml.v2"
)

const (
	navMenuHeight = 5
)

var (
	author     = flag.String("author", "", "List titles by this author.")
	title      = flag.String("title", "", "Title of the essay.")
	titleIndex = flag.Int("index", -1, "Index of the title.")
)

var (
	model          []interface{}
	currentPlaying map[interface{}]interface{}
	authorDetail   map[interface{}]interface{}
)

// Essay ...
type Essay struct {
	cmdCtx context.Context
	cmd    *exec.Cmd
	Lines  []string
}

// NewEssay ...
func NewEssay(input []byte) Essay {
	if input == nil {
		input = []byte("")
	}
	return Essay{
		Lines: strings.Split(string(input), "\n"),
	}
}

// NextDuration ...
func (e *Essay) NextDuration(index int) time.Duration {
	line := e.Lines[index]
	words := strings.Split(line, " ")
	return time.Second * time.Duration(len(words))
}

// Listen ...
func (e *Essay) Listen(index int) int {
	if index < 0 || index >= len(e.Lines) {
		return 0
	}
	line := e.Lines[index]
	duration := time.Duration(10) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	e.cmdCtx = ctx
	if e.cmd = exec.CommandContext(ctx, "say", line); e.cmd != nil {
		e.cmd.Run()
	}
	return index + 1
}

// Stop ...
func (e *Essay) Stop() {
	if e.cmdCtx != nil {
		e.cmd.Process.Kill()
		e.cmdCtx.Done()
	}
}

// Header ...
type Header struct {
	View  *termui.List
	Title string
}

// NewHeader ...
func NewHeader(title string) *Header {
	ls := termui.NewList()
	ls.Items = []string{
		"Menu: q quit. :p pause :r resume. o: open link.",
		fmt.Sprintf("Reading: [%s](fg-white) by [%s](fg-white)", title, authorDetail["fullName"].(string)),
	}
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "Welcome to Awesome Startup Essays!!"
	ls.Height = 4
	ls.Width = 50
	ls.Y = 0
	return &Header{View: ls, Title: title}
}

// Update ...
func (m *Header) Update(state string) {
	authorName := authorDetail["fullName"].(string)
	switch state {
	case "pause":
		m.View.Items[len(m.View.Items)-1] = fmt.Sprintf("[Paused](fg-red): [%s](fg-white) by %s", m.Title, authorName)
	case "resume":
		m.View.Items[len(m.View.Items)-1] = fmt.Sprintf("Reading: [%s](fg-white) by [%s](fg-white)", m.Title, authorName)
	}
	termui.Body.Align()
	termui.Render(m.View)
}

// Body ...
type Body struct {
	View  *termui.List
	essay Essay
}

// NewBody ...
func NewBody() *Body {
	ls := termui.NewList()
	ls.Items = []string{}
	ls.ItemFgColor = termui.ColorYellow
	ls.Border = false
	ls.Height = termui.TermHeight() - navMenuHeight
	ls.Width = 25
	ls.Y = 1
	return &Body{View: ls}
}

// Update ...
func (b *Body) Update(next *int, emptyLines int, essay *Essay) {
	curLineNo := *next
	line := strings.TrimSpace(essay.Lines[curLineNo])
	updateContent := func() {
		if emptyLines >= 1 {
			line = ""
		} else {
			line = fmt.Sprintf("%d [%s](fg-white)", curLineNo+1, line)
		}
		b.View.Items = append(b.View.Items, line)
		if len(b.View.Items) == b.View.Height {
			b.View.Items = []string{line}
		}
		b.View.Height = termui.TermHeight() - navMenuHeight
		termui.Body.Align()
		termui.Render(b.View)
	}
	updateContent()
}

func parseAuthorModel(author string) (map[interface{}]interface{}, []interface{}) {
	data, err := Asset(fmt.Sprintf("%s/index.yaml", author))
	if err != nil {
		panic(err)
	}
	var authors map[string]interface{}
	err = yaml.Unmarshal(data, &authors)
	if err != nil {
		fmt.Printf("Author not found.\n")
		return nil, nil
	}
	list := authors[author].([]interface{})
	first := list[0].(map[interface{}]interface{})
	return first, list[1:]
}

func printAuthorList(author string) {
	if model == nil {
		return
	}
	color.Println("@gTable of Contents")
	println()
	for index, title := range model {
		if titleMap, ok := title.(map[interface{}]interface{}); ok {
			color.Println("@w\t", index+1, " ", titleMap["title"].(string))
		}
	}
}

func main() {
	var err error
	var data []byte
	flag.Parse()
	if *author != "" {
		authorDetail, model = parseAuthorModel(*author)
	}
	if *author == "" && *title == "" && *titleIndex < 0 {
		flag.PrintDefaults()
		return
	}
	if *author != "" && *title == "" && *titleIndex < 0 {
		printAuthorList(*author)
		return
	}
	if *author != "" && *title != "" {
		var file string
		for _, modelItems := range model {
			modelItemMap := modelItems.(map[interface{}]interface{})
			if *title == modelItemMap["title"].(string) {
				currentPlaying = modelItemMap
				file = fmt.Sprintf("%s/%s", *author, modelItemMap["file"].(string))
			}
		}
		if file != "" {
			asset, err := Asset(file)
			if err != nil {
				panic(err)
			}
			data = asset
		}
	}
	if *author != "" && *titleIndex != -1 {
		var file string
		if *titleIndex <= 0 || *titleIndex > len(model) {
			flag.PrintDefaults()
			return
		}
		for index, modelItems := range model {
			modelItemMap := modelItems.(map[interface{}]interface{})
			if *titleIndex == index+1 {
				currentPlaying = modelItemMap
				*title = modelItemMap["title"].(string)
				file = fmt.Sprintf("%s/%s", *author, modelItemMap["file"].(string))
			}
		}
		if file != "" {
			asset, err := Asset(file)
			if err != nil {
				panic(err)
			}
			data = asset
		}
	}
	err = termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()
	file := *title
	header := NewHeader(file)
	body := NewBody()
	termui.Body.AddRows(
		termui.NewRow(termui.NewCol(12, 0, header.View)),
		termui.NewRow(termui.NewCol(12, 0, body.View)))
	termui.Body.Align()
	termui.Render(termui.Body)
	essay := NewEssay(data)
	next := 0
	resumeOnce, pauseOnce := sync.Once{}, sync.Once{}
	done := make(chan bool)
	play := func() {
		go func() {
			for {
				select {
				case <-done:
					resumeOnce = sync.Once{}
					essay.Stop()
					return
				case <-time.After(time.Millisecond * 100):
					if next >= len(essay.Lines) {
						essay.Stop()
						termui.StopLoop()
						return
					}
					empty := 0
					body.Update(&next, empty, &essay)
					next = essay.Listen(next)
				}
			}
		}()
	}
	termui.Handle("/sys/kbd/p", func(event termui.Event) {
		pauseOnce.Do(func() {
			header.Update("pause")
			done <- true
		})
	})
	termui.Handle("/sys/kbd/r", func(event termui.Event) {
		resumeOnce.Do(func() {
			pauseOnce = sync.Once{}
			header.Update("resume")
			play()
		})
	})
	termui.Handle("/sys/kbd/o", func(event termui.Event) {
		if link, ok := currentPlaying["link"].(string); ok {
			cmd := exec.Command("open", link)
			cmd.Start()
		}
	})
	termui.Handle("/sys/kbd/q", func(event termui.Event) {
		essay.Stop()
		termui.StopLoop()
	})
	// Autoplay
	play()
	termui.Loop()
}
