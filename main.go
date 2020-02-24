package main

import (
	"fmt"
	"syscall/js"
	"time"
)

// Abs returns the absolute value of x.
// Для int нету функции Abs, только для float, поэтому напишим свою
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type unitPixel struct {
	target int
	pixel  byte
}

func gaussMatrix(src []byte, width int, height int, radius int) []byte {

	tableSize := 2*radius + 1
	widthMax := width - 1
	heightMax := height - 1
	amountPixel := width * height
	// fmt.Println("amountPixel", amountPixel)
	jobs := make(chan unitPixel)
	// jobs := make(chan int, amountPixel)
	// done := make(chan bool)
	// target := make([]byte, width*height)
	offsetForY := 1
	// var acmLoopCount int = 0
	for y := 0; y < height; y++ {
		indexHorizontLine := y * width
		for indexPixelInLine := 0; indexPixelInLine < width; indexPixelInLine++ {
			indexPixel := indexPixelInLine + indexHorizontLine
			// acmLoopCount++
			go blurCore(
				src,
				indexPixel,
				indexPixelInLine,
				radius,
				widthMax,
				offsetForY,
				tableSize,
				jobs)
		}
	}
	for i := 0; i < amountPixel; i++ {
		a := <-jobs
		src[a.target] = a.pixel
	}
	offsetForX := width
	for x := 0; x < width; x++ {
		for indexVertical := 0; indexVertical < height; indexVertical++ {
			indexPixel := indexVertical*offsetForX + x
			go blurCore(
				src,
				indexPixel,
				indexVertical,
				radius,
				heightMax,
				offsetForX,
				tableSize,
				jobs)
		}
	}
	for i := 0; i < amountPixel; i++ {
		a := <-jobs
		src[a.target] = a.pixel
	}
	return src
}

func blurCore(
	src []byte,
	inIndex int,
	indexInLine int,
	radius int,
	max int,
	offset int,
	tableSize int,
	jobs chan unitPixel) {
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
	tmpPixelBit := byte(tmpPixel / tableSize)
	jobs <- unitPixel{inIndex, tmpPixelBit}
	// fmt.Println("run", i)
}

func calcAsm(value js.Value, args []js.Value) interface{} {
	callback := args[len(args)-1:][0]
	src := args[0]
	target := args[1]
	width := args[2].Int()
	height := args[3].Int()
	radius := args[4].Int()

	inBuf := make([]byte, width*height)
	t1 := time.Now()
	js.CopyBytesToGo(inBuf, src)
	t := time.Now()
	outBuf := gaussMatrix(inBuf, width, height, radius)
	fmt.Println("timerGo:", time.Since(t))
	js.CopyBytesToJS(target, outBuf)
	fmt.Println("timerGoAll:", time.Since(t1))
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
