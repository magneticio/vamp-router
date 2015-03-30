package haproxy

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

const (
	FILTERS_CORRECT_JSON = "../test/test_filters_correct.json"
	FILTERS_WRONG_JSON   = "../test/test_filters_wrong.json"
)

func TestFilters_ParseFilter(t *testing.T) {

	j, _ := ioutil.ReadFile(FILTERS_CORRECT_JSON)
	var filtersCorrect, filtersWrong []*Filter
	_ = json.Unmarshal(j, &filtersCorrect)

	i, _ := ioutil.ReadFile(FILTERS_WRONG_JSON)
	_ = json.Unmarshal(i, &filtersWrong)

	for _, filter := range filtersCorrect {
		if _, err := parseFilter(filter); err != nil {
			t.Errorf("Failed to correctly parse a filter %s", err.Error())
		}
	}

	for _, filter := range filtersWrong {
		if _, err := parseFilter(filter); err == nil {
			t.Errorf("Filter parsing should fail with incorrect filters")
		}
	}

}

func TestFilters_ParseFilterCondition(t *testing.T) {

	/*
	  these two notations should be equivalent. The full Haproxy condition
	  should pass through untouched
	*/

	type Condition struct {
		Input    string
		Expected string
	}

	tests := []struct {
		Input    string
		Expected string
	}{
		{"hdr_sub(user-agent) Android", "hdr_sub(user-agent) Android"},
		{"user-agent=Android", "hdr_sub(user-agent) Android"},
		{"User-Agent=Android", "hdr_sub(user-agent) Android"},
		{"user-agent = Android", "hdr_sub(user-agent) Android"},
		{"user-agent  =  Android", "user-agent  =  Android"},
		{"host = www.google.com", "hdr_str(host) www.google.com"},
		{"cookie MYCUSTOMER contains Value=good", "cook_sub(MYCUSTOMER) Value=good"},
		{"has cookie JSESSIONID", "cook(JSESSIONID) -m found"},
		{"misses cookie JSESSIONID", "cook_cnt(JSESSIONID) eq 0"},

		{"has header X-SPECIAL", "hdr_cnt(X-SPECIAL) gt 0"},
		{"misses header X-SPECIAL", "hdr_cnt(X-SPECIAL) eq 0"},
	}

	for i, condition := range tests {
		if result := parseFilterCondition(condition.Input); result != condition.Expected {
			t.Errorf("Failed to correctly parse filter condition %d. Got %s but expected %s", (i + 1), result, condition.Expected)
		}
	}
}
