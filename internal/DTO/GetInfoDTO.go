package dto

import "slices"

type InfoWithNumbersBucketDTO struct {
	LinksList []int64 `json:"links_list"`
	LinkList  *int64  `json:"link_list"`
}

func (n *InfoWithNumbersBucketDTO) Validate() bool {
	if n.LinkList != nil {
		if slices.Contains(n.LinksList, *n.LinkList) {
			return false
		}
		n.LinksList = append(n.LinksList, *n.LinkList)
		n.LinkList = nil
	}

	for _, v := range n.LinksList {
		if v < 0 {
			return false
		}
	}
	return true
}
