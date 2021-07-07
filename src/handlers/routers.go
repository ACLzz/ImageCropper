package handlers

import "github.com/gin-gonic/gin"

func MainRouter(r *gin.Engine) {
	r.MaxMultipartMemory = 16 << 20 				// 16mb
	r.POST("/upload_image", UploadImage)
}
