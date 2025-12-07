package dto

type InfoWithNumbersBucketDTO struct {
	LinksList []int64 `json:"links_list"`
	LinkList  *int64  `json:"link_list"`
}

func (n *InfoWithNumbersBucketDTO) Validate() bool {
	if len(n.LinksList) > 0 {
		for _, l := range n.LinksList {
			if l < 0 {
				return false
			}
		}
		return true
	}
	if n.LinkList != nil && *n.LinkList >= 0 {
		return true
	}
	return false
}

func (n *InfoWithNumbersBucketDTO) ProcessingDTO() {
	if n.LinkList != nil {
		n.LinksList = append(n.LinksList, *n.LinkList)
		n.LinkList = nil
	}
}
