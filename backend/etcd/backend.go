// Copyright 2015 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

const lastIPFile = "last_reserved_ip"

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 2 * time.Second
	endpoints      = []string{"127.0.0.1:2379"}
	defaultSessionTTL = 60
)

type Store struct {
	EtcdMutex	*concurrency.Mutex
	Key 		string
}

type Network struct {
        RangeStart net.IP        `json:"rangeStart"`
        RangeEnd   net.IP        `json:"rangeEnd"`
        Subnet     types.IPNet   `json:"subnet"`
        Gateway    net.IP        `json:"gateway"`
        Routes     []types.Route `json:"routes"`
        ResolvConf string        `json:"resolvConf"`
}


func ConnectStore(endpoints []string) (etcd *clientv3, err error) {
	config := clientv3.Config{
		Endpoints:	endpoints,
		DialTimeout: 	dialTimeout,
	}

	etcd, err = clientv3.New(config)
	if err != nil {
		panic(err)
	}
	return etcd, err
}

func NewEtcdLock(n *config.IPAMConfig) (*concurrency.Mutex, error) {
	config := clientv3.Config{
                Endpoints:      n.EtcdEndpoints,
                DialTimeout:    dialTimeout,
        }

	client, err = clientv3.New(config)
        if err != nil {
                panic(err)
        }

	opts := &concurrency.sessionOptions{ttl: defaultSessionTTL, ctx: client.Ctx()}

	s, err := concurrency.NewSession(client, opts)
	if err != nil {
		panic(err)
	}

	m := concurrency.NewMutex(s, n.EtcdPrefix)

	return m, nil
}

func New(n *config.IPAMConfig) (*Store, error) {
	// etcd, err := ConnectStore(n.Endpoints)
	// if err != nil {
	// 	panic(err)
	// }

	// network, err := NetConfigJson(n)
	// key, err := InitStore(n.Name, network, etcd)

	lk, err := NewEtcdLock(n)
	if err != nil {
		return nil, err
	}

	return &Store{*lk, key}, nil
}

func InitStore(k string, network []byte, etcd *clientv3) (store string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := etcd.Get(ctx, k)
	cancel()
	if err != nil {
		panic(err)
	}
	//for _, ev := range resp.Kvs {
	//    fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	//}

	return k, nil
}

func NetConfigJson(n *config.IPAMConfig) (config []byte, err error) {
	ip_set := IP_Settings{
		Gw:	n.Gateway,
		Net:	n.Subnet,
		Start:	n.RangeStart,
		End:	n.RangeEnd,
		Routes:	n.Routes,
	}
	conf, err = json.Marshal(ip_set)
	return conf, err
}

// func (s *Store) Lock() error {
// 
// 	Session := s.Consul.Session()
// 	kv := s.Consul.KV()
// 	var entry *api.SessionEntry
// 
// 	// create session
// 	id, _, err := Session.Create(entry, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// get pair object from consul
// 	pair, _, err := kv.Get(s.Key, nil)
// 	pair.Session = id
// 	if err != nil {
// 		panic(err)
// 	}
// 	// acquire is false
// 	acq := false
// 	attempts := 0
// 	// will try 10 times to get the lock - 10 seconds
// 	for acq != true {
// 		if attempts == 10 {
// 			panic("Wasn't able to acquire the lock in 10 seconds")
// 		}
// 		acq, _, err = kv.Acquire(pair, nil)
// 		if err != nil {
// 			panic(err)
// 		}
// 		attempts += 1
// 		time.Sleep(1000 * time.Millisecond)
// 	}
// 	return err
// }

func (s *Store) Lock() error {
	client := s.s.Client()
	if err := s.Mutex.Lock(client.Ctx()); err != nil {
		return err
	}

	return nil
}

// func (s *Store) Unlock() error {
// 	ses := s.etcd.Session()
// 	kv := s.etcd.KV()
// 	pair, _, err := kv.Get(s.Key, nil)
// 	kv.Release(pair, nil)
// 
// 	sessions, _, _ := ses.List(nil)
// 	for _, session := range sessions {
// 		ses.Destory(session.ID, nil)
// 	}
// 	if err != nil {
// 		panic(err)
// 	}
// 	return nil
// }

func (s *Store) Unlock() error {
	client := s.s.Client()
	if err := s.Mutex.Unlock(client.Ctx()); err != nil {
		return err
	}
}

func (s *Store) Reserve(id string, ip net.IP) (bool, error) {
	// fname := filepath.Join(s.dataDir, ip.String())
	// f, err := os.OpenFile(fname, os.O_RDWR|os.O_EXCL|os.O_CREATE, 0644)
	// if os.IsExist(err) {
	// 	return false, nil
	// }
	// if err != nil {
	// 	return false, err
	// }
	// if _, err := f.WriteString(strings.TrimSpace(id)); err != nil {
	// 	f.Close()
	// 	os.Remove(f.Name())
	// 	return false, err
	// }
	// if err := f.Close(); err != nil {
	// 	os.Remove(f.Name())
	// 	return false, err
	// }
	// // store the reserved ip in lastIPFile
	// ipfile := filepath.Join(s.dataDir, lastIPFile)
	// err = ioutil.WriteFile(ipfile, []byte(ip.String()), 0644)
	// if err != nil {
	// 	return false, err
	// }
	// return true, nil

	path := s.Key() + "/" + ip.String()
	pair, _ := G
}

// LastReservedIP returns the last reserved IP if exists
func (s *Store) LastReservedIP() (net.IP, error) {
	ipfile := filepath.Join(s.dataDir, lastIPFile)
	data, err := ioutil.ReadFile(ipfile)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(string(data)), nil
}

func (s *Store) Release(ip net.IP) error {
	return os.Remove(filepath.Join(s.dataDir, ip.String()))
}

// N.B. This function eats errors to be tolerant and
// release as much as possible
func (s *Store) ReleaseByID(id string) error {
	err := filepath.Walk(s.dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		if strings.TrimSpace(string(data)) == strings.TrimSpace(id) {
			if err := os.Remove(path); err != nil {
				return nil
			}
		}
		return nil
	})
	return err
}
