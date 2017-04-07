package oph

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestCallString(t *testing.T) {
	var outString sql.NullString
	str, err := CallString("TEST", 1, nil, &outString)
	if err != nil {
		t.Errorf("CallString: err=%s", err.Error())
	}
	if str != "CALL TEST(1,NULL,@1);SELECT @1" {
		t.Errorf("CallString: expected=CALL TEST(1,NULL,@1);SELECT @1, got=%s", str)
	}
}

func ExampleCallString() {
	var outString sql.NullString
	str, err := CallString("TEST", "test", 1, nil, &outString)
	if err != nil {
		fmt.Println("ERROR", err.Error())
	} else {
		fmt.Println(str)
	}
	// Output: CALL TEST('test',1,NULL,@1);SELECT @1
}

func ExampleCallStringInjection() {
	str, err := CallString("TEST", "test';DROP TABLE USERS")
	if err != nil {
		fmt.Println("ERROR", err.Error())
	} else {
		fmt.Println(str)
	}
	// Output: CALL TEST('test'';DROP TABLE USERS')
}
