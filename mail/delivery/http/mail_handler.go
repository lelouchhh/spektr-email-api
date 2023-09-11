package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"spektr-email-api/domain"
)

type MailHandler struct {
	MUsecase domain.MailUsecase
}

func NewMailHandler(g *gin.Engine, us domain.MailUsecase) {
	handler := &MailHandler{
		MUsecase: us,
	}

	g.POST("/feedback", handler.Feedback)
}

func (m *MailHandler) Feedback(c *gin.Context) {
	var mail domain.Mail
	err := c.ShouldBindJSON(&mail)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, "error")
		return
	}
	ctx := c.Request.Context()
	err = m.MUsecase.Feedback(ctx, mail)
	if err != nil {
		fmt.Println(err)

		c.JSON(500, "error")
		return
	}
	c.JSON(200, "ok")
}
