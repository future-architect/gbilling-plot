/**
 * Copyright (c) 2019-present Future Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package invoice

import (
	"sort"
	"time"
)

type Cost struct {
	Date    string  `bigquery:"date"`
	Project string  `bigquery:"project"`
	Cost    float64 `bigquery:"cost"`
}

type CostList []Cost

func (cl CostList) SummaryLowerProjects(thresholdRank int) CostList {

	sorts := cl.sortByTotalCost()

	var summaryTargets []Cost
	if len(sorts) > thresholdRank {
		summaryTargets = sorts[0 : len(sorts)-thresholdRank]
	}

	var result []Cost
	var summaryCosts = map[string]Cost{}
	for _, v := range cl {
		if CostList(summaryTargets).containsProject(v.Project) {
			cost, ok := summaryCosts[v.Date]
			if ok {
				cost.Cost += v.Cost
				summaryCosts[v.Date] = cost
				continue
			}
			summaryCosts[v.Date] = Cost{
				Date:    v.Date,
				Project: "Others",
				Cost:    v.Cost,
			}
			continue
		}
		result = append(result, v)
	}

	var others []Cost
	for _, v := range summaryCosts {
		others = append(others, v)
	}
	othersCosts := CostList(others).SortByDate()

	// merge
	result = append(result, othersCosts...)

	return result
}

func (cl CostList) MaxCost() float64 {
	max := 0.0
	for _, v := range cl {
		if max < v.Cost {
			max = v.Cost
		}
	}
	return max
}

func (cl CostList) Values() []float64 {
	result := make([]float64, 0, len(cl))
	for _, v := range cl {
		result = append(result, v.Cost)
	}
	return result
}

// return zero padded list to adopt month length
func (cl CostList) Padding() CostList {
	begin := time.Now().AddDate(0, 0, -29)
	end := time.Now()
	result := cl
	for d := begin; d.Before(end); d = d.AddDate(0, 0, 1) {
		isFind := false
		for i := 0; i < len(cl); i++ {
			if d.Format("2006-01-02") == cl[i].Date {
				isFind = true
				break
			}
		}
		if !isFind {
			result = append(result, Cost{d.Format("2006-01-02"), "", 0})
		}
	}
	return result
}

func (cl CostList) SortByDate() CostList {
	sort.Slice(cl, func(i, j int) bool {
		return cl[i].Date < cl[j].Date
	})
	return cl
}

func (cl CostList) containsProject(project string) bool {
	for _, v := range cl {
		if project == v.Project {
			return true
		}
	}
	return false
}

// SortByTotalCost is sum cost per project and sort asc by sum cost
func (cl CostList) sortByTotalCost() CostList {

	sumMap := map[string]float64{}
	for _, v := range cl {
		sumMap[v.Project] += v.Cost
	}

	var sumCosts []Cost
	for k, v := range sumMap {
		sumCosts = append(sumCosts, Cost{
			Project: k,
			Cost:    v,
		})
	}

	return CostList(sumCosts).sortByCost()
}

func (cl CostList) sortByCost() CostList {
	sort.Slice(cl, func(i, j int) bool {
		return cl[i].Cost < cl[j].Cost
	})
	return cl
}
