package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/centrypoint/refrigerator/back/go/main/db"
)

// ProductsAPI type
type ProductsAPI struct{}

// GetProductsHandler func
func (p *ProductsAPI) GetProductsHandler(
	apiInterface APIInterface,
	dbi DatabaseInterface,
	w http.ResponseWriter,
	r *http.Request) {

	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		w.WriteHeader(400)
		return
	}
	products, err := apiInterface.GetProducts(&http.Client{}, name)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	response, e := json.Marshal(products)
	if e != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(response)

	// Concurrent add products to database
	for _, v := range products {
		go dbi.InsertProductToDB(v)
	}
}

// GetProductImageHandler func
func (p *ProductsAPI) GetProductImageHandler(
	apiInterface APIInterface,
	dbi DatabaseInterface,
	w http.ResponseWriter,
	r *http.Request) {

	productID := r.URL.Query().Get("product_id")
	side := r.URL.Query().Get("side")
	if len(productID) < 1 {
		w.WriteHeader(400)
		return
	}

	if len(side) < 1 {
		side = "100"
	}

	// Fetch from database
	image, err := dbi.FetchPhotoFromDBByProductID(productID, side)
	if len(image.Data) < 10 || err != nil {
		image, err := apiInterface.GetProductImageFromAPI(&http.Client{}, side, productID)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		response, e := json.Marshal(map[string]string{"result": image})
		if e != nil {
			w.WriteHeader(500)
			return
		}

		w.Write(response)

		numProductID, e := strconv.Atoi(productID)

		if e == nil {
			numSide, e := strconv.Atoi(side)
			if e == nil {
				// Concurrent add photo to db
				go dbi.InsertPhotoToDB(db.Photo{
					ProductID: numProductID,
					Side:      numSide,
					Data:      image,
				})
			}
		}

	} else {
		response, e := json.Marshal(map[string]string{"result": image.Data})
		if e != nil {
			w.WriteHeader(500)
			return
		}

		w.Write(response)
	}

}

// GetProductShelflifeHandler func
func (p *ProductsAPI) GetProductShelflifeHandler(
	apiInterface APIInterface,
	dbi DatabaseInterface,
	w http.ResponseWriter,
	r *http.Request) {

	var productID = r.URL.Query().Get("product_id")
	var result int
	var response []byte
	var err error
	var shelf db.Shelflife

	if len(productID) < 1 {
		w.WriteHeader(400)
		return
	}

	if shelf, err = dbi.FetchShelflifeFromDBByProductID(productID); err != nil {
		var product db.Product
		err = nil
		if product, err = dbi.FetchProductFromDBByID(productID); err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		if result, err = apiInterface.GetProductShelfLife(&http.Client{}, product.URL); err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		if response, err = json.Marshal(map[string]int{"result": result}); err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		if _, err = w.Write(response); err != nil {
			log.Println(err)
			return
		}

		go dbi.InsertShelflifeToDB(db.Shelflife{
			ProductID: product.ProductID,
			Data:      result,
		})
	} else {
		if response, err = json.Marshal(map[string]int{"result": shelf.Data}); err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		if _, err = w.Write(response); err != nil {
			log.Println(err)
			return
		}
	}
}
