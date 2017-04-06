package main

import "sql"

func main() {
	var sqlObj sql.SqlAlchemy
	sqlObj.Select("session").F("*").W("session_id", "aaa", "=").Or("message_id", 123, "=").ExecQuery()
	sqlObj.Select("session").F("*").ExecQuery()
	sqlObj.Insert("session").V("session_id", "aaa").V("message_id", 123).Execute()
	sqlObj.Update("session").S("session_id", "aaa").S("message_id", 123).And("session_id", "aaa", "=").Or("message_id", 123, "=").Execute()
	sqlObj.Update("session").S("session_id", "aaa").S("message_id", 123).Execute()
	sqlObj.Delete("session").W("session_id", "aaa", "=").Or("message_id", 123, "=").Execute()
	sqlObj.Delete("session").Execute()
}
