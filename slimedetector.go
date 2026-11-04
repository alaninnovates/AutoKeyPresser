package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"sort"

	ort "github.com/yalue/onnxruntime_go"
)

const (
	InputW = 640
	InputH = 640

	NumClasses    = 3
	NumCandidates = 8400

	ConfidenceThreshold = 0.40
	IOUThreshold        = 0.45
)

var classNames = []string{
	"green_slime",
	"blue_slime",
	"pink_slime",
}

type Detection struct {
	ClassID    int
	Confidence float32
	X1, Y1     float32
	X2, Y2     float32
}

type Letterbox struct {
	Scale float32
	PadX  float32
	PadY  float32
}

type SlimeDetector struct {
	session *ort.AdvancedSession
	input   *ort.Tensor[float32]
	output  *ort.Tensor[float32]
}

func NewSlimeDetector(modelPath string) (*SlimeDetector, error) {
	ort.SetSharedLibraryPath(
		"/Users/alan/Documents/projects/AutoKeyPresser/slime-detector-test/onnxruntime-osx-arm64-1.31.1/lib/libonnxruntime.dylib",
	)

	if err := ort.InitializeEnvironment(); err != nil {
		return nil, err
	}

	input, err := ort.NewEmptyTensor[float32](
		ort.NewShape(1, 3, InputH, InputW),
	)
	if err != nil {
		return nil, err
	}

	output, err := ort.NewEmptyTensor[float32](
		ort.NewShape(1, 4+NumClasses, NumCandidates),
	)
	if err != nil {
		input.Destroy()
		return nil, err
	}

	session, err := ort.NewAdvancedSession(
		modelPath,
		[]string{"images"},
		[]string{"output0"},
		[]ort.Value{input},
		[]ort.Value{output},
		nil,
	)
	if err != nil {
		input.Destroy()
		output.Destroy()
		return nil, err
	}

	return &SlimeDetector{
		session: session,
		input:   input,
		output:  output,
	}, nil
}

func (d *SlimeDetector) Close() {
	if d.session != nil {
		d.session.Destroy()
	}
	if d.input != nil {
		d.input.Destroy()
	}
	if d.output != nil {
		d.output.Destroy()
	}

	ort.DestroyEnvironment()
}

func (d *SlimeDetector) ClassName(id int) string {
	if id >= 0 && id < len(classNames) {
		return classNames[id]
	}
	return fmt.Sprintf("class_%d", id)
}

func (d *SlimeDetector) Detect(img image.Image) ([]Detection, error) {
	inputData, lb := preprocess(img)

	copy(d.input.GetData(), inputData)

	if err := d.session.Run(); err != nil {
		return nil, err
	}

	bounds := img.Bounds()

	return decodeYOLO(
		d.output.GetData(),
		bounds.Dx(),
		bounds.Dy(),
		lb,
	), nil
}

// letterbox, normalize, convert to NCHW RGB tensor
func preprocess(img image.Image) ([]float32, Letterbox) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	scale := math.Min(
		float64(InputW)/float64(w),
		float64(InputH)/float64(h),
	)

	newW := int(math.Round(float64(w) * scale))
	newH := int(math.Round(float64(h) * scale))

	padX := (InputW - newW) / 2
	padY := (InputH - newH) / 2

	data := make([]float32, 3*InputW*InputH)

	plane := InputW * InputH

	for i := range data {
		data[i] = 114.0 / 255.0
	}

	for y := 0; y < newH; y++ {
		srcY := int(float64(y) / scale)

		for x := 0; x < newW; x++ {
			srcX := int(float64(x) / scale)

			r, g, b, _ := img.At(
				bounds.Min.X+srcX,
				bounds.Min.Y+srcY,
			).RGBA()

			dst := (padY+y)*InputW + padX + x

			data[dst] = float32(r>>8) / 255
			data[plane+dst] = float32(g>>8) / 255
			data[2*plane+dst] = float32(b>>8) / 255
		}
	}

	return data, Letterbox{
		Scale: float32(scale),
		PadX:  float32(padX),
		PadY:  float32(padY),
	}
}

func decodeYOLO(
	output []float32,
	imgW, imgH int,
	lb Letterbox,
) []Detection {
	features := 4 + NumClasses

	if len(output) != features*NumCandidates {
		log.Printf("Unexpected output length: %d", len(output))
		return nil
	}

	at := func(channel, candidate int) float32 {
		return output[channel*NumCandidates+candidate]
	}

	results := make([]Detection, 0)

	for i := 0; i < NumCandidates; i++ {
		classID := 0
		score := float32(0)

		for c := 0; c < NumClasses; c++ {
			s := at(4+c, i)
			if s > score {
				score = s
				classID = c
			}
		}

		if score < ConfidenceThreshold {
			continue
		}

		cx := at(0, i)
		cy := at(1, i)
		bw := at(2, i)
		bh := at(3, i)

		x1 := (cx - bw/2 - lb.PadX) / lb.Scale
		y1 := (cy - bh/2 - lb.PadY) / lb.Scale
		x2 := (cx + bw/2 - lb.PadX) / lb.Scale
		y2 := (cy + bh/2 - lb.PadY) / lb.Scale

		x1 = clamp(x1, 0, float32(imgW))
		y1 = clamp(y1, 0, float32(imgH))
		x2 = clamp(x2, 0, float32(imgW))
		y2 = clamp(y2, 0, float32(imgH))

		if x2 <= x1 || y2 <= y1 {
			continue
		}

		results = append(results, Detection{
			ClassID:    classID,
			Confidence: score,
			X1:         x1,
			Y1:         y1,
			X2:         x2,
			Y2:         y2,
		})
	}

	return nms(results, IOUThreshold)
}

func clamp(v, lo, hi float32) float32 {
	return float32(math.Max(
		float64(lo),
		math.Min(float64(v), float64(hi)),
	))
}

func iou(a, b Detection) float32 {
	x1 := max(a.X1, b.X1)
	y1 := max(a.Y1, b.Y1)
	x2 := min(a.X2, b.X2)
	y2 := min(a.Y2, b.Y2)

	w := max(0, x2-x1)
	h := max(0, y2-y1)

	intersection := w * h
	areaA := (a.X2 - a.X1) * (a.Y2 - a.Y1)
	areaB := (b.X2 - b.X1) * (b.Y2 - b.Y1)

	union := areaA + areaB - intersection

	if union <= 0 {
		return 0
	}

	return intersection / union
}

func nms(input []Detection, threshold float32) []Detection {
	sort.Slice(input, func(i, j int) bool {
		return input[i].Confidence > input[j].Confidence
	})

	kept := make([]Detection, 0)
	suppressed := make([]bool, len(input))

	for i := range input {
		if suppressed[i] {
			continue
		}

		kept = append(kept, input[i])

		for j := i + 1; j < len(input); j++ {
			if suppressed[j] ||
				input[i].ClassID != input[j].ClassID {
				continue
			}

			if iou(input[i], input[j]) > threshold {
				suppressed[j] = true
			}
		}
	}

	return kept
}
