package dto

import (
	"slices"
)

type RegistrationLinksRequest struct {
	Links []string `json:"links"`
	Link  *string  `json:"link"`
}

func (r *RegistrationLinksRequest) Validate() bool {
	if len(r.Links) > 0 {
		if r.Link != nil && *r.Link != "" {
			return true
		}
		return !slices.Contains(r.Links, "")
	}
	if r.Link != nil && *r.Link != "" {
		return true
	}
	return false
}

func (r *RegistrationLinksRequest) ProcessingDTO() {
	if r.Link != nil {
		r.Links = append(r.Links, *r.Link)
		r.Link = nil
	}
}
