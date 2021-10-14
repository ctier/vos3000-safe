// myki 21kixc@gmail.com
// 2021-05-29
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var certPath = "/etc/ssl/cert.pem"
var keyPath = "/etc/ssl/key.pem"
var (
	VOS_SAFE_USERNAME = os.Getenv("VOS_SAFE_USERNAME")
	VOS_SAFE_PASSWORD = os.Getenv("VOS_SAFE_PASSWORD")
)

type ViewFunc func(http.ResponseWriter, *http.Request)

func BasicAuth(f ViewFunc, user, passwd []byte) ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		basicAuthPrefix := "Basic "

		// 获取 request header
		auth := r.Header.Get("Authorization")
		// 如果是 http basic auth
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// 解码认证信息
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):],
			)
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 && bytes.Equal(pair[0], user) &&
					bytes.Equal(pair[1], passwd) {
					// 执行被装饰的函数
					f(w, r)
					return
				}
			}
		}

		// 认证失败，提示 401 Unauthorized
		// Restricted 可以改成其他的值
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// 401 状态码
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// 需要被保护的内容
//func HelloServer(w http.ResponseWriter, req *http.Request) {
//    io.WriteString(w, "hello, world!\n")
//}

func safe(w http.ResponseWriter, r *http.Request) {
	socket := strings.Split(r.RemoteAddr, ":")
	ip := socket[0] + "/32"
	fmt.Println("Match alright: ", ip)
	cmd := exec.Command("/sbin/iptables", "-A", "INPUT", "-s", ip, "-j", "ACCEPT")
	err := cmd.Run()
	fmt.Println(err)
	fmt.Fprintf(w, ip+" 已授权")
}

func main() {
	user := []byte(VOS_SAFE_USERNAME)
	passwd := []byte(VOS_SAFE_PASSWORD)

	// 装饰需要保护的 handler
	http.HandleFunc("/safe", BasicAuth(safe, user, passwd))

	log.Println("Listen :8000")
	err := http.ListenAndServeTLS(":8000", certPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
