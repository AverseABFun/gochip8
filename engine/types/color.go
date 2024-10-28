package types

import (
	"errors"
	"math"
)

type uint6 uint8 // shhh...

type Color struct {
	R uint6
	G uint6
	B uint6
}

func Lerp(a float64, b float64, f float64) float64 {
	return a*(1.0-f) + (b * f)
}

func ColorLerp(color1, color2 Color, t float64) Color {
	color1.R = uint6(math.Round(Lerp(float64(color1.R), float64(color2.R), t)))
	color1.G = uint6(math.Round(Lerp(float64(color1.G), float64(color2.G), t)))
	color1.B = uint6(math.Round(Lerp(float64(color1.B), float64(color2.B), t)))
	return color1
}

type Gradient struct {
	Colors []PaletteIndex
}

type PaletteIndex uint8

func (clr Color) IsValid() bool {
	return clr.R <= MAX_UINT6 && clr.G <= MAX_UINT6 && clr.B <= MAX_UINT6
}

var (
	InvalidColor = Color{R: 255, G: 255, B: 255}
)

const MAX_UINT6 uint6 = (1 << 6) - 1

var (
	ErrInvalidUint6 = errors.New("invalid uint6")
	ErrInvalidColor = errors.New("invalid color")
)

func FromRGB(r uint6, g uint6, b uint6) (Color, error) {
	if r > MAX_UINT6 {
		return Color{}, ErrInvalidUint6
	}
	if g > MAX_UINT6 {
		return Color{}, ErrInvalidUint6
	}
	if b > MAX_UINT6 {
		return Color{}, ErrInvalidUint6
	}
	return Color{R: r, G: g, B: b}, nil
}

func FromRGBNoErr(r uint6, g uint6, b uint6) Color {
	var out, err = FromRGB(r, g, b)
	if err != nil {
		panic(err)
	}
	return out
}