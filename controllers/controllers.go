package controllers

import (
	"Sto_kyc/config"
	"Sto_kyc/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"crypto/ecdsa"
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type User struct {
	Id int `json:"id"`
	*models.User
}

// Handler Error Json
type ActionErr struct {
	ActionError string `json:"fail"`
}

// Handler 200 OK JSON
type ActionSuccess struct {
	Action string `json:"success"`
}

func GetKycItems(c *gin.Context) {
	// data := config.V.Selectors
	c.JSON(http.StatusOK, config.V.Selectors)
	return
}

func Apply(c *gin.Context) {
	user := new(models.User)
	user.Name = c.PostForm("name")
	user.Address = c.PostForm("address")
	user.Email = c.PostForm("email")
	user.Selector = c.PostForm("selector")
	user.Status = 2
	user.Passport = ""

	ok, err := models.CheckUserExists(user)
	log.Println(user)
	log.Println(err)
	if ok == true {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}

	file, handler, err := c.Request.FormFile("passport")
	fmt.Println(handler.Filename)
	filename := handler.Filename
	user.Passport = user.Address + user.Selector + filename
	out, err := os.Create(config.V.ImagesDir + user.Passport)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}

	_, err = models.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}

	c.JSON(http.StatusOK, ActionSuccess{"申请已提交！"})
	return
}

func Query(c *gin.Context) {
	address := c.PostForm("address")
	selector := c.PostForm("selector")
	ok, err := models.CheckUserCertified(address, selector)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}
	if ok == false {
		c.JSON(http.StatusForbidden, ActionErr{"此用户未通过KYC认证！"})
		return
	} else {
		c.JSON(http.StatusOK, ActionSuccess{"此用户已通过KYC认证！"})
		return
	}
	c.JSON(http.StatusOK, ActionSuccess{"此用户已通过KYC认证！"})
	return
}

func GetCheckData(c *gin.Context) {
	user := new(User)
	id := c.PostForm("userId")
	msg := c.PostForm("msg")
	sig := c.PostForm("sig")
	selector := c.PostForm("selector")

	ok, err := verifyPermission(msg, sig)
	log.Println(msg)
	log.Println(sig)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}
	if ok != true {
		c.JSON(http.StatusForbidden, ActionErr{"非管理员不能操作！"})
		return
	}
	userid, _ := strconv.Atoi(id)
	user.Id, user.User, err = models.ReadUser(userid, selector)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

func Certify(c *gin.Context) {
	id := c.PostForm("userId")
	msg := c.PostForm("msg")
	sig := c.PostForm("sig")
	ok, err := verifyPermission(msg, sig)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
	}
	if ok != true {
		c.JSON(http.StatusForbidden, ActionErr{"非管理员不能操作！"})
		return
	}
	userid, _ := strconv.Atoi(id)
	err = models.UpdateUser(1, userid)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
	}
	c.JSON(http.StatusOK, ActionSuccess{"审批同意！"})
	return
}

func Reject(c *gin.Context) {
	id := c.PostForm("userId")
	msg := c.PostForm("msg")
	sig := c.PostForm("sig")
	ok, err := verifyPermission(msg, sig)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}
	if ok != true {
		c.JSON(http.StatusForbidden, ActionErr{"非管理员不能操作！"})
		return
	}
	userid, _ := strconv.Atoi(id)
	err = models.UpdateUser(0, userid)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}
	c.JSON(http.StatusOK, ActionSuccess{"审批拒绝！"})
	return
}

func verifyPermission(message, signature string) (bool, error) {
	address, err := EcRecover(message, signature)
	if err != nil {
		return false, err
	}
	log.Println(address)

	exp, _ := strconv.Atoi(message)
	exp = exp/1000 + 7200
	if time.Now().Unix() > int64(exp) {
		return false, errors.New("权限过期，请重新加载钱包！")
	}

	for _, v := range config.V.Managers {
		if v == address {
			return true, nil
		}
	}

	return false, errors.New("权限错误！")
}

func EcRecover(msg, sigHex string) (addr string, err error) {
	defer func() {
		if r := recover(); r != nil {
			addr = common.Address{}.String()
			err = r.(error)
		}
	}()

	data := []byte(msg)
	hash := SignHash(data)

	sig := hexutil.MustDecode(sigHex)
	// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	if sig[64] != 27 && sig[64] != 28 {
		return common.Address{}.String(), errors.New("nvalid Ethereum signature (V is not 27 or 28)")
	}
	sig[64] -= 27

	sigPiblicKey, err := SigToPub(hash, sig)

	if err != nil {
		panic(err)
	}

	commAddr := PubkeyToAddress(*sigPiblicKey)

	return commAddr.String(), nil
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	// Convert to btcec input format with 'recovery id' v at the beginning.
	btcsig := make([]byte, 65)
	btcsig[0] = sig[64] + 27
	copy(btcsig[1:], sig)

	pub, _, err := btcec.RecoverCompact(btcec.S256(), btcsig, hash)
	return (*ecdsa.PublicKey)(pub), err
}

// PubkeyToAddress converts ecdsa.PublicKey to ethereum address.
func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	pubBytes := crypto.FromECDSAPub(&p)
	return common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:])
}

func SignHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}
