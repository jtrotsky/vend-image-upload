// Package vendapi handles interactions with the Vend API.
package vendapi

import "time"

// Product is the basic Vend product structure.
// TODO: Any reason to use this and not default govend struct?
type Product struct {
	ID              *string        `json:"id,omitempty"`
	SourceID        *string        `json:"source_id,omitempty"`
	SourceVariantID *string        `json:"source_variant_id,omitempty"`
	VariantParentID *string        `json:"variant_parent_id,omitempty"`
	Name            *string        `json:"name,omitempty"`
	VariantName     *string        `json:"variant_name,omitempty"`
	Handle          *string        `json:"handle,omitempty"`
	SKU             *string        `json:"sku,omitempty"`
	Active          *bool          `json:"active,omitempty"`
	Source          *string        `json:"source,omitempty"`
	DeletedAt       *time.Time     `json:"deleted_at,omitempty"`
	Version         *int64         `json:"version,omitempty"`
	ImageURL        *string        `json:"image_url,omitempty"`
	Images          []ProductImage `json:"images,omitempty"`
}

// ProductImage contains each product image object.
type ProductImage struct {
	ID      *string `json:"id,omitempty"`
	URL     *string `json:"url,omitempty"`
	Version *int64  `json:"version,omitempty"`
}

/*
ENDPOINT:
<domain_prefix>.vendhq.com/api/2.0/products/<product_id>

EXAMPLE PAYLOAD:
{
  "data": {
    "id": "b8ca3a65-0125-11e4-fbb5-6fdf12990b2f",
    "source_id": null,
    "source_variant_id": null,
    "variant_parent_id": null,
    "name": "Mulching – Life Lessons ( Novel)",
    "variant_name": "Mulching â Life Lessons ( Novel)",
    "handle": "Mulcher-life",
    "sku": "20049",
    "active": true,
    "has_inventory": true,
    "is_composite": false,
    "description": "<p>To top it off, you always need a water feature.</p>",
    "image_url": "https://s3.amazonaws.com/vend-images/product/standard/9/6/96b22c38f6e84a638663fb1f768a86e19b393a16.jpg",
    "created_at": "2014-11-19T11:27:40+00:00",
    "updated_at": "2015-05-07T23:39:37+00:00",
    "deleted_at": null,
    "source": "USER",
		"supplier_code": null,
    "version": 99614596,
    "type": {
      "id": "b8ca3a65-0125-11e4-f728-ea1391a470a2",
      "name": "Mulch",
      "deleted_at": null,
      "version": 1097560
    },
    "supplier": {
      "id": "b8ca3a65-0125-11e4-f728-ea1391a185e6",
      "name": "Jim Jeff's Plumbing",
      "deleted_at": null,
      "version": 1707214
    },
    "brand": {
      "description": null,
      "source": "USER",
      "id": "b8ca3a65-0125-11e4-fbb5-598464cf6a01",
      "name": "Kooringal",
      "deleted_at": null,
      "version": 1295225
    },
    "variant_options": [
		 {
			"id": "a0369f1f-9025-11e4-f68e-e4d622274ed8",
			"name": "Consistency",
			"value": "Moist"
			}
		],
    "categories": [
      {
        "id": "a0369f1f-9025-11e4-f68e-e4d622274ed8",
        "name": "Large",
        "deleted_at": null,
        "version": 1757614
      },
      {
        "id": "a0369f1f-9025-11e4-f68e-e4d623e15741",
        "name": "Moist",
        "deleted_at": null,
        "version": 1757616
      },
      {
        "id": "a0369f1f-9025-11e4-f68e-e4d62994da4b",
        "name": "Good for the plants",
        "deleted_at": null,
        "version": 1757617
      }
    ],
    "images": [
      {
        "id": "a0369f1f-9025-11e4-f68e-d37417b49088",
        "url": "https://s3.amazonaws.com/vend-images/product/original/9/6/96b22c38f6e84a638663fb1f768a86e19b393a16.jpg",
        "version": 1190254
      }
    ],
    "image_thumbnail_url": "https://s3.amazonaws.com/vend-images/product/thumb/9/6/96b22c38f6e84a638663fb1f768a86e19b393a16.jpg"
  }
}
*/
