package mode

import (
	"image/color"

	"gocv.io/x/gocv"
)

type Person struct {
	hogDescriptor gocv.HOGDescriptor
}

func NewPerson() *Person {
	hogDescriptor := gocv.NewHOGDescriptor()
	hogDescriptor.SetSVMDetector(gocv.HOGDefaultPeopleDetector())

	return &Person{
		hogDescriptor,
	}
}

func (p *Person) FindPeople(frame *gocv.Mat) {
	// TODO(bga): Use DetectMultiScaleWithParams to tweak things.
	rects := p.hogDescriptor.DetectMultiScale(*frame)

	// TODO(bga): We need to coalesce overlapping rects to get a single one
	// that will cover a person.
	for _, r := range rects {
		gocv.Rectangle(frame, r, color.RGBA{0, 255, 0, 255}, 3)
	}
}
