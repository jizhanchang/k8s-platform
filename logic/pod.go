package logic

import (
	"context"
	"go.uber.org/zap"
	"k8s-platfrom/lib"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// new pod cell cover to  pod => dataCell  dataCell => pod

type podCell corev1.Pod

func (p podCell) GetName() string {

	return p.Name
}

func (p podCell) GetCreation() time.Time {

	return p.CreationTimestamp.Time
}

var Pod pod

type pod struct {
}

type getPodsResp struct {
	Data  []corev1.Pod
	Total int
}

func (p *pod) toCell(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

func (p *pod) fromCells(std []DataCell) []corev1.Pod {

	pods := make([]corev1.Pod, len(std))

	for i := range std {
		pods[i] = corev1.Pod(std[i].(podCell))
	}
	return pods
}

func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *getPodsResp, err error) {

	pods, err := lib.K8s.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error("fetch pods failed", zap.String("err", err.Error()))
		return nil, err
	}
	selectorData := &dataSelector{
		GenericDataList: p.toCell(pods.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: filterName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	// 过滤
	selectorData = selectorData.Filter()
	total := len(selectorData.GenericDataList)
	//排序分页
	selectorData = selectorData.Sort().Paginate()
	podsResp = &getPodsResp{
		Data:  p.fromCells(selectorData.GenericDataList),
		Total: total,
	}
	return podsResp, nil
}
