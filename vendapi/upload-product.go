// Package vendapi handles interactions with the Vend API.
package vendapi

// UploadProduct ...
type UploadProduct struct {
	Handle   *string `json:"handle,omitempty"`
	SKU      *string `json:"sku,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
}

/*
EXAMPLE PAYLOAD:
{
    "handle": "Mulcher-life",
    "sku": "20049",
    "image_url": "https://s3.amazonaws.com/vend-images/product/standard/9/6/96b22c38f6e84a638663fb1f768a86e19b393a16.jpg",
}
*/
