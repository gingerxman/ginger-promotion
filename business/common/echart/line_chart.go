package common

import (
	"encoding/json"
	"github.com/gingerxman/eel"
)

const lineChartTmpl = `
{
	"title" : {
		"text": ""
	},
	"tooltip" : {
		"trigger": "axis"
	},
	"toolbox": {
		"show": true,
		"feature": {
			"restore" : {"show": true},
			"saveAsImage" : {"show": true}
		}
	},
	"calculable" : true,
	"xAxis" : [
		{
			"type" : "category",
			"boundaryGap" : false,
			"data" : []
		}
	],
	"yAxis" : [
		{
			"type" : "value"
		}
	],
	"series" : [
		{
			"name":"",
			"type":"line",
			"smooth":true,
			"itemStyle": {"normal": {"areaStyle": {"type": "default"}}},
			"data":[]
		}
	]
}
`

type ChartPoint struct {
	X string
	Y float64
}

type LineChartInfo struct {
	Title string
	DataName string
	Points []*ChartPoint
}

func CreateLineChart(info *LineChartInfo) eel.Map {
	chartData := eel.Map{}
	err := json.Unmarshal([]byte(lineChartTmpl), &chartData)
	if err != nil {
		eel.Logger.Error(err)
		return chartData
	}
	
	xs := make([]string, 0)
	ys := make([]float64, 0)
	for _, point := range info.Points {
		xs = append(xs, point.X)
		ys = append(ys, point.Y)
	}
	chartData["title"].(eel.Map)["text"] = info.Title
	
	firstSeries := chartData["series"].([]interface{})[0].(map[string]interface{})
	firstSeries["name"] = info.DataName
	firstSeries["data"]= ys
	
	firstXAxis := chartData["xAxis"].([]interface{})[0].(map[string]interface{})
	firstXAxis["data"] = xs
	
	return chartData
}