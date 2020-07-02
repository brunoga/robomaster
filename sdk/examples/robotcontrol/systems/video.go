package systems

import (
	"image"
	"sync"

	"golang.org/x/image/colornames"

	"github.com/brunoga/robomaster/sdk"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type systemVideoEntity struct {
	ecs.BasicEntity
	*common.SpaceComponent
	*common.RenderComponent
}

type Video struct {
	entity  *systemVideoEntity
	client  *sdk.Client
	frameCh chan *image.NRGBA
}

func NewVideo(client *sdk.Client) *Video {
	return &Video{
		nil,
		client,
		make(chan *image.NRGBA, 1),
	}
}

func (v *Video) Add() {
	// We initialize everything we need internally.
}

func (v *Video) New(w *ecs.World) {
	spaceComponent := &common.SpaceComponent{
		Position: engo.Point{X: 0, Y: 0},
		Width:    1280,
		Height:   720,
	}

	rect := image.Rect(0, 0, 1280, 720)
	img := image.NewNRGBA(rect)
	obj := common.NewImageObject(img)

	renderComponent := &common.RenderComponent{
		Drawable: common.NewTextureSingle(obj),
	}

	v.entity = &systemVideoEntity{
		ecs.NewBasic(),
		spaceComponent,
		renderComponent,
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&v.entity.BasicEntity,
				v.entity.RenderComponent,
				v.entity.SpaceComponent)
		}
	}

	videoModule := v.client.VideoModule()
	videoModule.StartStream(v.videoHandler)
}

func (v *Video) Update(dt float32) {
	select {
	case img := <-v.frameCh:
		obj := common.NewImageObject(img)
		tex := common.NewTextureSingle(obj)
		v.entity.Drawable.Close()
		v.entity.Drawable = tex
	default:
		//do nothing
	}
}

func (v *Video) Remove(e ecs.BasicEntity) {
	if e.ID() == v.entity.ID() {
		v.entity = nil
	}
}

func (v *Video) Priority() int {
	// Makes sure video updates after other systems.
	return -10
}

func horizontalLine(img *image.NRGBA, y, x1, x2 int) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, colornames.Greenyellow)
	}
}

func verticalLine(img *image.NRGBA, x, y1, y2 int) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, colornames.Greenyellow)
	}
}

func (v *Video) videoHandler(frame *image.RGBA, wg *sync.WaitGroup) {
	frameCopy := *(*image.NRGBA)(frame)

	// Draw a simple crosshair.
	horizontalLine(&frameCopy, sdk.CameraVerticalResolutionPoints/2,
		(sdk.CameraHorizontalResolutionPoints/2)-50,
		(sdk.CameraHorizontalResolutionPoints/2)+50)
	verticalLine(&frameCopy, sdk.CameraHorizontalResolutionPoints/2,
		(sdk.CameraVerticalResolutionPoints/2)-50,
		(sdk.CameraVerticalResolutionPoints/2)+50)

	v.frameCh <- &frameCopy

	wg.Done()
}
