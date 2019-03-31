# go-session
golang session管理实现


#使用


###在创建HTTPServer时 创建全局session管理器

sessionManager:=session.NewSessionManager(1800, "cookie", "sessionid")


###在http的响应函数中调用SessionStart函数 初始化session会话 和使用
func apiLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionManager.SessionStart(w, r)
	session.Set("username","zhangsan")
    fmt.Println(session.Get("username"))

}
