package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	dsn := "root:{ROOT-PASSWORD}@tcp({HOST-IPADDRESS}:3306)/{MySQL-DATABASE}?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	fmt.Println(db)
	fmt.Println(err)

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	type User struct {
		//gorm.Model
		Username string //`gorm:"type:varchar(20); not null" json:"username" binding:"required"`
		Password string //`gorm:"type:varchar(20); not null" json:"password" binding:"required"`
	}

	//db.AutoMigrate(&User{})

	r := gin.Default()

	//增
	r.POST("/user/add", func(c *gin.Context) {
		var data User

		err := c.ShouldBindJSON(&data)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "faild",
				"data": gin.H{},
				"code": 404,
			})
		} else {
			//db.Create(&data)

			c.JSON(http.StatusOK, gin.H{
				"msg":  "success",
				"data": data,
				"code": 200,
			})
		}
		find := data.Username + data.Password
		fmt.Println(find)
	})

	//删
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []User

		id := c.Param("id")
		db.Where("id = ?", id).Find(&data)

		if len(data) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "faild",
				"code": 404,
			})
		} else {
			db.Delete(&data)
			c.JSON(http.StatusOK, gin.H{
				"msg":  "success",
				"code": 200,
			})
		}
	})

	//改
	r.PUT("/user/update/:id", func(c *gin.Context) {
		var data User

		id := c.Param("id")
		db.Select("id").Where("id = ?", id).Find(&data)

		if data.ID == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "faild",
				"code": 404,
			})
		} else {
			err := c.ShouldBindJSON(&data)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg":  "faild",
					"code": 404,
				})
			} else {
				db.Where("id = ?", id).Updates(&data)
				c.JSON(http.StatusOK, gin.H{
					"msg":  "success",
					"code": 200,
				})
			}
		}
	})

	//查
	r.GET("/user/find/:username", func(c *gin.Context) {
		name := c.Param("username")
		var datalist []User

		db.Where("username = ?", name).Find(&datalist)

		if len(datalist) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "faild",
				"data": gin.H{},
				"code": 404,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "success",
				"data": datalist,
				"code": 200,
			})
		}
	})

	r.Run(":80")
}
