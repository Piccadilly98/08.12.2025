package dto

import (
	"slices"
)

type RegistrationLinksMap struct {
	Links     map[string]string `json:"links"`
	Link      *string           `json:"link"`
	NumBucket int64
}

func (r *RegistrationLinksMap) Validate() bool {
	if r.Link != nil {
		r.Links[*r.Link] = ""
	}

	for k := range r.Links {
		if k == "" {
			return false
		}
	}
	return true
}

type RegistrationLinks struct {
	Links []string `json:"links"`
	Link  *string  `json:"link"`
}

func (r *RegistrationLinks) Validate() bool {
	if r.Link != nil {
		r.Links = append(r.Links, *r.Link)
	}

	return !slices.Contains(r.Links, "")
}
