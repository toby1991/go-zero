package discov

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toby1991/go-zero/core/discov/internal"
	"github.com/toby1991/go-zero/core/stringx"
)

func TestRegisterAccount(t *testing.T) {
	endpoints := []string{
		"localhost:2379",
	}
	user := "foo" + stringx.Rand()
	RegisterAccount(endpoints, user, "bar")
	account, ok := internal.GetAccount(endpoints)
	assert.True(t, ok)
	assert.Equal(t, user, account.User)
	assert.Equal(t, "bar", account.Pass)
}
