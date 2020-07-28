package main

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/num"
	"github.com/gonum/matrix/mat64"
	"github.com/sbinet/go-python"
	"io/ioutil"
)

const (
	window_size=10
	degree=2
) 
func init() {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	//m := python.PyImport_ImportModule("sys")
	//if m == nil {
	//	fmt.Println("import error")
	//	return
	//}
	//path := m.GetAttrString("path")
	//if path == nil {
	//	fmt.Println("get path error")
	//	return
	//}
	////加入当前目录，空串表示当前目录
	//currentDir := python.PyString_FromString("")
	//python.PyList_Insert(path, 0, currentDir)

	//size := python.PyList_GET_SIZE(path)
	//for i := 0; i < size; i++ {
	//	item := python.PyList_GET_ITEM(path, i)
	//	s := python.PyString_AsString(item)
	//	fmt.Println(s)
	//}

	//m = python.PyImport_ImportModule("wy2")
	//if m == nil {
	//	fmt.Println("import error")
	//	return
	//}
	//touchBaidu := m.GetAttrString("touch_baidu")
	//if touchBaidu == nil {
	//	fmt.Println("get touch_baidu error")
	//	return
	//}
	//res := touchBaidu.CallFunction()
	//if res == nil {
	//	fmt.Println("call touch_baidu error")
	//	return
	//}
	//statusCode := res.GetAttrString("status_code")
	//content := res.GetAttrString("content")
	//fmt.Println(python.PyInt_AS_LONG(statusCode))
	//fmt.Println(python.PyString_AS_STRING(content))
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println("read file err:", err.Error())
		return
	}
	cracks := make([]*MonitorPointCrack, 0)
	err = json.Unmarshal(data, &cracks)
	if err != nil {
		fmt.Println(err)
		return
	}
	times := make([]int, 0)
	vals := make([]float32, 0)
	for i, v := range cracks {
		vals = append(vals, v.Value)
		times = append(times, i)
	}
	//step1 for at least 2 times!
	vals1 := make([]float64, 0)
	times1 := make([]int, 0)
	for j := 0; j < 1; j++ {
		for i := 0; i < 3; i++ {
			vals1 = append(vals1, float64(vals[i]))
			times1 = append(times1, times[i])
		}
		for i := 3; i < len(vals); i++ {
			if vals[i]*1.5 < (vals[i-3]+vals[i-2]+vals[i-1])/3 {
				continue
			}
			vals1 = append(vals1, float64(vals[i]))
			times1 = append(times1, times[i])
		}
	}
	//step2 for at least 3 times!
	round1 := Moving_average(vals1, window_size)
	round2 := Moving_average(round1, window_size)
	round3 := Moving_average(round2, window_size)
	//step3
	//以天为单位[0, 1, 1.2, 2 ,...]
	times2 := make([]float64, 0)
	for i := 0; i < len(times1); i++ {
		times2 = append(times2, float64(times[i])-float64(times[0])/float64(60)/float64(60)/float64(24))
	}
	a := Vandermonde(round3, degree)
	b := mat64.NewDense(len(times2), 1,times2)
	c := mat64.NewDense(degree+1, 1, nil)
	qr := new(mat64.QR)
	qr.Factorize(a)
	err = c.SolveQR(qr, false, b)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("%.3f\n", mat64.Formatted(c))
	}
	//step4
}

func Moving_average(internal []float64,window_size int) []float64{
	window:=make([]float64,window_size)
	for i := 0; i < window_size; i++ {
		window[i]=1/float64(window_size)
	}
	return num.Convolve(internal, window, num.Same)
}
func Vandermonde(a []float64, degree int) *mat64.Dense {
	x := mat64.NewDense(len(a), degree+1, nil)
	for i := range a {
		for j, p := 0, 1.; j <= degree; j, p = j+1, p*a[i] {
			x.Set(i, j, p)
		}
	}
	return x
}