package video

import "image"

// Handler is a function that is called when a new frame is available. The frame
// is not guaranteed to be valid after the function returns so if you need to
// keep it to do something after that, make a copy.
type Handler func(frame *image.RGBA)
