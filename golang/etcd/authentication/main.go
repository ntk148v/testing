package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/namespace"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "root",
		Password:  "test",
	})

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// Create namespace
	clusterns := "/clusterns/"
	cli.KV = namespace.NewKV(cli.KV, clusterns)
	cli.Watcher = namespace.NewWatcher(cli.Watcher, clusterns)
	cli.Lease = namespace.NewLease(cli.Lease, clusterns)

	if _, err = cli.AuthEnable(context.TODO()); err != nil {
		panic(err)
	}

	// Create users with specific role
	// Setup
	var (
		user   string
		pass   string
		role   string
		userns string
	)

	for i := 1; i < 4; i++ {
		user = fmt.Sprintf("user%d", i)
		pass = fmt.Sprintf("pass%d", i)
		role = fmt.Sprintf("role%d", i)
		userns = fmt.Sprintf("%s%d/", clusterns, i)
		if _, err = cli.RoleAdd(context.TODO(), role); err != nil && err != rpctypes.ErrRoleAlreadyExist {
			panic(err)
		}
		if _, err = cli.RoleGrantPermission(
			context.TODO(),
			role,                               // role name
			userns,                             // key
			clientv3.GetPrefixRangeEnd(userns), // range end
			clientv3.PermissionType(clientv3.PermReadWrite),
		); err != nil {
			panic(err)
		}
		if _, err = cli.UserAdd(context.TODO(), user, pass); err != nil && err != rpctypes.ErrUserAlreadyExist {
			panic(err)
		}
		if _, err := cli.UserGrantRole(context.TODO(), user, role); err != nil {
			panic(err)
		}
		log.Printf("Create user %s with RW role %s on namespace %s\n", user, role, userns)
	}

	// Do something cool
	var userCli *clientv3.Client
	for i := 1; i < 4; i++ {
		user = fmt.Sprintf("user%d", i)
		pass = fmt.Sprintf("pass%d", i)
		userns = fmt.Sprintf("%d/", i)
		userCli, err = clientv3.New(clientv3.Config{
			Endpoints: []string{"127.0.0.1:2379"},
			Username:  user,
			Password:  pass,
		})
		defer userCli.Close()
		userCli.KV = namespace.NewKV(userCli.KV, clusterns)
		userCli.Watcher = namespace.NewWatcher(userCli.Watcher, clusterns)
		userCli.Lease = namespace.NewLease(userCli.Lease, clusterns)
		// Put to the correct key
		log.Println(fmt.Sprintf("%s%s", userns, "foo"))
		if putResp, err := userCli.Put(context.TODO(), fmt.Sprintf("%s%s", userns, "foo"), "bar"); err != nil {
			panic(err)
		} else {
			log.Printf("Put key %s value with user %s: %s", fmt.Sprintf("%s%s", userns, "foo"),
				user, putResp)
		}

		// Put to the wrong key, expect the permission denied error
		userns = fmt.Sprintf("%s%d/", clusterns, i+10)
		if _, err := userCli.Put(context.TODO(), fmt.Sprintf("%s%s", userns, "foo"), "bar"); err != nil {
			if err != rpctypes.ErrPermissionDenied {
				log.Fatal(err)
			} else {
				log.Println("Permission denied")
			}

		}
	}

	// Cleanup
	defer func() {
		for i := 1; i < 4; i++ {
			user = fmt.Sprintf("user%d", i)
			role = fmt.Sprintf("role%d", i)
			if _, err = cli.UserDelete(context.TODO(), user); err != nil {
				log.Fatal(err)
			}
			if _, err = cli.RoleDelete(context.TODO(), role); err != nil {
				log.Fatal(err)
			}
			if _, err = cli.Delete(context.TODO(), strconv.Itoa(i), clientv3.WithPrefix()); err != nil {
				log.Fatal(err)
			}
		}
		log.Println("Cleaned up")
	}()
}
