package dto

type GetInfoBucketDTO struct {
	Links     map[string]string `json:"links"`
	NumBucket int64             `json:"num_bucket"`
}

func CreateGetInfoBucketDTO(links map[string]string, numBucket int64) *GetInfoBucketDTO {
	if links == nil {
		return nil
	}
	res := &GetInfoBucketDTO{}
	res.Links = links
	res.NumBucket = numBucket
	return res
}
