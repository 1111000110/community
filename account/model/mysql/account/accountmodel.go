package account

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AccountModel = (*customAccountModel)(nil)

type (
	// AccountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountModel.
	AccountModel interface {
		accountModel
		withSession(session sqlx.Session) AccountModel
	}

	customAccountModel struct {
		*defaultAccountModel
	}
)

// NewAccountModel returns a model for the database table.
func NewAccountModel(conn sqlx.SqlConn) AccountModel {
	return &customAccountModel{
		defaultAccountModel: newAccountModel(conn),
	}
}

func (m *customAccountModel) withSession(session sqlx.Session) AccountModel {
	return NewAccountModel(sqlx.NewSqlConnFromSession(session))
}
