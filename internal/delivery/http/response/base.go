package response

import (
	"time"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/google/uuid"
)

// Response represents a Response body format
type Response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// NewResponse is a helper function to create a response body
func NewResponse(success bool, message string, data any) Response {
	return Response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// Meta represents metadata for a paginated response
type Meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit uint64 `json:"limit" example:"10"`
	Skip  uint64 `json:"skip" example:"0"`
}

// NewMeta is a helper function to create metadata for a paginated response
func NewMeta(total, limit, skip uint64) Meta {
	return Meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

// AuthResponse represents an authentication response body
type AuthResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

// NewAuthResponse is a helper function to create a response body for handling authentication data
func NewAuthResponse(token string) AuthResponse {
	return AuthResponse{
		AccessToken: token,
	}
}

// UserResponse represents a user response body
type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// NewUserResponse is a helper function to create a response body for handling user data
func NewUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// CollectionResponse represents a collection response body
type CollectionResponse struct {
	ID          uuid.UUID `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
	Name        string    `json:"name" example:"My Collection"`
	Description string    `json:"description,omitempty" example:"Collection description"`
	CreatedAt   time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// NewCollectionResponse is a helper function to create a response body for handling collection data
func NewCollectionResponse(collection *domain.Collection) CollectionResponse {
	return CollectionResponse{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		CreatedAt:   collection.CreatedAt,
		UpdatedAt:   collection.UpdatedAt,
	}
}

// SecretResponse represents a secret response body
type SecretResponse struct {
	ID           uuid.UUID             `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
	CollectionID uuid.UUID             `json:"collection_id" example:"fab8dfe9-7cd0-4cd7-a387-7d6835a910d3"`
	SecretType   domain.SecretTypeEnum `json:"secret_type" example:"password"`
	CreatedAt    time.Time             `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt    time.Time             `json:"updated_at" example:"1970-01-01T00:00:00Z"`
	CreatedBy    uuid.UUID             `json:"created_by" example:"f10ff052-b316-47f0-9788-ae8ebfa91b86"`
	UpdatedBy    uuid.UUID             `json:"updated_by" example:"f10ff052-b316-47f0-9788-ae8ebfa91b86"`
}

// NewSecretResponse is a helper function to create a response body for handling secret data
func NewSecretResponse(secret *domain.Secret) SecretResponse {
	return SecretResponse{
		ID:           secret.ID,
		CollectionID: secret.CollectionID,
		SecretType:   secret.SecretType,
		CreatedAt:    secret.CreatedAt,
		UpdatedAt:    secret.UpdatedAt,
		CreatedBy:    secret.CreatedBy,
		UpdatedBy:    secret.UpdatedBy,
	}
}

// // paymentResponse represents a payment response body
// type paymentResponse struct {
// 	ID   uint64             `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
// 	Name string             `json:"name" example:"Tunai"`
// 	Type domain.PaymentType `json:"type" example:"CASH"`
// 	Logo string             `json:"logo" example:"https://example.com/cash.png"`
// }

// // newPaymentResponse is a helper function to create a response body for handling payment data
// func newPaymentResponse(payment *domain.Payment) paymentResponse {
// 	return paymentResponse{
// 		ID:   payment.ID,
// 		Name: payment.Name,
// 		Type: payment.Type,
// 		Logo: payment.Logo,
// 	}
// }

// // categoryResponse represents a category response body
// type categoryResponse struct {
// 	ID   uint64 `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
// 	Name string `json:"name" example:"Foods"`
// }

// // newCategoryResponse is a helper function to create a response body for handling category data
// func newCategoryResponse(category *domain.Category) categoryResponse {
// 	return categoryResponse{
// 		ID:   category.ID,
// 		Name: category.Name,
// 	}
// }

// // productResponse represents a product response body
// type productResponse struct {
// 	ID        uint64           `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
// 	SKU       string           `json:"sku" example:"9a4c25d3-9786-492c-b084-85cb75c1ee3e"`
// 	Name      string           `json:"name" example:"Chiki Ball"`
// 	Stock     int64            `json:"stock" example:"100"`
// 	Price     float64          `json:"price" example:"5000"`
// 	Image     string           `json:"image" example:"https://example.com/chiki-ball.png"`
// 	Category  categoryResponse `json:"category"`
// 	CreatedAt time.Time        `json:"created_at" example:"1970-01-01T00:00:00Z"`
// 	UpdatedAt time.Time        `json:"updated_at" example:"1970-01-01T00:00:00Z"`
// }

// // newProductResponse is a helper function to create a response body for handling product data
// func newProductResponse(product *domain.Product) productResponse {
// 	return productResponse{
// 		ID:        product.ID,
// 		SKU:       product.SKU.String(),
// 		Name:      product.Name,
// 		Stock:     product.Stock,
// 		Price:     product.Price,
// 		Image:     product.Image,
// 		Category:  newCategoryResponse(product.Category),
// 		CreatedAt: product.CreatedAt,
// 		UpdatedAt: product.UpdatedAt,
// 	}
// }

// // orderResponse represents an order response body
// type orderResponse struct {
// 	ID           uint64                 `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
// 	UserID       uint64                 `json:"user_id" example:"1"`
// 	PaymentID    uint64                 `json:"payment_type_id" example:"1"`
// 	CustomerName string                 `json:"customer_name" example:"John Doe"`
// 	TotalPrice   float64                `json:"total_price" example:"100000"`
// 	TotalPaid    float64                `json:"total_paid" example:"100000"`
// 	TotalReturn  float64                `json:"total_return" example:"0"`
// 	ReceiptCode  string                 `json:"receipt_id" example:"4979cf6e-d215-4ff8-9d0d-b3e99bcc7750"`
// 	Products     []orderProductResponse `json:"products"`
// 	PaymentType  paymentResponse        `json:"payment_type"`
// 	CreatedAt    time.Time              `json:"created_at" example:"1970-01-01T00:00:00Z"`
// 	UpdatedAt    time.Time              `json:"updated_at" example:"1970-01-01T00:00:00Z"`
// }

// // newOrderResponse is a helper function to create a response body for handling order data
// func newOrderResponse(order *domain.Order) orderResponse {
// 	return orderResponse{
// 		ID:           order.ID,
// 		UserID:       order.UserID,
// 		PaymentID:    order.PaymentID,
// 		CustomerName: order.CustomerName,
// 		TotalPrice:   order.TotalPrice,
// 		TotalPaid:    order.TotalPaid,
// 		TotalReturn:  order.TotalReturn,
// 		ReceiptCode:  order.ReceiptCode.String(),
// 		Products:     newOrderProductResponse(order.Products),
// 		PaymentType:  newPaymentResponse(order.Payment),
// 		CreatedAt:    order.CreatedAt,
// 		UpdatedAt:    order.UpdatedAt,
// 	}
// }

// // orderProductResponse represents an order product response body
// type orderProductResponse struct {
// 	ID               uint64          `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
// 	OrderID          uint64          `json:"order_id" example:"1"`
// 	ProductID        uint64          `json:"product_id" example:"1"`
// 	Quantity         int64           `json:"qty" example:"1"`
// 	Price            float64         `json:"price" example:"100000"`
// 	TotalNormalPrice float64         `json:"total_normal_price" example:"100000"`
// 	TotalFinalPrice  float64         `json:"total_final_price" example:"100000"`
// 	Product          productResponse `json:"product"`
// 	CreatedAt        time.Time       `json:"created_at" example:"1970-01-01T00:00:00Z"`
// 	UpdatedAt        time.Time       `json:"updated_at" example:"1970-01-01T00:00:00Z"`
// }

// // newOrderProductResponse is a helper function to create a response body for handling order product data
// func newOrderProductResponse(orderProduct []domain.OrderProduct) []orderProductResponse {
// 	var orderProductResponses []orderProductResponse

// 	for _, orderProduct := range orderProduct {
// 		orderProductResponses = append(orderProductResponses, orderProductResponse{
// 			ID:               orderProduct.ID,
// 			OrderID:          orderProduct.OrderID,
// 			ProductID:        orderProduct.ProductID,
// 			Quantity:         orderProduct.Quantity,
// 			Price:            orderProduct.Product.Price,
// 			TotalNormalPrice: orderProduct.TotalPrice,
// 			TotalFinalPrice:  orderProduct.TotalPrice,
// 			Product:          newProductResponse(orderProduct.Product),
// 			CreatedAt:        orderProduct.CreatedAt,
// 			UpdatedAt:        orderProduct.UpdatedAt,
// 		})
// 	}

// 	return orderProductResponses
// }
