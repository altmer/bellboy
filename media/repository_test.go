package media

import (
	"github.com/altmer/bellboy/context"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"os"
	"testing"
)

var DB *sqlx.DB
var repo Repository

func setup() func() {
	testDBPath := "./test.db"
	testMediaPath := "./"

	viper.SetDefault("media_folder", testMediaPath)
	DB = context.NewDBConnection(testDBPath)

	repo = NewRepository(DB)

	return func() {
		DB.Close()
		os.Remove(testDBPath)
	}
}

func checkErrors(t *testing.T, expectedErr, err error) {
	if err == nil || expectedErr == nil {
		if err != expectedErr {
			t.Errorf("Expected error to be [%#v], got [%#v]", expectedErr, err)
		}
	} else {
		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected to fail with error [%#v], got [%#v]", expectedErr.Error(), err.Error())
		}
	}
}
