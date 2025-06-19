// controller/webhook.go
package controller

import (
	"log"
	"net/http"
	"os"
	"qiscus-agent-allocator/model"
	"qiscus-agent-allocator/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebhookPayload struct {
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func WebhookHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload WebhookPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		// Simpan ke database
		customer := model.Customer{
			Name:      payload.Name,
			RoomID:    payload.RoomID,
			Status:    "waiting",
			CreatedAt: time.Now(),
		}

		secretKey := os.Getenv("QISCUS_SECRET_KEY")
		appID := os.Getenv("QISCUS_APP_ID")

		valid, err := utils.ValidateRoomID(payload.RoomID, secretKey, appID)
		if err != nil || !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id tidak valid"})
			return
		}

		
		if err := db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save customer"})
			return
		}

		queue := model.Queue{
			CustomerID: customer.ID,
			RoomID:     payload.RoomID,
			Assigned:   false,
			CreatedAt:  time.Now(),
		}
		db.Create(&queue)

		// Langsung assign

		agents, err := utils.GetAvailableAgents(secretKey, appID, 2)
		if err != nil || len(agents) == 0 {
			log.Println("Tidak ada agent tersedia atau error:", err)
			return
		}

		err = utils.AssignAgentToRoom(payload.RoomID, agents[0].ID, secretKey, appID)
		if err != nil {
			log.Println("Gagal assign agent:", err)
			return
		}

		db.Model(&queue).Update("assigned", true)
		db.Model(&customer).Update("status", "assigned")

		c.JSON(http.StatusOK, gin.H{"message": "Agent assigned"})
	}
}
