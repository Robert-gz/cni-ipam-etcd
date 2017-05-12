package etcd

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/shwinpiocess/cni-ipam-etcd/backend/allocator"
)

var (
	dialTimeout    = 5 * time.Second
)

type Store struct {
	Client *clientv3.Client
	Key	string
}

func New(n *allocator.IPAMConfig) (*Store, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:      n.EtcdEndpoints,
		DialTimeout:    dialTimeout,
	})
	if err != nil {
		return err
	}

	return &Store{*client, n.Name}

}


func GetKV(k string, client *clientv3.Client) (*clientv3.GetResponse, error) {
	ctx, cancel := context.withTimeout(context.Background(), requestTimeout)
	resp, err := client.Get(ctx, k)
	cancel()
	if err != nil {
		return err
	}

	return resp, nil
}

func PutKV(k string, val string, client *clientv3.Client) (*clientv3.PutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Put(ctx, k, val)
	cancel()
	if err != nil {
		return err
	}

	return resp, nil
}

func (s *Store) Lock() error {
	session, err := concurrency.NewSession(s.Client)
	if err != nil {
		return err
	}

	mutex := concurrency.NewMutex(session, s.Key)
	ctx, cancel := context.WithCancel(context.TODO())

	if err := mutex.Lock(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Store) Unlock() error {
	session, err := concurrency.NewSession(s.Client)
	if err != nil {
		return err
	}

	mutex := concurrency.NewMutex(session, s.Key)
	ctx, cancel := context.WithCancel(context.TODO())

	if err := mutext.Unlock(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Store) Close() error {
}

func (s *Store) Reserve(id string, ip net.IP) (bool, error) {
	pair, _ := GetKV(s.Key, s.Client)
	if len(pair) != 0 {
		return false, nil
	}

	PutKV(s.Key, ip.String(), s.Client)

	return true, nil
}

func Release(ip net.IP) error {
}

func ReleaseByID(id string) error {
}

