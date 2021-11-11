package cmd

import (
	"os"
	"time"

	"github.com/zemirco/couchdb"
)

type AuthInterface interface {
	CreateAuth(string, string, *TokenDetails) error
	FetchAuth(string) (string, error)
	// DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}

type service struct {
	client *couchdb.Client
}

var _ AuthInterface = &service{}

var databaseName = os.Getenv("DATABASE_NAME")

func NewAuth(client *couchdb.Client) *service {
	return &service{client: client}
}

type AccessDetails struct {
	TokenUuid string
	UserId    string
	UserName  string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
type saveDetails struct {
	couchdb.Document
	Uuid      string        `json:"_id"`
	TokenType string        `json:"type"`
	UserId    string        `json:"userId"`
	UserName  string        `json:"Username`
	Expires   time.Duration `json:"AtExpires"`
}

//Save token metadata to Redis
func (tk *service) CreateAuth(userId string, userName string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	// rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	db := tk.client.Use("bulutzincir")
	atCreated := &saveDetails{
		Uuid:      td.TokenUuid,
		TokenType: "access",
		UserId:    userId,
		UserName:  userName,
		Expires:   at.Sub(now),
	}
	_, err := db.Post(atCreated)
	if err != nil {
		return err
	}
	// rtCreated := &saveDetails{
	// 	Uuid:      td.RefreshUuid,
	// 	TokenType: "refresh",
	// 	UserId:    userId,
	// 	Expires:   rt.Sub(now),
	// }
	// _, err = db.Post(rtCreated)
	// if err != nil {
	// 	return err
	// }
	return nil
}

//Check the metadata saved
func (tk *service) FetchAuth(tokenUuid string) (string, error) {
	// userid, err := tk.client.Get(tokenUuid).Result()
	db := tk.client.Use("bulutzincir")
	pulled := &saveDetails{}

	if err := db.Get(pulled, tokenUuid); err != nil {
		return "", err
	}
	return pulled.UserId, nil
}

//Once a user row in the token table
func (tk *service) DeleteTokens(authD *AccessDetails) error {
	//get the refresh uuid
	// refreshUuid := fmt.Sprintf("%s++%s", authD.TokenUuid, authD.UserId)
	//delete access token
	db := tk.client.Use("bulutzincir")

	access := &saveDetails{Uuid: authD.TokenUuid}
	if _, err := db.Delete(access); err != nil {
		return err
	}

	// refresh := &saveDetails{Uuid: refreshUuid}
	// if _, err := db.Delete(refresh); err != nil {
	// 	return err
	// }
	//delete refresh token

	//When the record is deleted, the return value is 1
	//if deletedAt != 1 || deletedRt != 1 {
	//return errors.New("something went wrong")
	//}
	return nil
}

// func (tk *service) DeleteRefresh(refreshUuid string) error {
// 	//delete refresh token
// 	//deleted, err := tk.client.Del(refreshUuid).Result()
// 	//if err != nil || deleted == 0 {
// 	//return err
// 	//}

// 	db := tk.client.Use("bulutzincir")

// 	refresh := &saveDetails{Uuid: refreshUuid}
// 	if _, err := db.Delete(refresh); err != nil {
// 		return err
// 	}
// 	return nil
// }
