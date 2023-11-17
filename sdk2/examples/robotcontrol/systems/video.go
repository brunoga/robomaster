package systems

import (
	"image"
	"unsafe"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/entities"
	"github.com/brunoga/robomaster/sdk2/module/camera"
	"github.com/brunoga/unitybridge/support/token"
)

type Video struct {
	videoEntity      *entities.Video
	frameCh          chan *image.NRGBA
	dataHandlerToken token.Token
	C                *sdk2.Client
}

func (v *Video) New(w *ecs.World) {
	rect := image.Rect(0, 0, 1280, 720)
	v.videoEntity = &entities.Video{BasicEntity: ecs.NewBasic()}
	v.videoEntity.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: 0},
		Width:    1280,
		Height:   720,
	}

	img := image.NewNRGBA(rect)

	obj := common.NewImageObject(img)

	v.videoEntity.RenderComponent = common.RenderComponent{
		Drawable: common.NewTextureSingle(obj),
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&v.videoEntity.BasicEntity,
				&v.videoEntity.RenderComponent,
				&v.videoEntity.SpaceComponent)
		}
	}

	v.frameCh = make(chan *image.NRGBA, 30)

	camera := v.C.Camera()

	index, err := camera.AddVideoCallback(v.DataHandler)
	if err != nil {
		panic(err)
	}

	v.dataHandlerToken = index
}

func (v *Video) Add() {}

func (v *Video) Remove(basic ecs.BasicEntity) {}

func (v *Video) Update(dt float32) {
	select {
	case img := <-v.frameCh:
		obj := common.NewImageObject(img)
		tex := common.NewTextureSingle(obj)
		v.videoEntity.Drawable.Close()
		v.videoEntity.Drawable = tex
	default:
		//do nothing
	}
}

func (v *Video) DataHandler(frame *camera.RGB) {
	// Create an image out of the data byte slice.
	img := common.ImageToNRGBA(frame, 1280, 720)

	v.frameCh <- img
}

func NRGBA(rgbData []byte) []byte {
	numPixels := len(rgbData) / 3

	nrgbaData := make([]byte, numPixels*4)

	intNRGBAData := *(*[]uint32)(unsafe.Pointer(&nrgbaData))
	intNRGBAData = intNRGBAData[:len(nrgbaData)/4]

	for i, j := 0, 0; i < len(rgbData); i, j = i+3, j+1 {
		intRGB := (*(*uint32)(unsafe.Pointer(&rgbData[i]))) |
			(0b11111111 << 24)
		intNRGBAData[j] = intRGB
	}

	return nrgbaData
}
