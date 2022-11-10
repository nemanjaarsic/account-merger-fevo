package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"main/models"
	"sort"
)

func errorCheck(err error) bool {
	return err != nil
}

func ImportAccounts(path string) []models.Account {
	content, err := ioutil.ReadFile(path)
	if errorCheck(err) {
		log.Fatal("Error when opening file: ", err)
	}
	var payload []models.Account
	err = json.Unmarshal(content, &payload)
	if errorCheck(err) {
		log.Fatal("Error during json unmarshal: ", err)
	}

	return payload
}

func MergeAccounts(accounts []models.Account) []models.Person {
	owners := map[string]string{}
	apps := map[string][]string{}
	parents := map[string]string{}
	unions := map[string][]string{}

	for _, a := range accounts {
		for _, e := range a.Emails {
			owners[e] = a.Name
			if _, found := apps[e]; !found {
				apps[e] = []string{}
			}
			apps[e] = append(apps[e], a.Application)
			parents[e] = e
		}
	}

	for _, a := range accounts {
		parentEmail := a.Emails[0]
		for i := 1; i < len(a.Emails); i++ {
			parents[a.Emails[i]] = parentEmail
		}
	}

	for _, a := range accounts {
		parentEmail := parents[a.Emails[0]]
		if _, found := unions[parentEmail]; !found {
			unions[parentEmail] = []string{}
		}

		for i := 1; i < len(a.Emails); i++ {
			unions[parentEmail] = append(unions[parentEmail], a.Emails[i])
			apps[parentEmail] = append(apps[parentEmail], apps[a.Emails[i]]...)
			apps[parentEmail] = removeDuplicates[string](apps[parentEmail])
			sort.Strings(apps[parentEmail])
		}
	}

	results := []models.Person{}
	for k, v := range unions {
		emails := v
		emails = append(emails, k)
		sort.Strings(emails)
		results = append(results, models.Person{
			Applications: apps[k],
			Emails:       emails,
			Name:         owners[k],
		})
	}

	return results
}

func removeDuplicates[T comparable](input []T) []T {
	keys := map[T]bool{}
	output := []T{}

	for _, element := range input {
		if _, found := keys[element]; !found {
			keys[element] = true
			output = append(output, element)
		}
	}
	return output
}

func PrintResults(results []models.Person) {
	for _, r := range results {
		log.Print(r)
	}
}
