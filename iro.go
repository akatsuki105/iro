package iro

import (
	"fmt"
	"image/color"
)

type ID int

const (
	ColorInvalid ID = iota
	ColorRGB555
)

type Color interface {
	color.Color
	fmt.Stringer
}

func RGBA(c Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{byte(r >> 24), byte(g >> 24), byte(b >> 24), byte(a >> 24)}
}

func normalize(c uint64, sourceDepth uint, targetDepth uint) uint64 {
	if sourceDepth == 0 || targetDepth == 0 {
		return 0
	}

	for sourceDepth < targetDepth {
		c = (c << sourceDepth) | c
		sourceDepth += sourceDepth
	}

	if targetDepth < sourceDepth {
		c >>= (sourceDepth - targetDepth)
	}

	return c
}
