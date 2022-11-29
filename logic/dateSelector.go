package logic

import (
	"sort"
	"strings"
	"time"
)

/*

定义用于排序\分页以及过滤的数据结构

*/
type dataSelector struct {
	GenericDataList []DataCell
	dataSelectQuery *DataSelectQuery
}

type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

/*
重构 Sort 方法
*/

func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *dataSelector) Less(i, j int) bool {

	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()

	return b.Before(a)
}

func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

/*
   条件过滤方法  只有name过滤其他没写
*/

func (d *dataSelector) Filter() *dataSelector {

	// 如果FilterName 为空 就全部返回
	if d.dataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	var filtered []DataCell
	for _, dat := range d.GenericDataList {
		matcher := true
		if !strings.Contains(dat.GetName(), d.dataSelectQuery.FilterQuery.Name) {
			matcher = false
			continue
		}
		if matcher {
			filtered = append(filtered, dat)
		}

	}
	d.GenericDataList = filtered

	return d
}

/*
  分页  没啥好说的
*/

func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectQuery.PaginateQuery.Limit
	page := d.dataSelectQuery.PaginateQuery.Page
	//如果传过来的参数有问题,就把所有的pod信息都返回,让前端处理
	if limit <= 0 || page <= 0 {
		return d
	}
	startIdx := limit * (page - 1)
	endIdx := limit * page

	if len(d.GenericDataList) < endIdx {
		return d
	}
	d.GenericDataList = d.GenericDataList[startIdx:endIdx]
	return d
}
