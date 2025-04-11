package dto

type ProductDTO struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Quantity int32    `json:"quantity"`
	Price    float32  `json:"price"`
	ImgUrls  []string `json:"img_urls"`
}

type GetProductsResponse struct {
	Products []ProductDTO `json:"products"`
}
