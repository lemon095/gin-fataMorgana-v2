package controllers

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FakeRealtimeActivity struct {
	UID    string  `json:"uid"`
	Time   string  `json:"time"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

func generateUID() string {
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

func maskUID(uid string) string {
	if len(uid) != 8 {
		return uid
	}
	return uid[:2] + "***" + uid[5:]
}

func randomTime() string {
	hour := rand.Intn(24)
	min := rand.Intn(60)
	return fmt.Sprintf("%02d:%02d", hour, min)
}

func randomAmount() float64 {
	return math.Round((rand.Float64()*800+200)*100) / 100 // 200~1000
}
func randomType() string {
	types := []string{"点赞", "关注", "收藏", "转发"}
	return types[rand.Intn(len(types))]
}

func generateFakeActivity() FakeRealtimeActivity {
	uid := generateUID()
	return FakeRealtimeActivity{
		UID:    maskUID(uid),
		Time:   randomTime(),
		Amount: randomAmount(),
		Type:   randomType(),
	}
}

// GetFakeRealtimeActivities 假数据实时动态接口
func GetFakeRealtimeActivities(c *gin.Context) {
	n := 10 // 默认生成10条
	if nParam := c.Query("n"); nParam != "" {
		if v, err := strconv.Atoi(nParam); err == nil && v > 0 && v <= 50 {
			n = v
		}
	}
	rand.Seed(time.Now().UnixNano())
	var list []FakeRealtimeActivity
	for i := 0; i < n; i++ {
		list = append(list, generateFakeActivity())
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    list,
	})
}
