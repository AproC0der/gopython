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
	degree1=2
	degree2=2
	v=1
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
	/*
	step1 去异常值至少2次!
	如果当前数据点数值*1.5<当前数据点前三个点的平均值，则将当前点数据做删除处理
	*/
	vals1 := make([]float64, 0)
	times1 := make([]int, 0)
	for j := 0; j < 2; j++ {
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
	/*
	step2降噪,至少3次!
	*/
	round1 := Moving_average(vals1, window_size)
	round2 := Moving_average(round1, window_size)
	round3 := Moving_average(round2, window_size)
	/*
	step3 多项式拟合
	times2以天为单位[0, 1, 1.2, 2 ,...]
	采用16阶多项式进行拟合
	*/
	times2 := make([]float64, 0)
	for i := 0; i < len(times1); i++ {
		times2 = append(times2, float64(times[i])-float64(times[0])/float64(60)/float64(60)/float64(24))
	}
	a := Vandermonde(round3, degree1)
	b := mat64.NewDense(len(times2), 1,times2)
	c := mat64.NewDense(degree1+1, 1, nil)
	qr := new(mat64.QR)
	qr.Factorize(a)
	err = c.SolveQR(qr, false, b)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("%.3f\n", mat64.Formatted(c))
	}
	//TODO c->[]float64
	poly_y:=Polyval(c,times2)
	for _, v := range poly_y {
		fmt.Println(v)
	}
	/*
	step4 数据向上处理
	如果当前点比前一点下降d,则后续所有的数据点数值加上1.5d
	*/
	for i:=0;i<len(poly_y);i++{
		if d:=poly_y[i]-poly_y[i+1];d>0{
			for j:=i+1;j<len(poly_y);j++{
				poly_y[j]+=d*1.5
			}
		}
	}
	/*
	step5 多项式拟合
	再次采用多项式拟合，之前degree1=16阶，这次degree2=20阶
	*/
	a := Vandermonde(poly_y, degree2)
	b := mat64.NewDense(len(times2), 1,times2)
	p1 := mat64.NewDense(degree2+1, 1, nil)
	qr := new(mat64.QR)
	qr.Factorize(a)
	err = p1.SolveQR(qr, false, b)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("%.3f\n", mat64.Formatted(p1))
	}
	//TODO p1->[]float64
	/*
	求导数
	*/
	pd1:=make([]float64,0)
	for i:=0;i<len(p1)-1;i++{
		pd1=append(pd1,p1[i]*(len(p1)-1))
	}
	pd1_y:=polyval(pd1,times2)
	for _, v := range pd1_y {
		fmt.Println(v)
	}
	/*
	求切线角
	*/
	angles:=make([]float64,0)
	for _, v := range pd1_y {
		a:=math.arctan(v*180/v/3.14)
		angles=append(angles,a)
	}
	for _, v := range angles {
		fmt.Println(v)
	}
	/*
	求加速度
	*/
	pd2:=make([]float64,0)
	for i:=0;i<len(pd1)-1;i++{
		pd2=append(pd2,p1[i]*(len(pd1)-1))
	}
	pd2_y:=polyval(pd2,times2)
	for _, v := range pd2_y {
		fmt.Println(v)
	}
}
/*
Moving_average -卷积
*/
func Moving_average(internal []float64,window_size int) []float64{
	window:=make([]float64,window_size)
	for i := 0; i < window_size; i++ {
		window[i]=1/float64(window_size)
	}
	return num.Convolve(internal, window, num.Same)
}
/*
Vandermonde -
*/
func Vandermonde(a []float64, degree int) *mat64.Dense {
	x := mat64.NewDense(len(a), degree+1, nil)
	for i := range a {
		for j, p := 0, 1.; j <= degree; j, p = j+1, p*a[i] {
			x.Set(i, j, p)
		}
	}
	return x
}
/*
Polyva -计算多项式值
eg. polyval([3,0,1], 5)  # 3 * 5^2 + 0 * 5^1 + 1=76 依此类推
*/
func Polyval(c []float64,times2 []float64) []float64{
	result:=make([]float64,0)
	for i:=0;i<len(times2);i++{
		var r float64
		for j,k:=len(c)-1,0;j>=0;j--,k++{
			r+=c[k]*math.pow(times2[i],j)
		}
		result=append(result,r)
	}
	return result
}