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
	"image/color"
	"math"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/future-architect/gbilling-plot/invoice"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func Draw(costList invoice.CostList) ([]byte, error) {

	p, w, err := newPlot()
	if err != nil {
		return nil, err
	}

	if err := addBarChart(p, w, costList); err != nil {
		return nil, err
	}

	// max y axis
	p.Y.Max = calcYAxisTop(costList, p.Y.Max)

	to, err := p.WriterTo(10*vg.Inch, 5*vg.Inch, "png")
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if _, err := to.WriteTo(&buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func newPlot() (*plot.Plot, vg.Length, error) {

	p, err := plot.New()
	if err != nil {
		return nil, vg.Length(0), err
	}

	p.Title.Text = "Stacked Bar Chart on Projects"
	p.Y.Label.Text = "yen (jp)"
	p.Legend.Top = true
	p.Legend.Left = true

	grid := plotter.NewGrid()
	grid.Horizontal.Color = color.Black
	grid.Horizontal.Dashes = []vg.Length{vg.Length(5)}
	p.Add(grid)

	w := vg.Points(15)

	today := time.Now()
	p.NominalX(today.AddDate(0, 0, -29).Format("01/02"), "", "", "", "", today.AddDate(0, 0, -24).Format("01/02"), "", "", "", "", today.AddDate(0, 0, -19).Format("01/02"), "", "", "", "", today.AddDate(0, 0, -14).Format("01/02"), "", "", "", "", today.AddDate(0, 0, -9).Format("01/02"), "", "", "", "", today.AddDate(0, 0, -4).Format("01/02"), "", "", "", today.Format("01/02"))

	return p, w, nil
}

func addBarChart(p *plot.Plot, w vg.Length, costList invoice.CostList) error {

	var beforeBarChar *plotter.BarChart
	currentProject := ""
	colorCount := 0

	plotCostList := make(invoice.CostList, 0, len(costList))

	for _, c := range costList {
		if currentProject == "" {
			currentProject = c.Project
		}

		if currentProject != c.Project {
			plotCostList = plotCostList.Padding().SortByDate()
			barChart, err := newBarChart(colorCount, plotCostList.Values(), w)
			if err != nil {
				return err
			}
			if beforeBarChar != nil {
				barChart.StackOn(beforeBarChar)
			}
			p.Add(barChart)
			p.Legend.Add(currentProject, barChart)

			beforeBarChar = barChart
			currentProject = c.Project

			// empty slice
			plotCostList = plotCostList[:0]
			colorCount++
		}
		plotCostList = append(plotCostList, c)
	}

	plotCostList = plotCostList.Padding().SortByDate()
	barChart, err := newBarChart(colorCount, plotCostList.Values(), w)
	if err != nil {
		return err
	}
	barChart.StackOn(beforeBarChar)
	p.Add(barChart)
	p.Legend.Add(currentProject, barChart)

	return nil
}

func newBarChart(counter int, plotValues plotter.Values, w vg.Length) (*plotter.BarChart, error) {
	barChart, err := plotter.NewBarChart(plotValues, w)
	if err != nil {
		return nil, err
	}

	barChart.Color = plotutil.Color(counter)
	barChart.LineStyle.Width = vg.Length(0)

	return barChart, nil
}

// Calculate Y Axis max length
func calcYAxisTop(costList invoice.CostList, yAxisTop float64) float64 {

	maxCost := costList.MaxCost()
	length := utf8.RuneCountInString(strconv.Itoa(int(maxCost)))

	var yAxisCriteria float64
	switch length {
	case 1, 2:
		yAxisCriteria = 100
	default:
		yAxisCriteria = math.Pow10(length - 1)
	}

	return float64(int(yAxisTop)/int(yAxisCriteria)+1) * yAxisCriteria
}
