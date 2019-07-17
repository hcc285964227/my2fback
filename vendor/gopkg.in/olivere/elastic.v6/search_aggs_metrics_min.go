// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// MinAggregation is a single-value metrics aggregation that keeps track and
// returns the minimum value among numeric values extracted from the
// aggregated documents. These values can be extracted either from
// specific numeric fields in the documents, or be generated by a
// provided script.
// See: https://www.elastic.co/guide/en/elasticsearch/reference/6.8/search-aggregations-metrics-min-aggregation.html
type MinAggregation struct {
	field           string
	script          *Script
	format          string
	missing         interface{}
	subAggregations map[string]Aggregation
	meta            map[string]interface{}
}

func NewMinAggregation() *MinAggregation {
	return &MinAggregation{
		subAggregations: make(map[string]Aggregation),
	}
}

func (a *MinAggregation) Field(field string) *MinAggregation {
	a.field = field
	return a
}

func (a *MinAggregation) Script(script *Script) *MinAggregation {
	a.script = script
	return a
}

func (a *MinAggregation) Format(format string) *MinAggregation {
	a.format = format
	return a
}

func (a *MinAggregation) Missing(missing interface{}) *MinAggregation {
	a.missing = missing
	return a
}

func (a *MinAggregation) SubAggregation(name string, subAggregation Aggregation) *MinAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *MinAggregation) Meta(metaData map[string]interface{}) *MinAggregation {
	a.meta = metaData
	return a
}

func (a *MinAggregation) Source() (interface{}, error) {
	// Example:
	//	{
	//    "aggs" : {
	//      "min_price" : { "min" : { "field" : "price" } }
	//    }
	//	}
	// This method returns only the { "min" : { "field" : "price" } } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["min"] = opts

	// ValuesSourceAggregationBuilder
	if a.field != "" {
		opts["field"] = a.field
	}
	if a.script != nil {
		src, err := a.script.Source()
		if err != nil {
			return nil, err
		}
		opts["script"] = src
	}
	if a.format != "" {
		opts["format"] = a.format
	}
	if a.missing != nil {
		opts["missing"] = a.missing
	}

	// AggregationBuilder (SubAggregations)
	if len(a.subAggregations) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.subAggregations {
			src, err := aggregate.Source()
			if err != nil {
				return nil, err
			}
			aggsMap[name] = src
		}
	}

	// Add Meta data if available
	if len(a.meta) > 0 {
		source["meta"] = a.meta
	}

	return source, nil
}
