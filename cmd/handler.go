package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

// ProfileHandler struct
type profileHandler struct {
	rd AuthInterface
	tk TokenInterface
}

func NewProfile(rd AuthInterface, tk TokenInterface) *profileHandler {
	return &profileHandler{rd, tk}
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"apiSecret"`
}

type Todo struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (h *profileHandler) Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if _, ok := h.rd.FetchAuth(u.ID); ok == nil {
		//do something here

		c.JSON(http.StatusUnauthorized, fmt.Sprintf("deneme %s", ok))
		return
	}
	ts, err := h.tk.CreateToken(u.ID, u.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := h.rd.CreateAuth(u.ID, u.Username, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *profileHandler) Logout(c *gin.Context) {
	//If metadata is passed and the tokens valid, delete them from the redis store
	metadata, _ := h.tk.ExtractTokenMetadata(c.Request)

	if metadata != nil {
		deleteErr := h.rd.DeleteTokens(metadata)
		if deleteErr != nil {
			c.JSON(http.StatusBadRequest, deleteErr.Error())
			return
		}
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func (h *profileHandler) CreateTodo(c *gin.Context) {
	var td Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	metadata, err := h.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userId, err := h.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	td.UserID = userId

	//you can proceed to save the  to a database

	c.JSON(http.StatusCreated, td)
}

func (h *profileHandler) ReturnIdentity(c *gin.Context) {
	var id User

	metadata, err := h.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userId, err := h.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	id.ID = userId
	id.Username = metadata.UserName

	//you can proceed to save the  to a database

	c.JSON(http.StatusOK, id)
}

func (h *profileHandler) CreateAccount(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	if u.Password != "cok_gizli_bunu_bilen_user_register_eder" {
		c.JSON(http.StatusUnauthorized, "You better get that secret! Youre not allowed in here!")
		return
	}
	u.ID = uuid.NewV4().String()
	ts, err := h.tk.CreateToken(u.ID, u.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := h.rd.CreateAuth(u.ID, u.Username, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
		return
	}
	tokens := map[string]string{
		"access_token": ts.AccessToken,
		// "refresh_token": ts.RefreshToken,
	}
	// m = make(map[string]User)
	// m[u.ID] = u
	c.JSON(http.StatusOK, `Registerd successfully. Welcom to TUBITAK BILGEM BAG.`)
	c.JSON(http.StatusOK, tokens)
}

// func (h *profileHandler) Refresh(c *gin.Context) {
// 	mapToken := map[string]string{}
// 	if err := c.ShouldBindJSON(&mapToken); err != nil {
// 		c.JSON(http.StatusUnprocessableEntity, err.Error())
// 		return
// 	}
// 	refreshToken := mapToken["refresh_token"]

// 	//verify the token
// 	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("REFRESH_SECRET")), nil
// 	})
// 	//if there is an error, the token must have expired
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, "Refresh token expired")
// 		return
// 	}
// 	//is token valid?
// 	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
// 		c.JSON(http.StatusUnauthorized, err)
// 		return
// 	}
// 	//Since token is valid, get the uuid:
// 	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
// 	if ok && token.Valid {
// 		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
// 		if !ok {
// 			c.JSON(http.StatusUnprocessableEntity, err)
// 			return
// 		}
// 		userId, roleOk := claims["user_id"].(string)
// 		if roleOk == false {
// 			c.JSON(http.StatusUnprocessableEntity, "unauthorized")
// 			return
// 		}
// 		//Delete the previous Refresh Token
// 		delErr := h.rd.DeleteRefresh(refreshUuid)
// 		if delErr != nil { //if any goes wrong
// 			c.JSON(http.StatusUnauthorized, "unauthorized")
// 			return
// 		}
// 		//Create new pairs of refresh and access tokens
// 		ts, createErr := h.tk.CreateToken(userId)
// 		if createErr != nil {
// 			c.JSON(http.StatusForbidden, createErr.Error())
// 			return
// 		}
// 		//save the tokens metadata to redis
// 		saveErr := h.rd.CreateAuth(userId, ts)
// 		if saveErr != nil {
// 			c.JSON(http.StatusForbidden, saveErr.Error())
// 			return
// 		}
// 		tokens := map[string]string{
// 			"access_token":  ts.AccessToken,
// 			"refresh_token": ts.RefreshToken,
// 		}
// 		c.JSON(http.StatusCreated, tokens)
// 	} else {
// 		c.JSON(http.StatusUnauthorized, "refresh expired")
// 	}
// }