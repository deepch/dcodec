package dcodec

import "C"
import (
	"errors"
	"image"
	"unsafe"

	"github.com/deepch/goav/avcodec"
	"github.com/deepch/goav/avutil"
)

type Encoder struct {
	codec *avcodec.Codec
	ctx   *avcodec.Context
	pkt   *avcodec.Packet
}

func NewEncoder() (encoder *Encoder, err error) {
	codec := &avcodec.Codec{}
	ctx := &avcodec.Context{}
	codec = avcodec.AvcodecFindEncoderByName("h264_cuvid")
	if codec == nil {
		codec = avcodec.AvcodecFindEncoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
		if codec == nil {
			err = errors.New("AvcodecFindEncoderByName, Unsupported codec!")
			return
		}
	}
	ctx = codec.AvcodecAllocContext3()
	ctx.SetBitRate(0)
	ctx.SetWidth(1920)
	ctx.SetHeight(1080)
	ctx.SetTimeBase(avutil.NewRational(1, 25))
	//ctx.SetFramerate(avutil.NewRational(25, 1))
	ctx.SetGopSize(10)
	//ctx.SetMaxBFrames(1)
	ctx.SetEncodeParams("preset", "ultrafast", 0)
	ctx.SetEncodeParams("profile", "baseline", 0)
	//ctx.
	ctx.SetPixFmt(avcodec.AV_PIX_FMT_YUV420P)
	if ctx.AvcodecOpen2(codec, nil) < 0 {
		err = errors.New("AvcodecOpen2, open error")
		return
	}
	encoder = &Encoder{codec: codec, ctx: ctx}
	return
}
func (self *Encoder) Encode(img *image.YCbCr) (payload []byte, err error) {
	var frameFinished int
	var videoFrame *avutil.Frame
	packet := avcodec.AvPacketAlloc()
	videoFrame = avutil.AvFrameAlloc()
	avutil.SetPicture(videoFrame, img)
	videoFrame.SetHeight(1080)
	videoFrame.SetWidth(1920)
	videoFrame.SetFormat(avcodec.AV_PIX_FMT_YUV420P)
	self.ctx.AvcodecEncodeVideo2(packet, (*avcodec.Frame)(unsafe.Pointer(videoFrame)), &frameFinished)
	if frameFinished > 0 {
		payload = C.GoBytes(unsafe.Pointer(packet.Data()), C.int(packet.Size()))
	} else {
		err = errors.New("No Frame")
	}
	avutil.AvFrameFree(videoFrame)
	packet.AvPacketUnref()
	return
}
func (self *Encoder) Close() {
	avcodec.AvcodecFreeContext(self.ctx)
	self.ctx.AvcodecClose()
}
