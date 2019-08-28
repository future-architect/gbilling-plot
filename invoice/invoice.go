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
	"cloud.google.com/go/bigquery"
	"context"
	"google.golang.org/api/iterator"
	"time"
)

type invoice struct {
	client *bigquery.Client
}

func NewInvoice(ctx context.Context, projectID string) (*invoice, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &invoice{client: client}, nil
}

func (i *invoice) FetchBilling(ctx context.Context, tableName string, period int) (CostList, error) {

	endDay := time.Now()
	startDay := endDay.AddDate(0, 0, -(period - 1))
	stmt := `
		SELECT
			CAST(sq.date AS string) AS date,
			sq.project AS project,
			sq.cost AS cost
		FROM (
			SELECT
				DATE(_PARTITIONTIME) AS date,
				project.id AS project,
				IFNULL(SUM(cost), 0) AS cost
			FROM
				` + "`" + tableName + "`" + `
			WHERE
				DATE(_PARTITIONTIME) BETWEEN ` + startDay.Format("'2006-01-02'") + `
				AND ` + endDay.Format("'2006-01-02'") + `
				AND project.id IS NOT NULL
			GROUP BY
				DATE(_PARTITIONTIME), project ) AS sq
		ORDER BY
			sq.project,
			sq.date
	`

	iter, err := i.client.Query(stmt).Read(ctx)
	if err != nil {
		return nil, err
	}

	var costList []Cost
	for {
		var c Cost
		err := iter.Next(&c)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		costList = append(costList, c)
	}
	return costList, nil
}
