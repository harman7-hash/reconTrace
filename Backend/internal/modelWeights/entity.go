package modelweights

import (
	"mime/multipart"
)

type ModelUploadForm struct {
	UserId string `form:"user_id" binding:"required"`
	Weight *multipart.FileHeader `form:"weight" binding:"required"`
}
