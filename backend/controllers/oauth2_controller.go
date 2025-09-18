package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"
	"xuan-ke-tong/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type OAuth2State struct {
	State     string    `json:"state"`
	ExpiresAt time.Time `json:"expiresAt"`
}

var oauth2States = make(map[string]OAuth2State)

// OAuth2配置结构
type OAuth2Config struct {
	AppID       string
	AppSecret   string
	TokenURL    string
	UserInfoURL string
}

// 从.env文件加载OAuth2配置
func loadOAuth2Config() (*OAuth2Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("加载环境配置文件失败: %v", err)
	}

	// 从环境变量中读取配置
	config := &OAuth2Config{
		AppID:       os.Getenv("OAUTH2_MARKET_APP_ID"),
		AppSecret:   os.Getenv("OAUTH2_MARKET_APP_SECRET"),
		TokenURL:    os.Getenv("OAUTH2_MARKET_TOKEN_URL"),
		UserInfoURL: os.Getenv("OAUTH2_MARKET_USERINFO_URL"),
	}

	return config, nil
}

// 集市授权的用户信息结构
type SSEUserInfo struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Intro     string `json:"intro"`
	Identity  string `json:"identity"`
}

type SSEUserInfoResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data SSEUserInfo `json:"data"`
}

// 集市授权的Token响应结构
type SSETokenResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	} `json:"data"`
}

func GenerateOAuth2State(c *gin.Context) {
	stateBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, stateBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成state失败",
		})
		return
	}

	state := hex.EncodeToString(stateBytes)

	// 存储state，有效期5分钟
	oauth2States[state] = OAuth2State{
		State:     state,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	c.JSON(http.StatusOK, gin.H{
		"state": state,
	})
}

// OAuth2Callback 处理OAuth2回调
func OAuth2Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少必要参数",
		})
		return
	}

	// 验证state
	storedState, exists := oauth2States[state]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的state",
		})
		return
	}

	// 检查state是否过期
	if time.Now().After(storedState.ExpiresAt) {
		delete(oauth2States, state)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "state已过期",
		})
		return
	}

	// 清理已使用的state
	delete(oauth2States, state)

	// 使用code获取access_token
	accessToken, err := getAccessToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取access_token出错: " + err.Error(),
		})
		return
	}

	// 使用access_token获取用户信息
	userInfo, err := getUserInfo(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户信息失败: " + err.Error(),
		})
		return
	}

	// 检查用户是否已存在
	var user models.User
	err = config.DB.Where("email = ?", userInfo.Email).First(&user).Error

	if err != nil {
		// 用户不存在，创建新用户
		user, err = createUserFromSSE(userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "创建用户失败: " + err.Error(),
			})
			return
		}
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成token失败",
		})
		return
	}

	// 返回成功响应，包含token和用户信息
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"nickname":  user.Nickname,
			"avatar":    user.Avatar,
			"role":      user.Role,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
		"message": "OAuth2登录成功",
	})
}

// getAccessToken 从SSE集市获取access_token
func getAccessToken(code string) (string, error) {
	// 加载OAuth2配置
	config, err := loadOAuth2Config()
	if err != nil {
		return "", fmt.Errorf("加载OAuth2配置失败: %v", err)
	}

	// 构建form-data请求参数
	formData := map[string][]string{
		"code":       {code},
		"app_id":     {config.AppID},
		"app_secret": {config.AppSecret},
	}

	// 发送POST请求获取token
	resp, err := http.PostForm(config.TokenURL, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 解析响应
	var tokenResp SSETokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	// 检查响应状态
	if tokenResp.Code != 200 {
		return "", fmt.Errorf("获取access_token出错: %s", tokenResp.Msg)
	}

	return tokenResp.Data.AccessToken, nil
}

// getUserInfo 从SSE集市获取用户信息
func getUserInfo(accessToken string) (*SSEUserInfo, error) {
	// 加载OAuth2配置
	config, err := loadOAuth2Config()
	if err != nil {
		return nil, fmt.Errorf("加载OAuth2配置失败: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置Authorization头
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析响应
	var userInfoResp SSEUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfoResp); err != nil {
		return nil, err
	}

	// 检查响应状态
	if userInfoResp.Code != 200 {
		return nil, fmt.Errorf("获取用户信息出错: %s", userInfoResp.Msg)
	}

	return &userInfoResp.Data, nil
}

// createUserFromSSE 根据SSE用户信息创建本地用户
func createUserFromSSE(sseUser *SSEUserInfo) (models.User, error) {
	// 生成复杂的默认密码
	defaultPassword := generateSecurePassword()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	// 生成用户名（基于邮箱前缀）
	username := sseUser.Email
	if len(username) > 50 {
		username = username[:50]
	}

	// 确保用户名唯一
	originalUsername := username
	counter := 1
	for {
		var existingUser models.User
		err := config.DB.Where("username = ?", username).First(&existingUser).Error
		if err != nil {
			break // 用户名不存在，可以使用
		}
		username = fmt.Sprintf("%s_%d", originalUsername, counter)
		counter++
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    sseUser.Email,
		Nickname: sseUser.Name,
		Avatar:   sseUser.AvatarURL,
		Role:     "user",
	}

	// 根据identity设置角色
	if sseUser.Identity == "teacher" {
		user.Role = "teacher"
	} else if sseUser.Identity == "organization" {
		user.Role = "user"
	} else {
		user.Role = "user" // 默认为普通用户
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

// generateSecurePassword 生成安全的随机密码
func generateSecurePassword() string {
	// 定义字符集
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	// 组合所有字符
	allChars := lowercase + uppercase + digits + symbols

	// 密码长度（16-20位）
	passwordLength := 16 + randInt(5)

	var password strings.Builder

	// 确保密码包含每种类型的字符
	password.WriteByte(lowercase[randInt(len(lowercase))])
	password.WriteByte(uppercase[randInt(len(uppercase))])
	password.WriteByte(digits[randInt(len(digits))])
	password.WriteByte(symbols[randInt(len(symbols))])

	// 填充剩余长度
	for i := 4; i < passwordLength; i++ {
		password.WriteByte(allChars[randInt(len(allChars))])
	}

	// 打乱密码字符顺序
	passwordStr := password.String()
	runes := []rune(passwordStr)

	for i := len(runes) - 1; i > 0; i-- {
		j := randInt(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// randInt 生成指定范围内的随机整数
func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}
