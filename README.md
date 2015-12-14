## vend-image-upload

### Summary
Intended to make uploading images to Vend easier to do in bulk, by reading a CSV file containing `sku, handle, image_url` for products and then proceeding to grab and upload the given product to Vend.

### Notes/Known issues
-URLs must end with an image extension (.jpg, .png, etc)
-CSV column headers must be in exact order (sku, handle, image_url)
