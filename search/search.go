package search

import (
	"errors"
	"regexp"
	"strings"

	"labdb.org/labdb/models"
)

func Matches(re *regexp.Regexp, target string, caseInsensitive bool) bool {
	normTarget := target
	if caseInsensitive {
		normTarget = strings.ToLower(normTarget)
	}
	return re.MatchString(normTarget)
}

func Search(term string, includeSequence bool, person string, types []string) ([]models.Entity, error) {
	normTerm := term
	caseInsensitive := false
	if normTerm[0] == '/' {
		if strings.HasSuffix(normTerm, "/i") {
			caseInsensitive = true
			normTerm = strings.Replace(normTerm, "/i", "", -1)
			normTerm = strings.Replace(normTerm, "/", "", -1)
		} else if strings.HasSuffix(normTerm, "/") {
			normTerm = strings.Replace(normTerm, "/", "", -1)
		} else {
			return []models.Entity{}, errors.New("malformed regular expression query")
		}
	} else {
		normTerm = "^" + strings.Replace(normTerm, "*", ".*", -1) + "$"
	}
	regex, err := regexp.Compile(normTerm)
	if err != nil {
		return []models.Entity{}, err
	}

	results := []models.Entity{}
	for _, t := range types {
		query := models.Db()
		normType := strings.ToLower(t)
		obj := models.Empty(normType)
		if person != "" {
			query = query.Where(obj.OwnerFieldName()+" = ?", person)
		}
		query.Order("id desc")
		queryResultsIter := models.RunQueryLazy(normType, query)
		var result models.Entity
		for ; queryResultsIter.HasNext(); result = queryResultsIter.Next() {
			if result == nil {
				continue
			}
			if Matches(regex, result.ShortDesc(), caseInsensitive) {
				results = append(results, result)
				continue
			}
			if Matches(regex, result.Desc(), caseInsensitive) {
				results = append(results, result)
				continue
			}
			if includeSequence && Matches(regex, result.GetSequence(), caseInsensitive) {
				results = append(results, result)
				continue
			}
		}
	}
	return results, nil
}
