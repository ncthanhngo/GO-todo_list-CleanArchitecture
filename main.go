package main

import (
	"GOEND/common"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	//GORM
	dsn := os.Getenv("DB_CONN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug()
	log.Println("Connected to database", db)

	//GIN START

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", CreateItem(db))
			items.GET("", ListItem(db))
			items.GET("/:id", GetItem(db))
			items.PATCH("/:id", UpdateItem(db))
			items.DELETE(":id", DeleteItem(db))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":3000") // listen and serve on 0.0.0.0:3000 (for windows "localhost:3000")
}

// createItem
func CreateItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var itemData TodoItemCreation
		if err := c.ShouldBind(&itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Lưu dữ liệu vào cơ sở dữ liệu
		if err := db.Create(&itemData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Trả về phản hồi với ID của mục vừa được tạo
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData.Id))
	}
}

// Get Iterm by ID
func GetItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var itemData TodoItem
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Lay dữ liệu vào cơ sở dữ liệu
		if err := db.Where("id=?", id).First(&itemData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Trả về phản hồi với ID của mục vừa được tạo
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData))
	}
}

// Update Iterm
func UpdateItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var updateData TodoItemUpdate
		if err := c.ShouldBind(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Lưu dữ liệu vào cơ sở dữ liệu
		if err := db.Where("id=?", id).Updates(&updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Trả về phản hồi true (theo nguyen tac SOLID)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

//Delete by ID

func DeleteItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// hard delete
		//if err := db.Table(TodoItem{}.TableName()).Where("id=?", id).Delete(nil).Error; err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}
		//soft delete
		// cách 1:
		deletedStatus := "Deleted"
		if err := db.Where("id=?", id).Updates(&TodoItemUpdate{Status: &deletedStatus}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		//cach 2: khai bao mot Map co key la String, value la any
		//if err := db.Where("id=?", id).Updates(map[string]interface{}{
		//	"status": "Deleted",}).Error; err != nil

		// Trả về phản hồi true (theo nguyen tac SOLID)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

//weak vs strong entity
//Strong: Quan trong
//- nhieu khoa ngoai tham chieu toi
//Weak: bang giua, thuong tao do quan he n-n

// List Iterm
func ListItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging
		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		paging.Process()
		var result []TodoItem
		if err := db.Table(TodoItem{}.TableName()).Select("id").Count(&paging.Total).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		} //dem theo ID
		if err := db.Table(TodoItem{}.TableName()).
			//Count(&paging.Total). > count(*) khong hieu qua doi voi các bang nhieu cột
			Offset((paging.Page - 1) * paging.Limit).
			Limit(paging.Limit).
			Order("id desc").Find(&result).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Trả về phản hồi true (theo nguyen tac SOLID)
		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, nil))
	}
}
