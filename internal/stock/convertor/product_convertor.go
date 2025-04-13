package convertor

import (
	client "github.com/FacundoChan/gorder-v1/common/client/order"
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"github.com/FacundoChan/gorder-v1/stock/entity"
	// "github.com/FacundoChan/gorder-v1/stock/infrastructure/persistent"
)

type ProductConvertor struct {
}

// type ProductStockModelConvertor struct {
// }
//
func (c *ProductConvertor) EntityToProto(p *entity.Product) *orderpb.Product {
	c.check(p)
	return &orderpb.Product{
		ID:       p.ID,
		Name:     p.Name,
		Quantity: p.Quantity,
		PriceID:  p.PriceID,
		Price:    p.Price,
		ImgUrls:  p.ImgUrls,
	}
}

func (c *ProductConvertor) ProtoToEntity(p *orderpb.Product) *entity.Product {
	c.check(p)
	return &entity.Product{
		ID:       p.ID,
		Name:     p.Name,
		Quantity: p.Quantity,
		PriceID:  p.PriceID,
		Price:    p.Price,
		ImgUrls:  p.ImgUrls,
	}
}

func (c *ProductConvertor) EntityToClient(p *entity.Product) client.Product {
	c.check(p)
	imgUrls := &p.ImgUrls
	return client.Product{
		Id:       p.ID,
		ImgUrls:  imgUrls,
		Name:     p.Name,
		Price:    p.Price,
		Quantity: p.Quantity,
	}
}

func (c *ProductConvertor) ClientToEntity(p client.Product) *entity.Product {
	imgs := *p.ImgUrls
	return &entity.Product{
		ID:       p.Id,
		Name:     p.Name,
		Quantity: p.Quantity,
		Price:    p.Price,
		ImgUrls:  imgs,
	}
}

func (c *ProductConvertor) EntitiesToProtos(products []*entity.Product) (res []*orderpb.Product) {
	for _, p := range products {
		res = append(res, c.EntityToProto(p))
	}
	return
}

func (c *ProductConvertor) ProtosToEntities(products []*orderpb.Product) (res []*entity.Product) {
	for _, p := range products {
		res = append(res, c.ProtoToEntity(p))
	}
	return
}

func (c *ProductConvertor) EntitiesToClients(products []*entity.Product) []client.Product {
	var res []client.Product
	for _, p := range products {
		res = append(res, c.EntityToClient(p))
	}
	return res
}

func (c *ProductConvertor) ClientsToEntities(products []client.Product) []*entity.Product {
	var res []*entity.Product
	for _, p := range products {
		res = append(res, c.ClientToEntity(p))
	}
	return res
}

// func (c *ProductStockModelConvertor) EntitiesToPersistents(productStockModel []entity.StockModel) []persistent.StockModel {
// 	var res []persistent.StockModel
// 	for _, p := range productStockModel {
// 		res = append(res, c.EntityToPersistent(p))
// 	}
// 	return res
// }
//
// func (c *ProductStockModelConvertor) EntityToPersistent(p entity.StockModel) persistent.StockModel {
// 	return persistent.StockModel{
// 		ID:          p.ID,
// 		ProductID:   p.ProductID,
// 		Name:        p.Name,
// 		Quantity:    p.Quantity,
// 		Price:       p.Price,
// 		Description: p.Description,
// 		ImgUrls:     p.ImgUrls,
// 	}
// }
//
// func (c *ProductStockModelConvertor) PersistentsToEntities(productStockModel []persistent.StockModel) []entity.StockModel {
// 	var res []entity.StockModel
// 	for _, p := range productStockModel {
// 		res = append(res, c.PersistentToEntity(p))
// 	}
// 	return res
// }
//
// func (c ProductStockModelConvertor) PersistentToEntity(p persistent.StockModel) entity.StockModel {
// 	return entity.StockModel{
// 		ID:          p.ID,
// 		ProductID:   p.ProductID,
// 		Name:        p.Name,
// 		Quantity:    p.Quantity,
// 		Price:       p.Price,
// 		Description: p.Description,
// 		ImgUrls:     p.ImgUrls,
// 	}
// }
//
func (c *ProductConvertor) check(p any) {
	if p == nil {
		panic("cannot convert nil product")
	}
}
