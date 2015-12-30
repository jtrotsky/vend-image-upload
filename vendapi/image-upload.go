// Package vendapi handles interactions with the Vend API.
package vendapi

// ImageUpload is the upper level data structure of the response.
type ImageUpload struct {
	Data Data `json:"data,omitempty"`
}

// Data is the information for each image contained in the response.
type Data struct {
	ID        *string `json:"id,omitempty"`
	ProductID *string `json:"product_id,omitempty"`
	Position  *int64  `json:"position,omitempty"`
	Status    *string `json:"status,omitempty"`
	Version   *int64  `json:"version,omitempty"`
}

/*
ENDPOINT:
<domain_prefix>.vendhq.com/api/2.0/<product_id>/actions/image_upload

RESPONSE PAYLOAD:
"data": {
	"id": "b8ca3a6e-7294-11e4-fd8a-e663ebb91b19",
	"product_id": "bc305bf5-da94-11e4-f3a2-b1ae3a5f4c50",
	"position": 1,
	"status": "processing",
	"version": 1245990
}
*/
