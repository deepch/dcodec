package dcodec

import "C"
import (
	"errors"
	"image"
	"unsafe"

	"github.com/deepch/goav/avcodec"
	"github.com/deepch/goav/avutil"
)

type Decoder struct {
	codec *avcodec.Codec
	ctx   *avcodec.Context
	pkt   *avcodec.Packet
}

func NewDecoder() (decoder *Decoder, err error) {
	codec := &avcodec.Codec{}
	ctx := &avcodec.Context{}
	codec = avcodec.AvcodecFindDecoderByName("h264_cuvid")
	if codec == nil {
		codec = avcodec.AvcodecFindDecoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
		if codec == nil {
			err = errors.New("AvcodecFindDecoderByName, Unsupported codec!")
			return
		}
	}
	ctx = codec.AvcodecAllocContext3()
	if ctx.AvcodecOpen2(codec, nil) < 0 {
		err = errors.New("AvcodecOpen2, open error")
		return
	}
	packet := avcodec.AvPacketAlloc()
	decoder = &Decoder{codec: codec, ctx: ctx, pkt: packet}
	return
}
func (self *Decoder) Decode(data []byte) (img *image.YCbCr, err error) {
	var frameFinished int
	var videoFrame *avutil.Frame
	self.pkt.AvPacketFromByteSlice(data)
	videoFrame = avutil.AvFrameAlloc()
	self.ctx.AvcodecDecodeVideo2((*avcodec.Frame)(unsafe.Pointer(videoFrame)), &frameFinished, self.pkt)
	if frameFinished > 0 {
		img, err = avutil.GetPicture(videoFrame)
		if err != nil {
			err = errors.New("Image Error")
		}
	} else {
		err = errors.New("No Image Frame")
	}
	avutil.AvFrameFree(videoFrame)
	self.pkt.AvPacketUnref()
	return
}
func (self *Decoder) Close() {
	avcodec.AvPacketFree(self.pkt)
	avcodec.AvcodecFreeContext(self.ctx)
	self.ctx.AvcodecClose()
}
