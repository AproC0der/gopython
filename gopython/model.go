package main

import "time"

// MonitorPointCrack 裂缝数据
type MonitorPointCrack struct {
	ID               int       `json:"id" xorm:"id pk autoincr"`
	MonitorPointNum  string    `json:"monitorPointNum"`
	MonitorPointName string    `json:"monitorPointName"`
	WTime            time.Time `json:"-"` //时间
	Value            float32   `json:"value"` // 裂缝数据(mm)
	IfType           int       `json:"ifType"`
	DataState        int       `json:"dataState"` // 数据状态 1为正常数据 0为异常数据
	Original         string    `json:"original"`
}
