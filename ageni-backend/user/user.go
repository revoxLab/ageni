package user

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/readonme/open-studio/common"
	"github.com/readonme/open-studio/common/log"
	signature "github.com/readonme/open-studio/common/wallet_signature"
	"github.com/readonme/open-studio/conf"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/dal/model"
	"github.com/readonme/open-studio/dal/query"
	"github.com/readonme/open-studio/lib"
	"gorm.io/gorm"
)

func GetUserByIds(ids []int64) []*model.User {
	user := query.Use(dal.StudioDB).User
	users, _ := user.WithContext(context.TODO()).Where(user.ID.In(ids...)).Find()
	return users
}

func GetUser(id int64) (*model.User, error) {
	user := query.Use(dal.StudioDB).User
	u, err := user.WithContext(context.TODO()).Where(user.ID.Eq(id)).First()
	return u, err
}

func WebWalletLogin(ctx context.Context, in *WalletLoginRequest) (userModel *model.User, token string, err error) {

	_, err = ExternalWalletLogin(ctx, in)
	if err != nil {
		return
	}

	user := query.Use(dal.StudioDB).User
	userModel, err = user.WithContext(context.TODO()).Where(user.WalletAddress.Eq(in.WalletAddress)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userModel = &model.User{
			WalletAddress: in.WalletAddress,
		}
		err := user.WithContext(context.TODO()).Create(userModel)
		if err != nil {
			return nil, "", err
		}
	}

	userModel.WalletAddress = in.WalletAddress
	token, err = lib.GenJwtToken(userModel, conf.Conf.JWTToken)
	return
}

const LenseSource = "Sign in to REVOX"

type WalletLoginRequest struct {
	From          string `protobuf:"bytes,8,opt,name=from,proto3" json:"from,omitempty"`
	Signature     string `protobuf:"bytes,9,opt,name=signature,proto3" json:"signature,omitempty"`
	WalletAddress string `protobuf:"bytes,10,opt,name=wallet_address,json=walletAddress,proto3" json:"wallet_address,omitempty"`
	WalletType    int32  `protobuf:"varint,11,opt,name=wallet_type,json=walletType,proto3" json:"wallet_type,omitempty"`
}

func ExternalWalletLogin(ctx context.Context, in *WalletLoginRequest) (string, error) {
	if len(in.Signature) == 0 || len(in.WalletAddress) == 0 {
		return "", common2.ErrParamCheck
	}
	source := LenseSource
	err := ValidateWalletSignature(in.Signature, in.WalletAddress, source)
	if err != nil {
		return "", err
	}

	return in.WalletAddress, nil
}

func ValidateWalletSignature(sign, address string, source string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = common2.ErrWalletSignature
			log.Errorf("ValidateWalletSignature sign:%s address:%s source:%s", sign, address, source)
		}
	}()

	if !has0xPrefix(address) {
		return common2.ErrWalletSignature
	}
	if !common.IsHexAddress(address) {
		return common2.ErrWalletSignature
	}
	pass, err := signature.VerifyMessage(source, sign, address)
	if !pass {
		if err != nil {
			log.Errorf("validate wallet sig failed err=%v", err)
		}
		return common2.ErrWalletSignature
	}
	return err
}

func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}
