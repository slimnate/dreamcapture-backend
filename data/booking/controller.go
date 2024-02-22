package booking

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	repo *Repository
}

func NewBookingController(repo *Repository) *BookingController {
	return &BookingController{
		repo: repo,
	}
}

func (c *BookingController) List(ctx *gin.Context) {
	bookings, err := c.repo.All()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	ctx.JSON(http.StatusOK, bookings)
}
