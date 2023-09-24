package models

type Dept struct {
	DeptCode int `json:"deptCode"`
	DeptName string `json:"deptName"`
	DeptShortName string `json:"deptShortName"`
}