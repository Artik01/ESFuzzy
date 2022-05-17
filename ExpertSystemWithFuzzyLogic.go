package main

import (
	"fmt"
	"math"
)

type (
	input struct {
		OBrightness float64
		LBrightness float64
		ITemperature float64
		WHTemperature float64
		Humidity float64
	}
	Knowledge struct {
		IF func(input) bool
		Formula func(input) float64
		THEN FuzzySet
		DecisionVar string
	}
	Variable map[string]FuzzySet
	FuzzySet []FuzzyValue
	FuzzyValue struct {
		x float64
		MDegree float64
	}
)

var (
	outsideBrightness = Variable{"dark":OBdark,"dim":OBdim,"avarage":OBavg,"bright":OBbright}
	lampsBrightness = Variable{"turned off":LBoff,"dark":LBdark,"dim":LBdim,"avarage":LBavg,"bright":LBbright}
	insideTemperature = Variable{"very cold":ITvcold,"cold":ITcold,"comfortable":ITcomf,"warm":ITwarm,"hot":IThot,"very hot":ITvhot}
	whTemperature = Variable{"very cold":WHTvcold,"cold":WHTcold,"comfortable":WHTcomf,"hot":WHThot,"very hot":WHTvhot}
	humidity = Variable{"dry":Hdry,"comfortable":Hcomf,"wet":Hwet}
)

//outsideBrightness
var (
	OBdark = FuzzySet{FV(0,1),FV(12,0)}
	OBdim = FuzzySet{FV(10,0),FV(30,1),FV(50,0)}
	OBavg = FuzzySet{FV(40,0),FV(60,1),FV(80,0)}
	OBbright = FuzzySet{FV(70,0),FV(100,1)}
)

//lampsBrightness
var (
	LBoff = FuzzySet{FV(0,1),FV(1,0)}
	LBdark = FuzzySet{FV(0,0),FV(3,1),FV(9,1),FV(12,0)}
	LBdim = FuzzySet{FV(10,0),FV(30,1),FV(50,0)}
	LBavg = FuzzySet{FV(40,0),FV(60,1),FV(80,0)}
	LBbright = FuzzySet{FV(70,0),FV(100,1)}
)

//insideTemperature
var (
	ITvcold = FuzzySet{FV(0,1),FV(10,1),FV(12,0)}
	ITcold = FuzzySet{FV(10,0),FV(12,1),FV(15,1),FV(17,0)}
	ITcomf = FuzzySet{FV(15,0),FV(17,1),FV(23,1),FV(25,0)}
	ITwarm = FuzzySet{FV(22,0),FV(25,1),FV(28,0)}
	IThot = FuzzySet{FV(26,0),FV(29,1),FV(32,0)}
	ITvhot = FuzzySet{FV(30,0),FV(35,1),FV(40,1)}
)

//whTemperature
var (
	WHTvcold = FuzzySet{FV(20,1),FV(35,0)}
	WHTcold = FuzzySet{FV(30,0),FV(40,1),FV(50,0)}
	WHTcomf = FuzzySet{FV(45,0),FV(50,1),FV(60,1),FV(65,0)}
	WHThot = FuzzySet{FV(60,0),FV(65,1),FV(70,1),FV(75,0)}
	WHTvhot = FuzzySet{FV(70,0),FV(80,1)}
)

//humidity
var (
	Hdry = FuzzySet{FV(0,1),FV(20,1),FV(30,0)}
	Hcomf = FuzzySet{FV(20,0),FV(30,1),FV(60,1),FV(70,0)}
	Hwet = FuzzySet{FV(60,0),FV(70,1),FV(100,1)}
)

var KB = []Knowledge{
	/*
	 * 01) IF brightness outside is dark AND (lamp brightness is turned off OR dark OR dim) THEN set lamp brightness to bright
	 * 02) IF brightness outside is dim AND (lamp brightness is turned off OR dark OR dim) THEN set lamp brightness to avg brightness
	 * 03) IF brightness outside is bright AND lamp brightness is NOT turned off THEN set lamp brightness to turned off\
	 * 04) IF brightness outside is avg brightness AND lamp brightness is NOT dim THEN set lamp brightness to dim
	 * */
	Knowledge{
		func(i input) bool {
			return outsideBrightness.Is(i.OBrightness, "dark") && (lampsBrightness.Is(i.LBrightness, "turned off") || lampsBrightness.Is(i.LBrightness, "dark") || lampsBrightness.Is(i.LBrightness, "dim"))
		},
		func(i input) float64 {
			return Min(outsideBrightness.Get(i.OBrightness, "dark"), Max(lampsBrightness.Get(i.LBrightness, "turned off"), lampsBrightness.Get(i.LBrightness, "dark"), lampsBrightness.Get(i.LBrightness, "dim")))
		},
		lampsBrightness["bright"],"Lamp brightness"},
	Knowledge{
		func(i input) bool {
			return outsideBrightness.Is(i.OBrightness, "dim") && (lampsBrightness.Is(i.LBrightness, "turned off") || lampsBrightness.Is(i.LBrightness, "dark") || lampsBrightness.Is(i.LBrightness, "dim"))
		},
		func(i input) float64 {
			return Min(outsideBrightness.Get(i.OBrightness, "dim"), Max(lampsBrightness.Get(i.LBrightness, "turned off"), lampsBrightness.Get(i.LBrightness, "dark"), lampsBrightness.Get(i.LBrightness, "dim")))
		},
		lampsBrightness["avarage"],"Lamp brightness"},
	Knowledge{
		func(i input) bool {
			return outsideBrightness.Is(i.OBrightness, "bright") && !lampsBrightness.Is(i.LBrightness, "turned off")
		},
		func(i input) float64 {
			return Min(outsideBrightness.Get(i.OBrightness, "bright"), 1-lampsBrightness.Get(i.LBrightness, "turned off"))
		},
		lampsBrightness["turned off"],"Lamp brightness"},
	Knowledge{
		func(i input) bool {
			return outsideBrightness.Is(i.OBrightness, "avarage") && !lampsBrightness.Is(i.LBrightness, "dim")
		},
		func(i input) float64 {
			return Min(outsideBrightness.Get(i.OBrightness, "avarage"), 1-lampsBrightness.Get(i.LBrightness, "dim"))
		},
		lampsBrightness["dim"],"Lamp brightness"},
	
	/*
	 * 05) IF (inside temperature is very cold OR cold) AND humidity is NOT wet THEN use air heater to heat up inside temperature to comfortable
	 * 06) IF (inside temperature is very cold OR cold) AND humidity is wet THEN use air heater to heat up inside temperature to warm
	 * 07) IF (inside temperature is hot OR very hot) AND humidity is dry THEN use AC to cool down inside temperature to comfortable
	 * 08) IF (inside temperature is hot OR very hot) AND humidity is NOT dry THEN use AC to cool down inside temperature to warm
	 * */
	Knowledge{
		func(i input) bool {
			return (insideTemperature.Is(i.ITemperature, "very cold") || insideTemperature.Is(i.ITemperature, "cold")) && !humidity.Is(i.Humidity, "wet")
		},
		func(i input) float64 {
			return Min(Max(insideTemperature.Get(i.ITemperature, "very cold"), insideTemperature.Get(i.ITemperature, "cold")), 1-humidity.Get(i.Humidity, "wet"))
		},
		insideTemperature["comfortable"],"Inside temperature"},
	Knowledge{
		func(i input) bool {
			return (insideTemperature.Is(i.ITemperature, "very cold") || insideTemperature.Is(i.ITemperature, "cold")) && humidity.Is(i.Humidity, "wet")
		},
		func(i input) float64 {
			return Min(Max(insideTemperature.Get(i.ITemperature, "very cold"), insideTemperature.Get(i.ITemperature, "cold")), humidity.Get(i.Humidity, "wet"))
		},
		insideTemperature["warm"],"Inside temperature"},
	Knowledge{
		func(i input) bool {
			return (insideTemperature.Is(i.ITemperature, "very hot") || insideTemperature.Is(i.ITemperature, "hot")) && humidity.Is(i.Humidity, "dry")
		},
		func(i input) float64 {
			return Min(Max(insideTemperature.Get(i.ITemperature, "very hot"), insideTemperature.Get(i.ITemperature, "hot")), humidity.Get(i.Humidity, "dry"))
		},
		insideTemperature["comfortable"],"Inside temperature"},
	Knowledge{
		func(i input) bool {
			return (insideTemperature.Is(i.ITemperature, "very hot") || insideTemperature.Is(i.ITemperature, "hot")) && !humidity.Is(i.Humidity, "dry")
		},
		func(i input) float64 {
			return Min(Max(insideTemperature.Get(i.ITemperature, "very hot"), insideTemperature.Get(i.ITemperature, "hot")), 1-humidity.Get(i.Humidity, "dry"))
		},
		insideTemperature["warm"],"Inside temperature"},
	
	/*
	 * 09) IF humidity is dry THEN use humidifier to set humidity to comfortable
	 * 10) IF humidity is wet THEN use dehumidifier to set humidity to comfortable
	 * */
	Knowledge{
		func(i input) bool {
			return humidity.Is(i.Humidity, "dry")
		},
		func(i input) float64 {
			return humidity.Get(i.Humidity, "dry")
		},
		humidity["comfortable"],"Humidity"},
	Knowledge{
		func(i input) bool {
			return humidity.Is(i.Humidity, "wet")
		},
		func(i input) float64 {
			return humidity.Get(i.Humidity, "wet")
		},
		humidity["comfortable"],"Humidity"},
	
	/*
	 * 11) IF WH temperature is NOT hot AND (inside temperature is cold OR very cold) AND humidity is NOT wet THEN set WH temperature to hot
	 * 12) IF WH temperature is NOT hot AND (inside temperature is cold OR very cold) AND humidity is wet THEN set WH temperature to comfortable
	 * 13) IF WH temperature is NOT comfortable AND (inside temperature is comfortable OR warm) THEN set WH temperature to comfortable
	 * 14) IF WH temperature is NOT cold AND (inside temperature is hot OR very hot) AND humidity is NOT dry THEN set WH temperature to cold
	 * 15) IF WH temperature is NOT cold AND (inside temperature is hot OR very hot) AND humidity is dry THEN set WH temperature to comfortable
	 * */
	Knowledge{
		func(i input) bool {
			return !whTemperature.Is(i.WHTemperature, "hot") && (insideTemperature.Is(i.ITemperature, "very cold") || insideTemperature.Is(i.ITemperature, "cold")) && !humidity.Is(i.Humidity, "wet")
		},
		func(i input) float64 {
			return Min(1-whTemperature.Get(i.WHTemperature, "hot"), Max(insideTemperature.Get(i.ITemperature, "very cold"), insideTemperature.Get(i.ITemperature, "cold")), 1-humidity.Get(i.Humidity, "wet"))
		},
		whTemperature["hot"],"Water heater temperature"},
	Knowledge{
		func(i input) bool {
			return !whTemperature.Is(i.WHTemperature, "hot") && (insideTemperature.Is(i.ITemperature, "very cold") || insideTemperature.Is(i.ITemperature, "cold")) && humidity.Is(i.Humidity, "wet")
		},
		func(i input) float64 {
			return Min(1-whTemperature.Get(i.WHTemperature, "hot"), Max(insideTemperature.Get(i.ITemperature, "very cold"), insideTemperature.Get(i.ITemperature, "cold")), humidity.Get(i.Humidity, "wet"))
		},
		whTemperature["comfortable"],"Water heater temperature"},
	Knowledge{
		func(i input) bool {
			return !whTemperature.Is(i.WHTemperature, "comfortable") && (insideTemperature.Is(i.ITemperature, "comfortable") || insideTemperature.Is(i.ITemperature, "warm"))
		},
		func(i input) float64 {
			return Min(1-whTemperature.Get(i.WHTemperature, "comfortable"), Max(insideTemperature.Get(i.ITemperature, "comfortable"), insideTemperature.Get(i.ITemperature, "warm")))
		},
		whTemperature["comfortable"],"Water heater temperature"},
	Knowledge{
		func(i input) bool {
			return !whTemperature.Is(i.WHTemperature, "cold") && (insideTemperature.Is(i.ITemperature, "very hot") || insideTemperature.Is(i.ITemperature, "hot")) && !humidity.Is(i.Humidity, "dry")
		},
		func(i input) float64 {
			return Min(1-whTemperature.Get(i.WHTemperature, "cold"), Max(insideTemperature.Get(i.ITemperature, "very hot"), insideTemperature.Get(i.ITemperature, "hot")), 1-humidity.Get(i.Humidity, "dry"))
		},
		whTemperature["cold"],"Water heater temperature"},
	Knowledge{
		func(i input) bool {
			return !whTemperature.Is(i.WHTemperature, "cold") && (insideTemperature.Is(i.ITemperature, "very hot") || insideTemperature.Is(i.ITemperature, "hot")) && humidity.Is(i.Humidity, "dry")
		},
		func(i input) float64 {
			return Min(1-whTemperature.Get(i.WHTemperature, "cold"), Max(insideTemperature.Get(i.ITemperature, "very hot"), insideTemperature.Get(i.ITemperature, "hot")), humidity.Get(i.Humidity, "dry"))
		},
		whTemperature["comfortable"],"Water heater temperature"}}

func FV(x float64, m float64) FuzzyValue {
	return FuzzyValue{x, m}
}

func main() {
	var in input = input{-1,-1,-1,-1,-1}
	for in.OBrightness+0.0001<0 || in.OBrightness-0.0001>100 {
		fmt.Print("Outside brightness [0-100%] = ")
		fmt.Scan(&(in.OBrightness))
		if in.OBrightness+0.0001<0 || in.OBrightness-0.0001>100 {
			fmt.Println("OUT OF RANGE!")
		}
	}
	for in.LBrightness+0.0001<0 || in.LBrightness-0.0001>100 {
		fmt.Print("Lamps brightness [0-100%] = ")
		fmt.Scan(&(in.LBrightness))
		if in.LBrightness+0.0001<0 || in.LBrightness-0.0001>100 {
			fmt.Println("OUT OF RANGE!")
		}
	}
	for in.ITemperature+0.0001<0 || in.ITemperature-0.0001>40 {
		fmt.Print("Inside temperature [0-40°C] = ")
		fmt.Scan(&(in.ITemperature))
		if in.ITemperature+0.0001<0 || in.ITemperature-0.0001>40 {
			fmt.Println("OUT OF RANGE!")
		}
	}
	for in.WHTemperature+0.0001<20 || in.WHTemperature-0.0001>80 {
		fmt.Print("Water heater temperature [20-80°C] = ")
		fmt.Scan(&in.WHTemperature)
		if in.WHTemperature+0.0001<20 || in.WHTemperature-0.0001>80 {
			fmt.Println("OUT OF RANGE!")
		}
	}
	for in.Humidity+0.0001<0 || in.Humidity-0.0001>100 {
		fmt.Print("Humidity [0-100%] = ")
		fmt.Scan(&in.Humidity)
		if in.Humidity+0.0001<0 || in.Humidity-0.0001>100 {
			fmt.Println("OUT OF RANGE!")
		}
	}
	
	var PreResults map[string][]FuzzySet = make(map[string][]FuzzySet, 5)
	for _, k := range KB {
		if k.IF(in) {
			res:=k.THEN.Multiply(k.Formula(in))
			if _, ok := PreResults[k.DecisionVar]; !ok {
				PreResults[k.DecisionVar]=[]FuzzySet{res}
			} else {
				PreResults[k.DecisionVar]=append(PreResults[k.DecisionVar], res)
			}
		}
	}
	
	var Results map[string]float64 = make(map[string]float64, 5)
	
	for k, sets := range PreResults {
		Results[k]=Union(sets).Centroid()
	}
	
	
	if r, ok := Results["Lamp brightness"]; ok {
		if r<1 {
			fmt.Printf("Lamps should be turned off.\n")
		} else if r>in.LBrightness {
			fmt.Printf("Lamps brightness should be increased up to %.0f%%.\n", r)
		} else if r<in.LBrightness {
			fmt.Printf("Lamps brightness should be decreased down to %.0f%%.\n", r)
		}
	}
	if r, ok := Results["Inside temperature"]; ok {
		if r>in.ITemperature {
			fmt.Printf("Inside temperature should be heated up til %.1f°C, using air heater.\n", r)
		} else if r<in.ITemperature {
			fmt.Printf("Inside temperature should be cooled down til %.1f°C, using AC.\n", r)
		}
	}
	if r, ok := Results["Humidity"]; ok {
		if r>in.Humidity {
			fmt.Printf("Humidity should be increased up to %.0f%%, using humidifier.\n", r)
		} else if r<in.Humidity {
			fmt.Printf("Humidity should be decreased down to %.0f%%, using dehumidifier.\n", r)
		}
	}
	if r, ok := Results["Water heater temperature"]; ok {
		if r>in.WHTemperature {
			fmt.Printf("Water heater temperature should be heated up til %.1f°C.\n", r)
		} else if r<in.WHTemperature {
			fmt.Printf("Water heater temperature should be cooled down til %.1f°C.\n", r)
		}
	}
	
	fmt.Println("\n\nTo exit this programm close the window or press Ctrl+C.")
	for {}
}

func (fz FuzzySet) Multiply(b float64) FuzzySet {
	res := make(FuzzySet, len(fz))
	copy(res, fz)
	for i:=0; i < len(res); i++ {
		res[i].MDegree*=b
	}
	return res
}

func (fz FuzzySet) Get(x float64) float64 {
	for _,y := range fz {
		if y.x-0.001<=x && y.x+0.001>=x {
			return y.MDegree
		}
	}
	if x<fz[0].x {
		return fz[0].MDegree
	}
	if x>fz[len(fz)-1].x {
		return fz[len(fz)-1].MDegree
	}
	var a int
	for i, y := range fz {
		if y.x<x {
			a=i
		} else {
			break
		}
	}
	dx:=fz[a+1].x-fz[a].x
	dm:=fz[a+1].MDegree-fz[a].MDegree
	dm*=((x-fz[a].x)/dx)
	return fz[a].MDegree+dm
}

func Union(sets []FuzzySet) FuzzySet {
	if len(sets) == 0 {
		return nil
	}
	var min, max float64
	
	min, max = sets[0][0].x, sets[0][0].x
	for _, fz := range sets {
		for _, y := range fz {
			if y.x<=min {
				min = y.x
			}
			if y.x>=max {
				max=y.x
			}
		}
	}
	
	var new FuzzySet
	for x:=min; x<=max+0.001; x+=0.1 {
		var m float64 = -1
		for _, fz := range sets {
			m=math.Max(m, fz.Get(x))
		}
		new = append(new, FV(x, m))
	}
	return new
}

func (fz FuzzySet) Centroid() float64 {
	sum1, sum2:=0.0,0.0
	for _, x := range fz {
		sum1+=x.x * x.MDegree
		sum2+=x.MDegree
	}
	return sum1/sum2
}

func (v Variable) Is(x float64, class string) bool {
	if fz, ok := v[class]; ok {
		return fz.Get(x)>0.00001
	} else {
		return false
	}
}

func (v Variable) Get(x float64, class string) float64 {
	return v[class].Get(x)
}

func Max(a float64, args ...float64) float64 {
	for _, b := range args {
		a=math.Max(a, b)
	}
	return a
}

func Min(a float64, args ...float64) float64 {
	for _, b := range args {
		a=math.Min(a, b)
	}
	return a
}
