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

func NewEncoder(w, h, fps, bitrate int, gpu string) (encoder *Encoder, err error) {
	codec := &avcodec.Codec{}
	ctx := &avcodec.Context{}
	codec = avcodec.AvcodecFindEncoderByName("nvenc_h264")
	if codec == nil {
		codec = avcodec.AvcodecFindEncoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
		if codec == nil {
			err = errors.New("AvcodecFindEncoderByName, Unsupported codec!")
			return
		} else {
			//	log.Println("============================================================Open Codec encoder nvenc_h264===============================================================")
		}
	} else {
		//log.Println("============================================================Open Codec encoder nvenc_h264===============================================================")
	}
	ctx = codec.AvcodecAllocContext3()
	//	ctx.SetBitRate(10000000) //// size of gop * 8 == bit_rate
	//log.Fatalln(fps)
	ctx.SetWidth(w)
	ctx.SetHeight(h)
	if gpu != "" {
		ctx.SetEncodeParams("gpu", gpu, 0)
	}
	ctx.SetTimeBase(avutil.NewRational(1, fps))
	ctx.SetGopSize(fps)
	ctx.SetPixFmt(avcodec.AV_PIX_FMT_YUV420P)
	if ctx.AvcodecOpen2(codec, nil) < 0 {
		//log.Println("============================================================Open Codec encoder ERROR===============================================================")
		avcodec.AvcodecFreeContext(ctx)
		codec = avcodec.AvcodecFindEncoderByName("libx264")
		if codec != nil {
			//log.Println("============================================================SWITCH Codec encoder libx264===============================================================")
			ctx = codec.AvcodecAllocContext3()
			ctx.SetBitRate(0)
			ctx.SetWidth(w)
			ctx.SetHeight(h)
			ctx.SetTimeBase(avutil.NewRational(1, 1))
			ctx.SetGopSize(1)
			ctx.SetEncodeParams("preset", "ultrafast", 0)
			ctx.SetEncodeParams("profile", "baseline", 0)
			ctx.SetThreadCount(20)
			ctx.SetPixFmt(avcodec.AV_PIX_FMT_YUV420P)
			if ctx.AvcodecOpen2(codec, nil) < 0 {
				err = errors.New("AvcodecOpen2, open error")
				return
			}
		} else {
			err = errors.New("AvcodecOpen2, open error")
			return
		}
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
	//log.Println("============================================================!!!!!!!!!!!!!!!!!!!!!!!!!!!! Close Codecs !!!!!!!!!!!!!!!!===============================================================")
	avcodec.AvcodecFreeContext(self.ctx)
	self.ctx.AvcodecClose()
	//avcodec.Codec = nul
}
