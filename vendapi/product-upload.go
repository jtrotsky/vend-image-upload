// Package vendapi handles interactions with the Vend API.
package vendapi

// UploadProduct contains the fields needed to post an image to a product in Vend.
type UploadProduct struct {
	// TODO: Pointer.
	ID       string  `json:"id,omitempty"`
	Handle   *string `json:"handle,omitempty"`
	SKU      *string `json:"sku,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
}

/*
EXAMPLE PAYLOAD:
{
    "handle": "Mulcher-life",
    "sku": "20049",
    "image_url": "https://www.vendhq.com/images/icns/vend-logo.svg",
}
*/
