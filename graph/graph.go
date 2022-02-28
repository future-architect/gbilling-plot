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
package graph

import (
	"bytes"
	"log"
	"strings"

	"github.com/future-architect/gbilling-plot/invoice"
	chart "github.com/wcharczuk/go-chart/v2"
)

func GetChartValues(costList invoice.CostList) ([]chart.Value, error) {
	var values []chart.Value
	total := 0.0
	for _, cost := range costList {
		values = append(values, chart.Value{
			Value: cost.Cost,
			Label: strings.Replace(cost.Project, "monom-", "", -1),
		})
		total += cost.Cost
	}

	values = append(values, chart.Value{
		Value: total,
		Label: "total",
	})

	return values, nil
}

func Draw(costList invoice.CostList) ([]byte, error) {

	values, _ := GetChartValues(costList)

	log.Println(values)
	graph := chart.BarChart{
		Title: "Cumulative cost by project",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 45,
			},
		},
		Width:        1024,
		Height:       512,
		BarWidth:     50,
		UseBaseValue: true,
		BaseValue:    0.0,
		Bars:         values,
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)

	return buffer.Bytes(), err
}
