package dto

type GetBucketsRequest struct {
	LinksList []int64 `json:"links_list"`
	LinkList  *int64  `json:"link_list"`
}

func (g *GetBucketsRequest) Validate() bool {
	if len(g.LinksList) > 0 {
		for _, l := range g.LinksList {
			if l < 0 {
				return false
			}
		}
		return true
	}
	if g.LinkList != nil && *g.LinkList >= 0 {
		return true
	}
	return false
}

func (g *GetBucketsRequest) ProcessingDTO() {
	if g.LinkList != nil {
		g.LinksList = append(g.LinksList, *g.LinkList)
		g.LinkList = nil
	}
}
