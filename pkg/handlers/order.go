package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/tsagita/go-vue-postgre/pkg/database"
	"github.com/tsagita/go-vue-postgre/pkg/domain"
	// "encoding/json"
)

var db *sql.DB

func init() {
	var err error
	db, err = database.Connection()
	if err != nil {
		fmt.Printf("Error")
	}
}

func List(c *gin.Context) {
	productName := c.Query("product")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	items, totalData, err := getOrderItems(productName, startDate, endDate, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "status": "ok", "data": items, "totalData": totalData})
}

func getOrderItems(productName, startDate, endDate string, page int) ([]domain.OrderItem, int, error) {
	query := `
		SELECT
			oi.order_id,
			oi.product,
			o.order_name,
			o.created_at,
			oi.price_per_unit as delivery_amount,
			oi.price_per_unit * oi.quantity as total_amount,
			c.name,
			cc.company_name
		FROM order_items oi
		LEFT JOIN orders o ON o.id = oi.order_id
		LEFT JOIN customers c ON c.user_id = o.customer_id
		LEFT JOIN customer_companies cc ON cc.id = c.company_id
		WHERE o.created_at IS NOT NULL
	`

	args := []interface{}{}

	if productName != "" {
		query += " AND oi.product ILIKE $1"
		args = append(args, "%"+productName+"%")
	}

	if startDate != "" {
		query += " AND o.created_at >= $1"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND o.created_at <= $1"
		args = append(args, endDate)
	}

	// Count total data without pagination
	totalDataQuery := "SELECT COUNT(*) FROM (" + query + ") as total"
	var totalData int
	if err := db.QueryRow(totalDataQuery, args...).Scan(&totalData); err != nil {
		return nil, 0, err
	}

	// Add pagination to the main query
	offset := (page - 1) * 5
	query += fmt.Sprintf(" ORDER BY o.created_at DESC LIMIT 5 OFFSET $%d", len(args)+1)
	args = append(args, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var result []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		var deliveryAmt sql.NullFloat64
		var TotalAmount sql.NullFloat64
		if err := rows.Scan(&item.OrderID, &item.Product, &item.OrderName, &item.CreatedAt, &deliveryAmt, &item.TotalAmount, &item.CustomerName, &item.CompanyName); err != nil {
			return nil, 0, err
		}

		if deliveryAmt.Valid {
			item.DeliveryAmt = deliveryAmt
		} else {
			item.DeliveryAmt = sql.NullFloat64{Valid: false}
		}

		if TotalAmount.Valid {
			item.TotalAmount = TotalAmount
		} else {
			item.TotalAmount = sql.NullFloat64{Valid: false}
		}

		result = append(result, item)
	}

	return result, totalData, nil
}
