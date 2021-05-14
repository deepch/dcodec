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

func NewDecoder(resize, crop, gpu string) (decoder *Decoder, err error) {
	codec := &avcodec.Codec{}
	ctx := &avcodec.Context{}
	codec = avcodec.AvcodecFindDecoderByName("h264_cuvid")
	if codec == nil {
		codec = avcodec.AvcodecFindDecoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
		if codec == nil {
			err = errors.New("AvcodecFindDecoderByName, Unsupported codec!")
			return
		} else {
			//		log.Println("============================================================Open Codec decoder soft===============================================================")
		}
	} else {

	}
	ctx = codec.AvcodecAllocContext3()
	if crop != "" {
		ctx.SetEncodeParams("crop", crop, 0)
	}
	if resize != "" {
		ctx.SetEncodeParams("resize", resize, 0)
	}
	if gpu != "" {
		ctx.SetEncodeParams("gpu", gpu, 0)
	}
	if ctx.AvcodecOpen2(codec, nil) < 0 {
		err = errors.New("AvcodecOpen2, open error")
		return
	}
	packet := avcodec.AvPacketAlloc()
	decoder = &Decoder{codec: codec, ctx: ctx, pkt: packet}
	return
}
func (self *Decoder) Decode(data []byte) (videoFrame *avutil.Frame, err error) {
	var frameFinished int
	self.pkt.AvPacketFromByteSlice(data)
	videoFrame = avutil.AvFrameAlloc()
	self.ctx.AvcodecDecodeVideo2((*avcodec.Frame)(unsafe.Pointer(videoFrame)), &frameFinished, self.pkt)
	if frameFinished > 0 {
		if err != nil {
			err = errors.New("Image Error")
		}
	} else {
		err = errors.New("No Image Frame")
	}

	self.pkt.AvPacketUnref()
	return
}
func (self *Decoder) GetPicturev5(name string, videoFrame *avutil.Frame) (img *image.YCbCr, err error) {
	img, err = avutil.GetPicturev5(videoFrame)
	return
}
func (self *Decoder) GetPicture(videoFrame *avutil.Frame, buffer *[]float32) (buf *[]float32, img *image.YCbCr, LastFrameFloat *float32, err error) {
	buf, img, LastFrameFloat, err = avutil.GetPicture(videoFrame, buffer)

	return
}
func (self *Decoder) GetPictureGray(videoFrame *avutil.Frame, buffer *[]float32) (buf *[]float32, img *image.YCbCr, LastFrameFloat *float32, err error) {
	buf, img, LastFrameFloat, err = avutil.GetPictureGray(videoFrame, buffer)

	return
}
func (self *Decoder) Close() {
	avcodec.AvcodecFreeContext(self.ctx)
	self.ctx.AvcodecClose()
}
