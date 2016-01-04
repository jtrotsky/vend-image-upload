// Package vendapi handles interactions with the Vend API.
package vendapi

// ProductUpload contains the fields needed to post an image to a product in Vend.
type ProductUpload struct {
	ID       string `json:"id,omitempty"`
	Handle   string `json:"handle,omitempty"`
	SKU      string `json:"sku,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

/*
EXAMPLE PAYLOAD:
{
		"id": "e52b2846-e925-11e5-f98b-40a345e4e0bb"
    "handle": "Mulcher-life",
    "sku": "20049",
    "image_url": "https://www.vendhq.com/images/icns/vend-logo.svg",
}
*/
