package main

import (
	"fmt"
	"image"
	"image/draw"
	"math"

	"gocv.io/x/gocv"
)

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
}

type Line struct {
	X, Y, W, H int
}

type FisherDetector struct{}

func NewFisherDetector() *FisherDetector {
	return &FisherDetector{}
}

func (d *FisherDetector) FindFishingMinigameBounds(img gocv.Mat) (*Bounds, error) {
	if img.Empty() {
		return nil, fmt.Errorf("input image is empty")
	}

	screenW := img.Cols()
	screenH := img.Rows()

	// restrict to low-middle portion of screen
	roiY := screenH * 45 / 100
	roi := img.Region(image.Rect(0, roiY, screenW, screenH))
	defer roi.Close()

	// conv hsv
	hsv := gocv.NewMat()
	defer hsv.Close()
	gocv.CvtColor(roi, &hsv, gocv.ColorBGRToHSV)

	whiteMask := gocv.NewMat()
	defer whiteMask.Close()

	gocv.InRangeWithScalar(
		hsv,
		gocv.NewScalar(0, 0, 175, 0),
		gocv.NewScalar(180, 100, 255, 0),
		&whiteMask,
	)

	//debugPath := "debug_white_mask.png"
	//if ok := gocv.IMWrite(debugPath, whiteMask); !ok {
	//	log.Printf("Warning: could not write %s", debugPath)
	//} else {
	//	fmt.Printf("Debug image saved to %s\n", debugPath)
	//}

	// connect small hz gaps
	kernelW := max(15, screenW/45)
	hKernel := gocv.GetStructuringElement(
		gocv.MorphRect,
		image.Pt(kernelW, 1),
	)
	defer hKernel.Close()

	horizontal := gocv.NewMat()
	defer horizontal.Close()

	gocv.MorphologyEx(
		whiteMask,
		&horizontal,
		gocv.MorphOpen,
		hKernel,
	)

	// hz line contours
	contours := gocv.FindContours(
		horizontal,
		gocv.RetrievalExternal,
		gocv.ChainApproxSimple,
	)
	defer contours.Close()

	//debugPath = "debug_horizontal.png"
	//if ok := gocv.IMWrite(debugPath, horizontal); !ok {
	//	log.Printf("Warning: could not write %s", debugPath)
	//} else {
	//	fmt.Printf("Debug image saved to %s\n", debugPath)
	//}

	var lines []Line

	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		r := gocv.BoundingRect(contour)
		contour.Close()

		if r.Dx() < screenW/5 {
			continue
		}
		if r.Dy() > 10 || r.Dy() < 1 {
			continue
		}

		lines = append(lines, Line{
			X: r.Min.X,
			Y: r.Min.Y + roiY,
			W: r.Dx(),
			H: r.Dy(),
		})
	}

	// top/bottom line pairs
	var best *Bounds
	bestScore := math.MaxFloat64

	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines); j++ {
			a, b := lines[i], lines[j]

			top, bottom := a, b
			if top.Y > bottom.Y {
				top, bottom = bottom, top
			}

			height := bottom.Y - top.Y
			if height < screenH/50 || height > screenH/5 {
				continue
			}

			// The two lines should overlap horizontally
			left := max(top.X, bottom.X)
			right := min(top.X+top.W, bottom.X+bottom.W)
			width := right - left

			if width < screenW/5 {
				continue
			}

			shorter := min(top.W, bottom.W)
			if shorter == 0 || width < shorter*75/100 {
				continue
			}

			// The reference UI is approximately 8.5:1
			ratio := float64(width) / float64(height)
			if ratio < 5 || ratio > 12 {
				continue
			}

			lengthPenalty := math.Abs(
				float64(top.W-bottom.W) / float64(shorter),
			)
			ratioPenalty := math.Abs(ratio-8.5) / 8.5

			score := ratioPenalty + lengthPenalty

			if score < bestScore {
				bestScore = score
				best = &Bounds{
					X:      left,
					Y:      top.Y,
					Width:  width,
					Height: height,
					Right:  right,
					Bottom: bottom.Y,
				}
			}
		}
	}

	if best == nil {
		return nil, fmt.Errorf("fishing bar not found")
	}

	// pad to inside the bar
	pad := 10
	best.X = max(0, best.X+pad)
	best.Y = max(0, best.Y+pad)
	best.Right = min(screenW, best.Right-pad)
	best.Bottom = min(screenH, best.Bottom-pad)
	best.Width = best.Right - best.X
	best.Height = best.Bottom - best.Y

	return best, nil
}

func (d *FisherDetector) FindFishCoord(img image.Image, bar *Bounds) (*Coordinate, error) {
	// find pixel that is color white (255, 255, 255) in the bar region
	for y := bar.Y; y < bar.Bottom; y++ {
		for x := bar.X; x < bar.Right; x++ {
			pixel := img.At(x, y)
			b, g, r, _ := pixel.RGBA()
			if r>>8 == 255 && g>>8 == 255 && b>>8 == 255 {
				return &Coordinate{X: x, Y: y}, nil
			}
		}
	}

	return nil, fmt.Errorf("fish not found")
}

func (d *FisherDetector) FindRedBarCoord(img image.Image, bar *Bounds) (*Coordinate, error) {
	// find pixel that is color (#6E72ED) in the bar region
	for y := bar.Y; y < bar.Bottom; y++ {
		for x := bar.X; x < bar.Right; x++ {
			pixel := img.At(x, y)
			b, g, r, _ := pixel.RGBA()
			if r>>8 == 110 && g>>8 == 114 && b>>8 == 237 {
				return &Coordinate{X: x, Y: y}, nil
			}
		}
	}

	return nil, fmt.Errorf("red bar not found")
}

func (d *FisherDetector) ImageToMat(src image.Image) (gocv.Mat, error) {

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	nrgbaImg, ok := src.(*image.NRGBA)
	if !ok {
		nrgbaImg = image.NewNRGBA(bounds)
		draw.Draw(nrgbaImg, bounds, src, bounds.Min, draw.Src)
	}

	tightBytes := make([]byte, width*height*4)

	for y := 0; y < height; y++ {
		srcStart := y * nrgbaImg.Stride
		srcEnd := srcStart + (width * 4)
		dstStart := y * (width * 4)
		dstEnd := dstStart + (width * 4)

		copy(tightBytes[dstStart:dstEnd], nrgbaImg.Pix[srcStart:srcEnd])
	}

	mat, err := gocv.NewMatFromBytes(height, width, gocv.MatTypeCV8UC4, tightBytes)
	if err != nil {
		return gocv.Mat{}, err
	}

	gocv.CvtColor(mat, &mat, gocv.ColorRGBAToBGRA)

	return mat, nil
}
