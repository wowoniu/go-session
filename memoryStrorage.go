package session

import (
	"errors"
	"sync"
	"time"
)

//内存版session
type SessionMemoryStorage struct {
	ActiveSessionTable map[string]*Session
	Expire             int
	SessionData        map[string]*SessionDataTable
}

type SessionDataTable struct {
	Items map[string]string
}

//存储初始化
func (s *SessionMemoryStorage) InitStorage(expire int) {
	s.Expire = expire
	s.ActiveSessionTable = make(map[string]*Session)
	s.SessionData = make(map[string]*SessionDataTable)
	go s.GC()
}

func (s *SessionMemoryStorage) SessionStart(sessionID string) (*Session, error) {
	var sessionInstance *Session
	if sessionInstance, isExisted := s.ActiveSessionTable[sessionID]; isExisted {
		//s.SessionData[sessionID].LastVisitTime = time.Now()
		sessionInstance.LastVisit = time.Now()
		return sessionInstance, nil
	}
	//会话初次初始化
	sessionInstance = &Session{
		SessionID: sessionID,
		Mutex:     &sync.Mutex{},
		Storage:   s,
		LastVisit: time.Now(),
	}
	s.ActiveSessionTable[sessionID] = sessionInstance
	s.SessionData[sessionID] = &SessionDataTable{
		Items: make(map[string]string),
	}
	return sessionInstance, nil
}

func (s *SessionMemoryStorage) Get(sessionID string, key string) (string, error) {
	dataTalbes, isExisted := s.SessionData[sessionID]
	if !isExisted {
		return "", errors.New("Invalid Session")
	}
	value, isExisted := dataTalbes.Items[key]
	if !isExisted {
		return "", errors.New("Invalid Session Key")
	}
	//fmt.Println("读取：", key, " val:", value, " Last:", dataTalbes.LastVisitTime)
	return value, nil
}

func (s *SessionMemoryStorage) Set(sessionID string, key string, val string) error {
	dataTalbes, isExisted := s.SessionData[sessionID]
	if !isExisted {
		return errors.New("Invalid Session")
	}
	dataTalbes.Items[key] = val
	return nil
}

func (s *SessionMemoryStorage) Destroy(sessionID string) {

}

//垃圾回收
func (s *SessionMemoryStorage) GC() {
	var ticker = time.Tick(1 * time.Minute)
	for {
		select {
		case <-ticker:
			//遍历
			for sessionID, session := range s.ActiveSessionTable {
				if time.Now().Sub(session.LastVisit) > time.Duration(s.Expire)*time.Second {
					//fmt.Println("垃圾回收:", sessionID)
					//fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
					//fmt.Println(session.LastVisit)
					delete(s.ActiveSessionTable, sessionID)
					delete(s.SessionData, sessionID)
				}
			}
		}

	}
}
