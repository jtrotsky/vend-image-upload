## vend-image-upload

### Summary
Intended to make uploading images to Vend easier to do in bulk by taking a CSV file containing `sku, handle, image_url` for products and then proceeding to grab and upload the given product to Vend.

### Howto

Run the tool locally and point it to the CSV file of products/images to upload.

```
$ ./vend-image-upload -d honestmulch -t <access-token> -f ~/Desktop/products_with_images.csv
```

### Limitations/Known issues
* URLs must end with *an* image extension of some sort (`.jpg, .png`, etc)
* CSV column headers must be in exact order (`sku, handle, image_url`)
* Image will be uploaded regardless of "position". If a product in Vend already has an image
in position 1 then the new image will be posted in position 2.
