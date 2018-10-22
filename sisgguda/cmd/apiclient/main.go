package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	h "net/http"
	"os"
	"strconv"

	"github.com/lamg/regapi"
)

var apihost, user, pass, authHd string
var rd *bufio.Reader

func Login() {
	var e error

	if user == "" {
		print("user:")
		user = readString()

		print("pass:")
		pass = readString()
	}

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	c := &regapi.Credentials{
		User: user,
		Pass: pass,
	}
	var bs []byte
	bs, e = json.Marshal(c)
	var r *h.Response

	if e == nil {
		var rq *h.Request
		if e == nil {
			bf := bytes.NewReader(bs)
			rq, e = h.NewRequest(h.MethodPost, apihost+"/auth", bf)
		}
		if e == nil {
			r, e = h.DefaultClient.Do(rq)
		}

		if e != nil {
			println("e: ", e)
			return
		}
	} else {
		println("e: ", e)
		return
	}

	if r.StatusCode != 200 {
		printBody(r)
		return
	}

	println(r.Status)

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		println("e: ", e)
		return
	}

	authHd = string(body)
	println("OK")
}

func AddUser() {
	print("username: ")
	username := readString()

	print("description: ")
	description := readString()

	print("rol:")
	rol := readString()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	s := fmt.Sprintf(`{"username":"%s","description":"%s","rol":"%s"}`, username, description, rol)
	println(s)

	bs := []byte(s)

	q, e := h.NewRequest(h.MethodPost, apihost+"/user", bytes.NewReader(bs))
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	q.Header.Set("Content-Type", "aplication/json")
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func UpdUser() {
	print("userId: ")
	userId := readInt()

	print("description: ")
	description := readString()

	print("rol:")
	rol := readString()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	s := fmt.Sprintf(`{"description":"%s","rol":"%s"}`, description, rol)
	println(s)

	bs := []byte(s)

	url := fmt.Sprintf("%s/user/%d", apihost, userId)
	println("url =", url)
	q, e := h.NewRequest(h.MethodPut, url, bytes.NewReader(bs))
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	q.Header.Set("Content-Type", "aplication/json")
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func DelUser() {
	print("userId: ")
	userId := readInt()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	url := fmt.Sprintf("%s/user/%d", apihost, userId)
	println("url =", url)
	q, e := h.NewRequest(h.MethodDelete, url, nil)
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func AddGroup() {
	print("group name: ")
	name := readString()

	print("description: ")
	description := readString()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	s := fmt.Sprintf(`{"name":"%s","description":"%s"}`, name, description)
	println(s)

	bs := []byte(s)

	q, e := h.NewRequest(h.MethodPost, apihost+"/group", bytes.NewReader(bs))
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	q.Header.Set("Content-Type", "aplication/json")
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func UpdGroup() {
	print("groupId: ")
	groupId := readInt()

	print("name: ")
	name := readString()

	print("description: ")
	description := readString()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	s := fmt.Sprintf(`{"name":"%s","description":"%s"}`, name, description)
	println(s)

	bs := []byte(s)

	url := fmt.Sprintf("%s/group/%d", apihost, groupId)
	println("url =", url)
	q, e := h.NewRequest(h.MethodPut, url, bytes.NewReader(bs))
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	q.Header.Set("Content-Type", "aplication/json")
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func DelGroup() {
	print("groupId: ")
	groupId := readInt()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	url := fmt.Sprintf("%s/group/%d", apihost, groupId)
	println("url =", url)
	q, e := h.NewRequest(h.MethodDelete, url, nil)
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func AddSpecialist() {
	print("Group Id: ")
	groupId := readInt()

	print("User Id: ")
	userId := readInt()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	s := fmt.Sprintf(`{"groupId":%d,"userId":%d}`, groupId, userId)
	println(s)

	bs := []byte(s)

	q, e := h.NewRequest(h.MethodPost, apihost+"/gspecialist", bytes.NewReader(bs))
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	q.Header.Set("Content-Type", "aplication/json")
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func DelSpecialist() {
	print("Group Id: ")
	groupId := readInt()

	print("User Id: ")
	userId := readInt()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	s := fmt.Sprintf(`{"groupId":%d,"userId":%d}`, groupId, userId)
	println(s)

	bs := []byte(s)

	q, e := h.NewRequest(h.MethodDelete, apihost+"/gspecialist", bytes.NewReader(bs))
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	q.Header.Set("Content-Type", "aplication/json")
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func AddGroupAdUser() {
	print("Group Id: ")
	groupId := readInt()

	print("AdUsername: ")
	username := readString()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	s := fmt.Sprintf(`{"groupId":%d,"adUser":"%s"}`, groupId, username)
	println(s)

	bs := []byte(s)

	q, e := h.NewRequest(h.MethodPost, apihost+"/group", bytes.NewReader(bs))
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	q.Header.Set("Content-Type", "aplication/json")
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func DelGroupAdUser() {
	print("Group Id: ")
	groupId := readInt()

	print("User Id: ")
	userId := readInt()

	tr := &h.Transport{
		Proxy: nil,
	}
	h.DefaultClient.Transport = tr

	s := fmt.Sprintf(`{"groupId":%d,"userId":%d}`, groupId, userId)
	println(s)

	bs := []byte(s)

	q, e := h.NewRequest(h.MethodDelete, apihost+"/gspecialist", bytes.NewReader(bs))
	if e != nil {
		panic(e)
	}

	q.Header.Set(regapi.AuthHd, authHd)
	q.Header.Set("Content-Type", "aplication/json")
	r, e := h.DefaultClient.Do(q)

	println(r.Status)

	if e == nil {
		printBody(r)
	} else {
		fmt.Printf(e.Error(), '\n')
	}
}

func readString() string {
	s, e := rd.ReadString('\n')
	if e != nil {
		panic(e)
	}
	return s[:len(s)-1]
}

func readInt() int {
	s := readString()
	num, e := strconv.Atoi(s)
	if e != nil {
		panic(e)
	}
	return num
}

func main() {
	flag.StringVar(&apihost, "apihost", "http://localhost:8080", "Api Host")
	flag.StringVar(&user, "u", "", "username")
	flag.StringVar(&pass, "p", "", "password")
	flag.Parse()

	rd = bufio.NewReader(os.Stdin)

	for {
		println("Menu.")
		println("0. Login")

		println("1. Add User")
		println("2. Update User")
		println("3. Delete User")

		println("4. Add Group")
		println("5. Update Group")
		println("6. Delete Group")

		println("7. Add Specialist")
		println("8. Delete Specialist")

		println("9. Add GroupAdUser")
		println("10. Delete GroupAdUser")

		print("Enter op: ")
		op := readInt()

		if op == 0 {
			Login()
		} else if op == 1 {
			AddUser()
		} else if op == 2 {
			UpdUser()
		} else if op == 3 {
			DelUser()
		} else if op == 4 {
			AddGroup()
		} else if op == 5 {
			UpdGroup()
		} else if op == 6 {
			DelGroup()
		} else if op == 7 {
			AddSpecialist()
		} else if op == 8 {
			DelSpecialist()
		} else if op == 9 {
			AddGroupAdUser()
		} else if op == 10 {
			DelGroupAdUser()
		}
	}
}

func printBody(r *h.Response) {
	body, e := ioutil.ReadAll(r.Body)
	if e == nil {
		r.Body.Close()
		fmt.Println(string(body))
	}
}
