package main

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/screenshot"
	"gocv.io/x/gocv"
)

type FisherService struct {
	ctx            context.Context
	cancel         context.CancelFunc
	fisherDetector FisherDetector
}

func NewFisherService() *FisherService {
	fisherDetector := NewFisherDetector()

	return &FisherService{
		fisherDetector: *fisherDetector,
	}
}

func FastCaptureRegion(x, y, w, h int) (gocv.Mat, error) {
	bounds := image.Rect(x, y, x+w, y+h)
	img, err := screenshot.Capture(bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy())
	if err != nil {
		return gocv.Mat{}, fmt.Errorf("screenshot capture failed: %w", err)
	}
	mat, err := gocv.NewMatFromBytes(h, w, gocv.MatTypeCV8UC4, img.Pix)
	if err != nil {
		return gocv.Mat{}, fmt.Errorf("failed to load bytes into mat: %w", err)
	}
	gocv.CvtColor(mat, &mat, gocv.ColorRGBAToBGRA)
	return mat, nil
}

func (s *FisherService) start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	go func() {
		for {
			name, _, pid := robotgo.GetActiveApp()
			if name != "Roblox" {
				fmt.Println("Active window is not Roblox")
				break
			}
			x, y, w, h := robotgo.GetBounds(pid)
			robotgo.MouseDown("left")
			time.Sleep(600 * time.Millisecond)
			robotgo.MouseUp("left")
			time.Sleep(2500 * time.Millisecond)

			select {
			case <-s.ctx.Done():
				return
			case <-time.After(1 * time.Millisecond):
				foundBarOnce := false
			OuterLoop:
				for {
					select {
					case <-s.ctx.Done():
						return
					case <-time.After(1 * time.Millisecond):
						//img, err := robotgo.CaptureImg(x, y, w, h)
						//mat, _ := s.fisherDetector.ImageToMat(img)
						mat, err := FastCaptureRegion(x, y, w, h)
						if err != nil {
							fmt.Println("Capture error:", err)
							continue
						}
						img, err := mat.ToImage()
						if err != nil {
							fmt.Println("Mat to Image error:", err)
							continue
						}
						minigameBounds, err := s.fisherDetector.FindFishingMinigameBounds(mat)
						if err != nil {
							//fmt.Println("Minigame detection error:", err)
							if foundBarOnce {
								//fmt.Println("Breaking because bar found - select loop")
								break OuterLoop
							}
							continue
						}
						mat.Close()
						//fmt.Println("found bar")
						foundBarOnce = true
						fishCoord, err := s.fisherDetector.FindFishCoord(img, minigameBounds)
						if err != nil {
							fmt.Println("Fish detection error:", err)
							continue
						}
						redBarCoord, err := s.fisherDetector.FindRedBarCoord(img, minigameBounds)
						if err != nil {
							fmt.Println("Bar detection error:", err)
							continue
						}
						if redBarCoord.X < fishCoord.X {
							robotgo.MouseDown("left")
						} else if redBarCoord.X > fishCoord.X {
							robotgo.MouseUp("left")
						}
					}
				}
			}
			time.Sleep(100 * time.Millisecond)
			robotgo.MouseUp("left")
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

func (s *FisherService) stop() {
	s.cancel()
}
