package systems

import (
	"image"
	"sync"
	"unsafe"

	"git.bug-br.org.br/bga/robomasters1/app/example/entities"
	"git.bug-br.org.br/bga/robomasters1/app/video"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Video struct {
	videoEntity      *entities.Video
	elapsed          float32
	frameCh          chan *image.NRGBA
	dataHandlerIndex int
	wg               *sync.WaitGroup
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

	appVideo, err := video.New()
	if err != nil {
		panic(err)
	}

	index, err := appVideo.AddDataHandler(v.DataHandler)
	if err != nil {
		panic(err)
	}

	v.dataHandlerIndex = index

	appVideo.StartVideo()
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

func (v *Video) DataHandler(data []byte, wg *sync.WaitGroup) {
	// Create an image out of the data byte slice.
	img := image.NewNRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{1280, 720},
		},
	)
	img.Pix = NRGBA(data)

	v.frameCh <- img

	wg.Done()
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
