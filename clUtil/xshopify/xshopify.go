package xshopify

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	xhr2 "github.com/cxi7448/cxhttp/clUtil/xhr"
)

// const SHOPIFY_HOST = "https://your-store.myshopify.com"
const SHOPIFY_HOST = "https://your-development-store.myshopify.com"

type ShopiFy struct {
	AccessToken string
}

type Product struct {
	Title       string
	BodyHtml    string
	Vendor      string
	ProductType string
	Status      string
}

func NewShopiFy(access string) *ShopiFy {
	return &ShopiFy{
		AccessToken: access,
	}
}

func (this *ShopiFy) CreateProduct(product *Product) {
	xhr := xhr2.NewXhr(this.buildUrl("/admin/api/2024-07/products.json"))
	result := clJson.M{}
	data := clJson.M{
		"product": clJson.M{
			"title":        product.Title,
			"body_html":    product.BodyHtml,
			"vendor":       product.Vendor,
			"product_type": product.ProductType,
			"status":       product.Status,
		},
	}
	xhr.SetJSON()
	xhr.SetHeaders(clJson.M{
		"x-shopify-access-token": this.AccessToken,
	})
	err := xhr.Post(data, &result)
	fmt.Println(err)
}
func (this *ShopiFy) buildUrl(uri string) string {
	return fmt.Sprintf("%v%v", SHOPIFY_HOST, uri)
}

func NewProduct(title string) *Product {
	return &Product{Title: title}
}

func (this *Product) SetProductType(product_type string) {
	this.ProductType = product_type
}

func (this *Product) SeStatus(status string) {
	this.Status = status
}
func (this *Product) SetBody(body_html string) {
	this.BodyHtml = body_html
}

func (this *Product) SetVendor(vendor string) {
	this.Vendor = vendor
}

//
//import requests
//import json
//
//# 配置您的Shopify和Facebook信息
//SHOPIFY_STORE_URL = 'https://your-store.myshopify.com/'
//SHOPIFY_ACCESS_TOKEN = 'your_shopify_access_token'
//FACEBOOK_ACCESS_TOKEN = 'your_facebook_access_token'
//FACEBOOK_GRAPH_URL = 'https://graph.facebook.com/'
//
//# 定义将要发送到Shopify的商品数据
//shopify_product_data = {
//"product": {
//"title": "New Product",
//"body_html": "<strong>Good stuff</strong>",
//"vendor": "Your Company",
//# ... 其他必要的商品属性
//}
//}
//
//# 定义将要发送到Facebook的商品数据
//facebook_product_data = {
//"message": "New Product Available",
//"link": "https://your-store.myshopify.com/products/new-product",
//# ... 其他自定义分享内容
//}
//
//# 将商品推送到Shopify
//headers = {
//'Content-Type': 'application/json',
//'X-Shopify-Access-Token': SHOPIFY_ACCESS_TOKEN,
//}
//response = requests.post(f'{SHOPIFY_STORE_URL}/admin/api/2021-04/products.json',
//headers=headers, data=json.dumps(shopify_product_data))
//
//# 将商品推送到Facebook
//response = requests.post(f'{FACEBOOK_GRAPH_URL}/me/feed',
//params={'access_token': FACEBOOK_ACCESS_TOKEN},
//data=json.dumps(facebook_product_data))
//
//# 打印响应
//print(response.json())
