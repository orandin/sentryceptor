package main

import "github.com/orandin/go-ruler"

func filterTags(sf []SentryFilterRules, tags []interface{}) []interface{} {
	tagsFiltered := tags[:0]
	tagMap := make(map[string]interface{}, 1)

	for _, tag := range tags {
		key := tag.([]interface{})[0].(string)
		tagMap[key] = tag.([]interface{})[1]

		if keep(sf, tagMap) {
			tagsFiltered = append(tagsFiltered, tag)
		}

		delete(tagMap, key)
	}
	return tagsFiltered
}

func filterBreadcrumbs(sf []SentryFilterRules, breadcrumbs map[string]interface{}) []interface{} {
	breadcrumbsValues := breadcrumbs["values"].([]interface{})
	breadcrumbsFiltered := breadcrumbsValues[:0]

	for _, breadcrumb := range breadcrumbsValues {
		if keep(sf, breadcrumb.(map[string]interface{})) {
			breadcrumbsFiltered = append(breadcrumbsFiltered, breadcrumb)
		}
	}
	return breadcrumbsFiltered
}

func filterMap(sf []SentryFilterRules, list map[string]interface{}) map[string]interface{} {
	tmpMap := make(map[string]interface{}, 1)

	for key, value := range list {
		tmpMap[key] = value

		if !keep(sf, tmpMap) {
			delete(list, key)
		}
	}
	return list
}

func keep(sf []SentryFilterRules, v map[string]interface{}) bool {
	for _, filterRules := range sf {
		engine := ruler.NewRuler(filterRules.Conditions)

		if ok := engine.Test(v); ok {
			return false
		}
	}
	return true
}
