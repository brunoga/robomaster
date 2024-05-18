package mode

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// ColorObject locates objects in a frame based on their color.
type ColorObject struct {
	hsvLower  gocv.Scalar
	hsvUpper  gocv.Scalar
	minRadius float32
}

// NewColorObject returns a new ColorObject instance that locates objects with
// a color that falls between the lower bound HSV values (hl, sl, vl) and the
// upper bound HSV values (hu, su, vu), only considering objects with a radius
// greater than minRadius.
func NewColorObject(hl, sl, vl, hu, su, vu float64,
	minRadius float32) *ColorObject {
	return &ColorObject{
		gocv.NewScalar(hl, sl, vl, 0),
		gocv.NewScalar(hu, su, vu, 0),
		minRadius,
	}
}

// FindLargestObject locates the largest object that satisfies our parameters in
// a frame and returns its center position (x and y) and the radius of a circle
// that fully encloses it.
func (c *ColorObject) FindLargestObject(
	frame *gocv.Mat) (float32, float32, float32, error) {
	// Create a scratch frame with the same size and type as the original frame.
	scratch := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), frame.Type())
	defer scratch.Close()

	// Copy frame to scratch while applying gaussian blur to reduce image noise.
	gocv.GaussianBlur(*frame, &scratch, image.Point{X: 11, Y: 11}, 0,
		0, gocv.BorderDefault)

	// Convert scratch from BGR to HSV so we can apply our lower and upper bound
	// filters.
	gocv.CvtColor(scratch, &scratch, gocv.ColorBGRToHSV)

	// Try to filter out everything but our colored ball.
	gocv.InRangeWithScalar(scratch, c.hsvLower, c.hsvUpper, &scratch)

	// Erode then dilate the image to better approximate our ball shape.
	gocv.ErodeWithParams(scratch, &scratch, gocv.NewMat(),
		image.Point{X: -1, Y: -1}, 2, int(gocv.BorderDefault))
	gocv.Dilate(scratch, &scratch, gocv.NewMat())

	// Find the contours of anything that is left in the scratch image.
	contours := gocv.FindContours(scratch, gocv.RetrievalExternal,
		gocv.ChainApproxSimple)

	if contours.Size() > 0 {
		// We found at least one object. Find the biggest one.
		biggestContour := findBiggestAreaContour(contours)

		// Get the center position and radius of the minimum enclosing circle
		// that contains the object with the largest area.
		x, y, radius := gocv.MinEnclosingCircle(biggestContour)
		if radius >= c.minRadius {
			// Return coordinates and radius of what is hopefully our ball.
			return x, y, radius, nil
		}
	}

	// No object found.
	return -1, -1, -1, fmt.Errorf("no suitable object found")
}

func findBiggestAreaContour(contours gocv.PointsVector) gocv.PointVector {
	maxArea := 0.0
	maxIdx := -1
	for i := 0; i < contours.Size(); i++ {
		area := gocv.ContourArea(contours.At(i))
		if area > maxArea {
			maxArea = area
			maxIdx = i
		}
	}

	return contours.At(maxIdx)
}
