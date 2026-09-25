// Package imageproc yuklangan rasmlardan responsive (srcset) variantlar yasaydi.
package imageproc

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // WebP dekoder
)

// VariantWidths — frontend srcset uchun kengliklar (asl rasmdan kichiklari yaratiladi).
var VariantWidths = []int{640, 1024, 1600, 2400}

type Variant struct {
	Width int
	File  string // dir ichidagi fayl nomi
}

type Result struct {
	Width    int
	Height   int
	Variants []Variant
}

// Dimensions — rasmni to'liq dekod qilmasdan o'lchamini oladi.
func Dimensions(data []byte) (int, int, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

// Process rasmni dekod qiladi va dir ichiga base-w{width}.{jpg|png} variantlarini yozadi.
// Shaffof rasmlar PNG, qolganlari JPEG (sifat 82) sifatida saqlanadi.
func Process(data []byte, dir, base string) (*Result, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("rasmni o‘qib bo‘lmadi: %w", err)
	}
	b := src.Bounds()
	res := &Result{Width: b.Dx(), Height: b.Dy()}

	opaque := isOpaque(src)
	for _, w := range VariantWidths {
		if w >= res.Width {
			break
		}
		h := int(float64(res.Height) * float64(w) / float64(res.Width))
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Over, nil)

		ext := ".jpg"
		if !opaque {
			ext = ".png"
		}
		name := fmt.Sprintf("%s-w%d%s", base, w, ext)
		if err := writeImage(filepath.Join(dir, name), dst, opaque); err != nil {
			return nil, err
		}
		res.Variants = append(res.Variants, Variant{Width: w, File: name})
	}
	return res, nil
}

func writeImage(path string, img image.Image, opaque bool) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if opaque {
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 82})
	}
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	return enc.Encode(f, img)
}

// isOpaque — rasmda shaffof piksellar bor-yo'qligini tekshiradi.
func isOpaque(img image.Image) bool {
	if o, ok := img.(interface{ Opaque() bool }); ok {
		return o.Opaque()
	}
	return true
}
