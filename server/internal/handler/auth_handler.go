package handler

import (
	"log"
	"net/http"

	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"server/internal/models"
)

// AuthHandlerは認証関連のAPIハンドラをまとめます。
// この構造体は、ハンドラが依存する外部リソース（DBやFirebase認証クライアント）を保持します。
type AuthHandler struct {
	DB           *gorm.DB
	FirebaseAuth *auth.Client
}

// NewAuthHandlerは新しいAuthHandlerのインスタンスを作成します。
// データベースクライアントとFirebase認証クライアントを引数として受け取り、ハンドラに注入します。
// Firebase Authクライアントの初期化は、メインアプリケーションの起動時に一度だけ行われるべきです。
func NewAuthHandler(db *gorm.DB, firebaseAuthClient *auth.Client) *AuthHandler {
	return &AuthHandler{DB: db, FirebaseAuth: firebaseAuthClient}
}

// RegisterRequestは/registerエンドポイントへのリクエストボディの構造体です。
// クライアントから送られるFirebase IDトークンと、アプリケーション固有のユーザー情報を含みます。
type RegisterRequest struct {
	IDToken         string `json:"id_token" binding:"required"`
	LoginID         string `json:"login_id" binding:"required,alphanum"`
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
	EMail           string `json:"email" binding:"required,email"`
}

// Registerは新しいユーザーを登録します (POST /register)。
// クライアントから送られてきたFirebase IDトークンを検証し、その後にユーザー情報をデータベースに保存します。
// このエンドポイント自体に認証ミドルウェアは適用されません。
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	// リクエストボディのJSONをRegisterRequest構造体にバインドし、バリデーションを実行します。
	// `binding`タグに基づいて自動的に必須チェックや形式チェックが行われます。
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: Invalid request body for /register: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- 1. Firebase IDトークンの検証 ---
	// クライアントから送られたIDトークンを検証し、Firebaseの正当なユーザーであることを確認します。
	// 検証が成功すると、そのユーザーのFirebase UIDを安全に抽出できます。
	token, err := h.FirebaseAuth.VerifyIDToken(c.Request.Context(), req.IDToken)
	if err != nil {
		log.Printf("WARN: Firebase ID token verification failed for /register: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired Firebase ID token"})
		return
	}
	firebaseUID := token.UID

	// --- 2. 既存ユーザーの重複チェック ---
	// 同じFirebase UIDまたはLoginIDを持つユーザーが既にデータベースに存在しないか確認し、重複登録を防ぎます。
	var existingUser models.User
	// FirebaseUID で既に存在するかチェック
	if h.DB.Where("firebase_uid = ?", firebaseUID).First(&existingUser).Error == nil {
		log.Printf("WARN: Attempted to register existing Firebase UID: %s", firebaseUID)
		c.JSON(http.StatusConflict, gin.H{"error": "User with this Firebase UID already exists"})
		return
	}
	// LoginID で既に存在するかチェック
	if h.DB.Where("login_id = ?", req.LoginID).First(&existingUser).Error == nil {
		log.Printf("WARN: Attempted to register existing LoginID: %s", req.LoginID)
		c.JSON(http.StatusConflict, gin.H{"error": "LoginID already taken"})
		return
	}

	// --- 3. 新規ユーザーの作成とデータベース保存 ---
	// 検証済みの情報とリクエストデータを使って新しいUserモデルを作成します。
	newUser := models.User{
		FirebaseUID:     firebaseUID,
		LoginID:         req.LoginID,
		DisplayName:     req.DisplayName,
		ProfileImageURL: req.ProfileImageURL,
		EMail:           req.EMail,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// ユーザー情報をデータベースに保存します。
	if result := h.DB.Create(&newUser); result.Error != nil {
		log.Printf("ERROR: Failed to create user in DB: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user due to database error"})
		return
	}

	// --- 4. 登録成功レスポンス ---
	// 登録が成功したことをクライアントに通知します。
	c.JSON(http.StatusCreated, gin.H{
		"message":      "User registered successfully",
		"user_id":      newUser.ID,
		"firebase_uid": newUser.FirebaseUID,
	})
}
