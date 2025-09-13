package utils

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	name string
	tick string

	progress, limit, width, tickSize int
}

func NewDefaultProgressBar(name string, stop, width int) ProgressBar {
	tickSize := stop / width
	return ProgressBar{name, "=", 0, stop, width, tickSize}
}

func (pb *ProgressBar) render() {
	progress := min(max(pb.progress, 0), pb.limit)

	percentage := float32(progress) / float32(pb.limit)
	fill := int(float32(pb.width-2) * percentage)
	empty := (pb.width - 2) - fill
	var bar string = "[" + strings.Repeat(pb.tick, fill) + strings.Repeat(" ", empty) + "]"

	fmt.Printf("\r%s: %s %3d%%", pb.name, bar, int(percentage*100))

	if progress == pb.limit {
		fmt.Print("\n")
	}
}

func (pb *ProgressBar) Step() {
	if pb.progress == pb.limit {
		return
	}

	pb.progress++
	pb.render()
}
