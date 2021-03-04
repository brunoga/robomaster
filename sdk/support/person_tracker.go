package support

import (
	"gocv.io/x/gocv"
)

type PersonTracker struct {
	hogDescriptor gocv.HOGDescriptor
}

func NewPersonTracker() *PersonTracker {
	hogDescriptor := gocv.NewHOGDescriptor()
	hogDescriptor.SetSVMDetector(gocv.HOGDefaultPeopleDetector())

	return &PersonTracker{
		hogDescriptor,
	}
}

func (p *PersonTracker) FindPeople(frame *gocv.Mat) {
	p.hogDescriptor.DetectMultiScaleWithParams(frame)

}
