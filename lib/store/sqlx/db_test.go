package sqlx

import (
	"fmt"
	"testing"
	"time"
)

func TestCommonDB_QueryRow(t *testing.T) {
	db := NewMySQL("root:asdfasdf@tcp(192.168.0.166:3306)/nest_label?parseTime=true")
	result := struct {
		Name  string
		Total int
		Price float32
		//Total int    `db:"totalx"`
		//Name  string `db:"book"`
	}{}
	err := db.Query(&result, "select book, count(0) total from book group by book order by total desc")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if result.Total != 1433 {
		t.Fatalf("期望结果 %v - 实际结果 %v", 1433, result.Total)
	} else {
		fmt.Println(result)
	}
}

type Time struct {
	time.Time
}

// String returns current time object as string.
func (t *Time) String() string {
	if t == nil {
		return ""
	}
	if t.IsZero() {
		return ""
	}
	return t.Format("Y-m-d H:i:s")
}

func TestCommonDB_QueryRows(t *testing.T) {
	dataSourceName := "root:asdfasdf@tcp(192.168.0.166:3306)/nest_label?parseTime=true"
	db := NewMySQL(dataSourceName)
	//db := sqlx.NewMysql(dataSourceName)
	//var book struct {
	//	Name  string
	//	Total int
	//	Price float32
	//	//Total int    `db:"totalx"`
	//	//Name  string `db:"book"`
	//}
	//var books []struct {
	//	Total         int    `db:"totalx"`
	//	Name          string `db:"book"`
	//	NotExistField int    `db:"y"`
	//}
	//var accountKinds []struct {
	//	Id   int
	//	Name string
	//}

	//var adminUsers []struct {
	//	UserId    int       `db:"user_id"`
	//	AdminId   int       `db:"admin_id"`
	//	Txt       string    `db:"txt"`
	//	CreatedAt time.Time `db:"created_at"`
	//}
	var createdTime time.Time
	//errBook := db.Query(&book, "select book, count(0) total from book group by book order by total desc")
	//errBooks := db.Query(&books, "select book, count(0) totalx, 1 as x, 2 as y from book group by book order by totalx desc")
	//errAccountKinds := db.Query(&accountKinds, "select id, value as name from nest_user.account_kind")
	//errAdminUsers := db.Query(&adminUsers, "select user_id, admin_id, txt from nest_admin.admin_user")
	errAdminUsers := db.Query(&createdTime, "select created_at from nest_admin.admin_user")
	//errAdminUsers := db.QueryRows(&createdTime, "select created_at from nest_admin.admin_user")

	//if errBook != nil {
	//	t.Fatalf("%v", errBook)
	//}
	//if errBooks != nil {
	//	t.Fatalf("%v", errBooks)
	//}
	//if errAccountKinds != nil {
	//	t.Fatalf("%v", errAccountKinds)
	//}
	if errAdminUsers != nil {
		t.Fatalf("%v", errAdminUsers)
	}

	//fmt.Println(book)
	//
	//for _, book := range books {
	//	fmt.Println(book)
	//}
	//
	//for _, accountKind := range accountKinds {
	//	fmt.Println(accountKind)
	//}

	//for _, adminUser := range adminUsers {
	//	fmt.Println(adminUser)
	//}

	fmt.Println(createdTime)
}