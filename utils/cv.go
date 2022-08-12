package utils

import (
	"bytes"
	"encoding/base64"
	"image"

	"gocv.io/x/gocv"
)

type Vecb []uint8

func GetVecbAt(m gocv.Mat, row int, col int) Vecb {
	ch := m.Channels()
	v := make(Vecb, ch)

	for c := 0; c < ch; c++ {
		v[c] = m.GetUCharAt(row, col*ch+c)
	}

	return v
}

func (v Vecb) SetVecbAt(m gocv.Mat, row int, col int) {
	ch := m.Channels()

	for c := 0; c < ch; c++ {
		m.SetUCharAt(row, col*ch+c, v[c])
	}
}

func getPixel(mat gocv.Mat, row, col int) uint8 {
	v := GetVecbAt(mat, row, col)
	return v[0]
}

func setPixel(mat gocv.Mat, row, col int, val uint8) {
	v := GetVecbAt(mat, row, col)
	for idx, _ := range v {
		v[idx] = val
	}
	v.SetVecbAt(mat, row, col)
}

func mat2base64(mat gocv.Mat) string {
	// 图像 mat 转 base64
	imb, _ := gocv.IMEncode(gocv.JPEGFileExt, mat)
	defer mat.Close()
	return base64.StdEncoding.EncodeToString(imb.GetBytes())
}

func binaryImage(b []byte) gocv.Mat {
	//mat := gocv.IMRead("./c.jpg", gocv.IMReadColor)
	mat, _ := gocv.IMDecode(b, gocv.IMReadColor)
	defer mat.Close()

	gocv.CvtColor(mat, &mat, gocv.ColorBGRToGray)

	threshold := gocv.NewMat()
	defer threshold.Close()

	gocv.AdaptiveThreshold(mat, &threshold, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, 21, 1)

	bilateral := gocv.NewMat()
	//defer bilateral.Close()
	gocv.BilateralFilter(threshold, &bilateral, -1, 0.3, 10)

	interferenceLine(bilateral)

	// w := gocv.NewWindow("w")
	// w.IMShow(bilateral)
	// gocv.WaitKey(0)

	// print(mat2base64(bilateral))

	return bilateral
}

func interferenceLine(mat gocv.Mat) {
	// img, _ := mat.ToImage()
	// bounds := img.Bounds()

	// newRgba := image.NewRGBA(bounds)

	// size := mat.Size()
	// h := size[0]
	// w := size[1]

	for y := 1; y < mat.Rows()-1; y++ {
		for x := 1; x < mat.Cols()-1; x++ {
			count := 0
			if getPixel(mat, y, x-1) > 245 {
				count++
			}
			if getPixel(mat, y-1, x) > 245 {
				count++
			}
			if getPixel(mat, y, x+1) > 245 {
				count++
			}
			if getPixel(mat, y+1, x) > 245 {
				count++
			}

			if count > 2 {
				setPixel(mat, y, x, 255)
			}
		}
	}
}

func cvtImageToMat(img image.Image) (gocv.Mat, error) {
	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()
	bytes := make([]byte, 0, x*y*3)

	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8), byte(g>>8), byte(r>>8))
		}
	}
	return gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
}

/*
	Parameters:
		mat: gocv.Mat，原图像
		format: string，要转换的类型格式，比如.png .jpeg
	Returns:
		img: image.Image 转换之后的图像
		err: error
*/
func cvtMatToImage(mat gocv.Mat, format string) (img image.Image, err error) {
	// 把mat转成字节，指定图片格式format
	buf, err := gocv.IMEncode(gocv.FileExt(format), mat)
	if err != nil {
		return nil, err
	}
	// 根据图片的字节，创建出reader
	reader := bytes.NewReader(buf.GetBytes())
	// 解码，将字节解码成图片
	dest, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return dest, err
}

// 二值化并降噪图片，输出图片base64
func ImageProcess(ib []byte) string {
	return mat2base64(binaryImage(ib))
}
