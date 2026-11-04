package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-vgo/robotgo"
)

type SlimeService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	detector SlimeDetector
}

func NewSlimeService() *SlimeService {
	detector, err := NewSlimeDetector("slime_model.onnx")
	if err != nil {
		fmt.Println("Failed to create SlimeDetector:", err)
		return nil
	}

	return &SlimeService{
		detector: *detector,
	}
}

func (s *SlimeService) start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	go func() {
		logicalW, logicalH := robotgo.GetScreenSize()
		physicalW, physicalH := robotgo.GetScaleSize(robotgo.GetMainId())
		scaleX := float64(logicalW) / float64(physicalW)
		scaleY := float64(logicalH) / float64(physicalH)
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(1 * time.Millisecond):
				name, _, pid := robotgo.GetActiveApp()
				if name != "Roblox" {
					fmt.Println("Active window is not Roblox")
					break
				}
				x, y, w, h := robotgo.GetBounds(pid)
				windowCx, windowCy := float64(x+w/2)*scaleX, float64(y+h/2)*scaleY
				//fmt.Printf("Roblox PID: %d - Bounds: x=%d, y=%d, w=%d, h=%d\n", pid, x, y, w, h)
				img, _ := robotgo.CaptureImg(x, y, w, h)
				detections, err := s.detector.Detect(img)
				if err != nil {
					fmt.Println("Detection error:", err)
					break
				}

				//buf := new(bytes.Buffer)
				//_ = jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})
				//imgMat, _ := gocv.IMDecode(buf.Bytes(), gocv.IMReadColor)
				//defer imgMat.Close()
				if len(detections) == 0 {
					//fmt.Println("No slimes detected")
					break
				}
				closestSlimeCx, closestSlimeCy := float64(0), float64(0)
				closestDistance := math.MaxFloat64
				for _, detection := range detections {
					if detection.Confidence < 0.5 {
						continue
					}

					//rect := image.Rect(
					//	int(detection.X1),
					//	int(detection.Y1),
					//	int(detection.X2),
					//	int(detection.Y2),
					//)
					//
					//gocv.Rectangle(
					//	&imgMat,
					//	rect,
					//	color.RGBA{0, 255, 0, 0},
					//	2,
					//)
					//fmt.Printf("Detected %s with confidence %.2f at (%.2f, %.2f, %.2f, %.2f)\n",
					//	classNames[detection.ClassID], detection.Confidence,
					//	detection.X1, detection.Y1, detection.X2, detection.Y2)
					cx := float64((detection.X1 + detection.X2) / 2)
					cy := float64((detection.Y1 + detection.Y2) / 2)
					// only click in the bounds of roblox MINUS a padding of 150 pixels to avoid clicking on the border. account for scale in the OOB check
					padding := float64(150)
					if cx*scaleX+float64(x) < float64(x)+padding || cx*scaleX+float64(x) > float64(x+w)-padding ||
						cy*scaleY+float64(y) < float64(y)+padding || cy*scaleY+float64(y) > float64(y+h)-padding {
						//fmt.Println("Detected slime OOB")
						continue
					}
					distance := (cx-windowCx)*(cx-windowCx) + (cy-windowCy)*(cy-windowCy)
					if distance < closestDistance {
						//fmt.Printf("Detected slime is closer: %.2f < %.2f\n", distance, closestDistance)
						closestDistance = distance
						closestSlimeCx = cx
						closestSlimeCy = cy
					}
				}
				robotgo.Click("left")
				if closestDistance == 999999 {
					//fmt.Println("No slimes detected")
					break
				}
				robotgo.Move(int(closestSlimeCx*scaleX+float64(x)), int(closestSlimeCy*scaleY+float64(y)))
				//mousex, mousey := robotgo.Location()
				//print("Mouse position: ", mousex, mousey, "\n")
				//print("Detected slimes, saving screenshot to screenshot.png\n")
				//gocv.IMWrite("screenshot.png", imgMat)
			}
		}
	}()
}

func (s *SlimeService) stop() {
	s.cancel()
}
