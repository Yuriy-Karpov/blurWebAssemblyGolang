package main

import (
	"syscall/js"
)

// Abs returns the absolute value of x.
// Для int нету функции Abs, только для float, поэтому напишим свою
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func gaussMatrix(src []byte, width int, height int, radius int) []byte {

	tableSize := 2*radius + 1
	widthMax := width - 1
	heightMax := height - 1

	// fmt.Println(tableSize)
	// fmt.Println(widthMax)
	// fmt.Println(heightMax)

	// target := make([]byte, width*height)
	for y := 0; y < height; y++ {
		indexHorizontLine := y * width
		for indexPixelInLine := 0; indexPixelInLine < width; indexPixelInLine++ {
			indexPixel := indexPixelInLine + indexHorizontLine
			sumPixel := blurCore(src, indexPixel, indexPixelInLine, radius, widthMax, 1)
			middleColorPixel := sumPixel / tableSize
			src[indexPixel] = byte(middleColorPixel)
		}
	}
	offset := width
	for x := 0; x < width; x++ {
		for indexVertical := 0; indexVertical < height; indexVertical++ {
			indexPixel := indexVertical*offset + x
			sumPixel := blurCore(src, indexPixel, indexVertical, radius, heightMax, offset)
			middleColorPixel := sumPixel / tableSize
			src[indexPixel] = byte(middleColorPixel)
		}
	}
	return src
}

func blurCore(src []byte, inIndex int, indexInLine int, radius int, max int, offset int) int {
	var tmpPixel int = 0
	for r := -radius; r <= radius; r++ {
		targetPixel := indexInLine + r
		if targetPixel <= 0 {
			t := r
			for indexInLine+t < 0 {
				t++
			}
			tmpPixel += int(src[inIndex+t*offset])
			continue
		}
		if targetPixel > max {
			t := r
			for indexInLine+t*offset > max {
				t--
			}
			tmpPixel += int(src[inIndex+t*offset])
			continue
		}

		var pxIndex int
		if r < 0 {
			pxIndex = inIndex - (Abs(r) * offset)
		} else if r > 0 {
			pxIndex = inIndex + (Abs(r) * offset)
		} else {
			pxIndex = inIndex
		}
		tmpPixel += int(src[pxIndex])
	}
	return tmpPixel
}

func calcAsm(value js.Value, args []js.Value) interface{} {
	callback := args[len(args)-1:][0]
	src := args[0]
	target := args[1]
	width := args[2].Int()
	height := args[3].Int()
	radius := args[4].Int()

	inBuf := make([]byte, width*height)
	js.CopyBytesToGo(inBuf, src)
	outBuf := gaussMatrix(inBuf, width, height, radius)
	js.CopyBytesToJS(target, outBuf)
	callback.Invoke(js.Null(), target)
	return nil
}

func registerCallbacks() {
	js.Global().Set("calcAsm", js.FuncOf(calcAsm))
}

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}
