package concurrency

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"sync"
// 	"time"

// 	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
// 	"github.com/Nishad4140/SkillSync_UserService/internal/adapters"
// 	"github.com/Nishad4140/SkillSync_UserService/internal/service"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type Concurrency struct {
// 	DB       *gorm.DB
// 	adapters adapters.AdapterInterface
// 	mu       sync.Mutex
// 	service  *service.UserService
// }

// func NewConcurrency(DB *gorm.DB, adapters adapters.AdapterInterface, service *service.UserService) *Concurrency {
// 	return &Concurrency{
// 		DB:       DB,
// 		adapters: adapters,
// 		service:  service,
// 	}
// }

// type Users struct {
// 	ClientId     uuid.UUID
// 	FreelancerId uuid.UUID
// }

// func (c *Concurrency) Concurrency() {
// 	ticker := time.NewTicker(1 * time.Minute)
// 	go func() {
// 		for range ticker.C {
// 			c.mu.Lock()
// 			if err := c.DB.Exec(`
// 			UPDATE users SET is_blocked = true WHERE report_count > 50
// 			`).Error; err != nil {
// 				log.Print("error while performing concurrency", err)
// 				break
// 			}

// 			currentDate := time.Now()

// 			threeDaysLater := currentDate.AddDate(0, 0, 3)

// 			dateString := threeDaysLater.Format("2006-01-02")

// 			var users Users

// 			if err := c.DB.Raw("SELECT client_id, freelancer_id FROM projects WHERE end_date = $1", dateString).Scan(&users).Error; err != nil {
// 				log.Print("error while performing concurrency", err)
// 				break
// 			}

// 			if _, err := c.service.NotiClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
// 				UserId:       users.ClientId.String(),
// 				Notification: fmt.Sprintf(`{"message":"Please be aware of your interview scheduled at %v by the company %s for the position %s with roomId %s"}`, shortlist.InterviewDate.String(), jobData.Company, jobData.Designation, shortlist.RoomId),
// 			}); err != nil {
// 				return err
// 			}
// 		}
// 		c.mu.Unlock()
// 	}()
// }
