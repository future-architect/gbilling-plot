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
