package session

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

//全局session管理器
type SessionManager struct {
	Storage     SessionStorage
	Expire      int
	SessionType string
	SessionKey  string
}

//session 存储接口
type SessionStorage interface {
	//存储初始化
	InitStorage(int)
	//会话初始化
	SessionStart(sessionID string) (*Session, error)
	Set(sessionID string, key string, val string) error
	Get(sessionID string, key string) (string, error)
	Destroy(sessionID string)
	GC()
}

//会话session实例
type Session struct {
	SessionID string
	Storage   SessionStorage
	Mutex     *sync.Mutex
	LastVisit time.Time
}

var GSessionManager *SessionManager

/****session管理器*****/
func NewSessionManager(expire int, sessionType string, sessionKey string) *SessionManager {
	if GSessionManager == nil {
		//初始化存储器
		sessionStorageInstance := &SessionMemoryStorage{}
		sessionStorageInstance.InitStorage(expire)
		//实例化全局管理器
		GSessionManager = &SessionManager{
			Storage:     sessionStorageInstance,
			Expire:      expire,
			SessionType: sessionType,
			SessionKey:  sessionKey,
		}
	}

	return GSessionManager
}

func (m *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) (*Session, error) {
	if m.SessionType == "cookie" {
		sessionID := ""
		//尝试从cookie获取session
		if sessionIDCookie, err := r.Cookie(m.SessionKey); err != nil {
			//生成sessionID
			sessionID = m.createSessionID()
			//fmt.Println("设置COOKIE:", cookie)
		} else {
			sessionID = sessionIDCookie.Value
		}
		//刷新cookie
		cookie := http.Cookie{Name: m.SessionKey, Value: sessionID, Path: "/", MaxAge: m.Expire}
		http.SetCookie(w, &cookie)
		//fmt.Println("SESSIONID:", sessionID)
		session, err := m.Storage.SessionStart(sessionID)
		return session, err
	}
	//其他类型支持 TODO
	return nil, errors.New("不支持的sessionType")
}

//生成唯一ID TODO
func (m SessionManager) createSessionID() string {
	return fmt.Sprintf("%v", time.Now().Unix())
}

/************session实例*****************/
func (m *Session) Get(key string) (string, error) {
	return m.Storage.Get(m.SessionID, key)
}

func (m *Session) Set(key string, val string) {
	m.Mutex.Lock()
	m.Storage.Set(m.SessionID, key, val)
	m.Mutex.Unlock()
}

func (m *Session) Destroy() {
	m.Storage.Destroy(m.SessionID)
}
