package controllers

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// FakeActivityRequest 假数据请求参数
type FakeActivityRequest struct {
	Count int `json:"count" binding:"min=1,max=50"` // 返回数据条数，默认10条，最大50条
}

type FakeRealtimeActivity struct {
	UID    string  `json:"uid"`
	Time   string  `json:"time"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

func generateUID() string {
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

// maskUID 对UID进行脱敏处理
func maskUID(uid string) string {
	if len(uid) < 5 {
		return uid // 如果UID太短，直接返回
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
	var req FakeActivityRequest

	// 解析请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 设置默认值和限制
	if req.Count <= 0 || req.Count > 50 {
		req.Count = 10
	}

	rand.Seed(time.Now().UnixNano())
	var list []FakeRealtimeActivity
	for i := 0; i < req.Count; i++ {
		list = append(list, generateFakeActivity())
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    list,
	})
}
