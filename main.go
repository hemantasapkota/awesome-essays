package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gizak/termui"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	usage         = `awesome-essays -file <<filename>>`
	navMenuHeight = 5
)

//Essay
type Essay struct {
	cmdCtx context.Context
	cmd    *exec.Cmd

	Lines []string
}

func NewEssay(input []byte) Essay {
	if input == nil {
		input = []byte("")
	}
	return Essay{
		Lines: strings.Split(string(input), "\n"),
	}
}

func (e *Essay) NextDuration(index int) time.Duration {
	line := e.Lines[index]
	words := strings.Split(line, " ")
	return time.Second * time.Duration(len(words))
}

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

func (e *Essay) Stop() {
	if e.cmdCtx != nil {
		e.cmd.Process.Kill()
		e.cmdCtx.Done()
	}
}

//Menu
type Header struct {
	View  *termui.List
	Title string
}

func NewHeader(title string) *Header {
	ls := termui.NewList()
	ls.Items = []string{
		"Menu: q quit. :p pause :r resume. o: open link.",
		fmt.Sprintf("Reading: [%s](fg-white)", title),
	}
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "Welcome to Awesome Startup Essays!!"
	ls.Height = 4
	ls.Width = 50
	ls.Y = 0
	return &Header{View: ls, Title: title}
}

func (m *Header) Update(state string) {
	switch state {
	case "pause":
		m.View.Items[len(m.View.Items)-1] = fmt.Sprintf("[Paused](fg-red): [%s](fg-white)", m.Title)
	case "resume":
		m.View.Items[len(m.View.Items)-1] = fmt.Sprintf("Reading: [%s](fg-white)", m.Title)
	}
	termui.Body.Align()
	termui.Render(m.View)
}

//Body
type Body struct {
	View  *termui.List
	essay Essay
}

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

func (b *Body) Update(next *int, essay *Essay) {
	emptyCount, lineNo := 0, 0
	updateContent := func() {
		formatLine := func(lineNo int, line string) (out string) {
			out = fmt.Sprintf("%d [%s](fg-white)", lineNo, line)
			return
		}

		if len(b.View.Items) == b.View.Height {
			b.View.Items = []string{formatLine(lineNo, essay.Lines[*next])}
		} else {
			line := essay.Lines[*next]
			lineNo = (*next + 1) - emptyCount

			if strings.TrimSpace(line) != "" {
				line = formatLine(lineNo, line)
			} else {
				emptyCount++
			}
			b.View.Items = append(b.View.Items, line)
		}

		b.View.Height = termui.TermHeight() - navMenuHeight
		termui.Body.Align()
		termui.Render(b.View)
	}

	updateContent()
}

func main() {
	file := flag.String("file", "", "Filename")
	flag.Parse()

	if *file == "" {
		fmt.Println(usage)
		return
	}

	var err error
	var data []byte

	if *file != "" {
		data, err = ioutil.ReadFile(*file)
		if err != nil {
			return
		}
	}

	err = termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	header := NewHeader(*file)
	body := NewBody()

	termui.Body.AddRows(
		termui.NewRow(termui.NewCol(12, 0, header.View)),
		termui.NewRow(termui.NewCol(12, 0, body.View)))

	termui.Body.Align()
	termui.Render(termui.Body)

	essay := NewEssay(data)
	next := 0
	body.Update(&next, &essay)
	next = essay.Listen(next)

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

					body.Update(&next, &essay)
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

	termui.Handle("/sys/kbd/q", func(event termui.Event) {
		essay.Stop()
		termui.StopLoop()
	})

	// Autoplay
	play()
	termui.Loop()
}
