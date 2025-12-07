package dto

type GetBucketInfoResponse struct {
	Links     map[string]string `json:"links"`
	NumBucket int64             `json:"num_bucket"`
}

func CreateGetInfoBucketDTO(links map[string]string, numBucket int64) *GetBucketInfoResponse {
	if links == nil {
		return nil
	}
	res := &GetBucketInfoResponse{}
	res.Links = links
	res.NumBucket = numBucket
	return res
}
