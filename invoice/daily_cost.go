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

func (list CostList) MaxCost() float64 {
	max := 0.0
	for _, v := range list {
		if max < v.Cost {
			max = v.Cost
		}
	}
	return max
}

func (list CostList) Values() []float64 {
	result := make([]float64, 0, len(list))
	for _, v := range list {
		result = append(result, v.Cost)
	}
	return result
}

// return zero padded list to adopt month length
func (list CostList) Padding() CostList {
	begin := time.Now().AddDate(0, 0, -29)
	end := time.Now()
	for d := begin; d.Before(end); d = d.AddDate(0, 0, 1) {
		isFind := false
		for i := 0; i < len(list); i++ {
			if d.Format("2006-01-02") == list[i].Date {
				isFind = true
				break
			}
		}
		if !isFind {
			list = append(list, Cost{d.Format("2006-01-02"), "", 0})
		}
	}
	return list
}

func (list CostList) Sort() CostList {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Date < list[j].Date
	})
	return list
}
