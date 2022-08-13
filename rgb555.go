package iro

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// RGB555 is GameBoy color format.
//
// uint16 = Bit0-4(R) | Bit5-9(G) | Bit10-14(B)
type RGB555 uint16

func (r RGB555) String() string {
	rgba := r.Color()
	return fmt.Sprintf("#%02X%02X%02X", rgba.R, rgba.G, rgba.B)
}

// RGB gets color intensity from 0 to 31(0b11111)
func (r RGB555) intensity() (red byte, green byte, blue byte) {
	red, green, blue = byte(uint16(r)&0b11111), byte(uint16(r>>5)&0b11111), byte(uint16(r>>10)&0b11111)
	return red, green, blue
}

// Color gets 8bit RGB data
func (r RGB555) Color() color.RGBA {
	return RGBA(r)
}

func (r RGB555) RGBA() (uint32, uint32, uint32, uint32) {
	r5, g5, b5 := r.intensity()
	red, green, blue := normalize(uint64(r5), 5, 32), normalize(uint64(g5), 5, 32), normalize(uint64(b5), 5, 32)
	return uint32(red), uint32(green), uint32(blue), uint32(0xffffffff)
}

// intensity is 0 ~ 16(whiteout)
func (r RGB555) Brighten(intensity int) RGB555 {
	a := uint(r & 0x1f)
	c := (a + ((0x1f-a)*uint(intensity))/16) & 0x1f

	a = uint(r & 0x3e0)
	c |= (a + ((0x3e0-a)*uint(intensity))/16) & 0x3e0

	a = uint(r & 0x7c00)
	c |= (a + ((0x7c00-a)*uint(intensity))/16) & 0x7c00

	return RGB555(c)
}

// intensity is 0 ~ 16(blackout)
func (r RGB555) Darken(intensity int) RGB555 {
	a := uint(r & 0x1f)
	c := (a - (a*uint(intensity))/16) & 0x1f

	a = uint(r & 0x3e0)
	c |= (a - (a*uint(intensity))/16) & 0x3e0

	a = uint(r & 0x7c00)
	c |= (a - (a*uint(intensity))/16) & 0x7c00

	return RGB555(c)
}

// Mix two RGB555.
// Weight parameters are 0.0 ~ 1.0
func MixRGB555(aWeight float64, aColor RGB555, bWeight float64, bColor RGB555) RGB555 {
	ar := float64(aColor & 0x001f)
	ag := float64((aColor & 0x03e0) >> 5)
	ab := float64((aColor & 0x7c00) >> 10)

	br := float64(bColor & 0x001f)
	bg := float64((bColor & 0x03e0) >> 5)
	bb := float64((bColor & 0x7c00) >> 10)

	r := RGB555(math.Min(aWeight*ar+bWeight*br, 0x1f))
	g := RGB555(math.Min(aWeight*ag+bWeight*bg, 0x1f))
	b := RGB555(math.Min(aWeight*ab+bWeight*bb, 0x1f))

	return r | (g << 5) | (b << 10)
}

func (r RGB555) RGB565() uint16 {
	red, green, blue := r.intensity()
	return (uint16(red) << 11) | (uint16(green) << 6) | uint16(blue)
}

// fn is convert func. When fn is nil, apply Color().
func RGB555ToImage(data []RGB555, w, h int, fn func(RGB555) color.RGBA) *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var c color.RGBA
			if fn != nil {
				c = fn(data[y*w+x])
			} else {
				c = data[y*w+x].Color()
			}
			result.Set(x, y, c)
		}
	}

	return result
}
