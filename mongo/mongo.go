package mongo

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

type mongoSession struct {
	dialInfo *mgo.DialInfo
	session  *mgo.Session
}

type sessionPool map[string]mongoSession

type sessionMgr struct {
	pool      sessionPool
	locker    sync.RWMutex
	allClosed bool //prevent double close operation
}

var mgr = sessionMgr{}

func (m *sessionMgr) addSession(name string, s mongoSession) {
	m.locker.Lock()
	defer m.locker.Unlock()
	if m.pool == nil {
		m.pool = make(map[string]mongoSession)
	}
	m.pool[name] = s
}

func (m *sessionMgr) getSession(name string) (mongoSession, error) {
	m.locker.Lock()
	defer m.locker.Unlock()

	empty := mongoSession{}

	if m.pool == nil {
		err := errors.New("session pull is null")
		return empty, err
	}
	if v, ok := m.pool[name]; !ok {
		err := errors.New(fmt.Sprintf("session %s not found", name))
		return empty, err
	} else {
		return v, nil
	}
}

func (m *sessionMgr) fetchMgoSession(name string, copy bool) (*mgo.Session, error) {
	m.locker.Lock()
	defer m.locker.Unlock()
	if m.pool == nil {
		err := errors.New("session pull is null")
		return nil, err
	}
	if v, ok := m.pool[name]; !ok {
		err := errors.New(fmt.Sprintf("session %s not found", name))
		return nil, err
	} else if copy {
		return v.session.Copy(), nil
	} else {
		return v.session, nil
	}
}

func (m *sessionMgr) removeOneSession(name string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	if m.pool == nil {
		log.Println("session pool has not bee initialized when trying to delete session!")
	}
	delete(m.pool, name)
}

func (m *sessionMgr) cleanAllSessions() {
	m.locker.Lock()
	defer m.locker.Unlock()
	for k, _ := range m.pool {
		delete(m.pool, k)
	}
}

func (m *sessionMgr) closeAllSessions() {
	m.locker.Lock()
	defer m.locker.Unlock()
	if !m.allClosed {
		for _, v := range m.pool {
			v.session.Close()
		}
		m.allClosed = true
	}
}

func StartUp(sessionName, hosts, dbName, userName, password string) error {
	if _, err := mgr.getSession(sessionName); err == nil {
		return nil
	}
	multiHosts := strings.Split(hosts, ",")
	return createSession(sessionName, multiHosts, dbName, userName, password)
}

//close all existing db sessions
func CloseAll() {
	mgr.closeAllSessions()
}

//close all existing db sessions and drop all session instances from session manager
func Shutdown() {
	mgr.closeAllSessions()
	mgr.cleanAllSessions()
}

func CopySession(sessionName string) (*mgo.Session, error) {
	if session, err := mgr.fetchMgoSession(sessionName, true); err != nil {
		return nil, err
	} else {
		return session, nil
	}
}

func GetSessionDBName(sessionName string) (string, error) {
	if session, err := mgr.getSession(sessionName); err != nil {
		return "", err
	} else {
		return session.dialInfo.Database, nil
	}
}

func WithConnection(d, c, s string, op DBOperator) error {
	session, err := CopySession(s)
	if err != nil {
		return err
	}
	defer session.Close()
	if op != nil {
		c := session.DB(d).C(c)
		return op(c)
	}
	return nil
}

func createSession(sessionName string, hosts []string, dbName, username, password string) error {
	myDialInfo := &mgo.DialInfo{
		Addrs:    hosts,
		Timeout:  60 * time.Second,
		Database: dbName,
		Username: username,
		Password: password,
	}
	dbSession, err := mgo.DialWithInfo(myDialInfo)
	if err != nil {
		fmt.Println("meet error while DialWithInfo")
		return err
	}
	dbSession.SetSafe(&mgo.Safe{})

	session := mongoSession{
		dialInfo: myDialInfo,
		session:  dbSession,
	}
	//mongoSession.mongoSession.SetMode(mgo.Monotonic, true)
	mgr.addSession(sessionName, session)
	return nil
}
