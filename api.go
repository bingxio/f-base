// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Api struct{} // Bridge

func logging(m string) { log.Printf("\t%s", m) }

// OpenApiServer : Monitor local port
func OpenApiServer() error {
	logging("F-Base runs in the background and monitors port: " + strconv.Itoa(port))
	logging("listening... ")

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), Api{})
	return err
}

// ServeHTTP : Serve http requests
func (a Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logging("GET: " + r.RequestURI)
	switch r.RequestURI {
	case "/version":
		_, _ = w.Write([]byte(Version))
	case "/save":
		err := GlobalMem.Write()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte("OK"))
	case "/reload":
		err := GlobalEm.Reload()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte("OK"))
	case "/t":
		var elem []interface{}
		for _, v := range GlobalMem.Tb {
			elem = append(elem, struct {
				Name     string `json:"name"`
				At       uint8  `json:"at"`
				Rows     uint64 `json:"rows"`
				CreateAt uint32 `json:"create_at"`
			}{
				Name:     v.NameToString(),
				At:       v.At,
				Rows:     v.Rows,
				CreateAt: v.Created,
			})
		}
		b, err := json.Marshal(elem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(b)
	case "/p":
		var elem []interface{}
		for i := 0; i < len(GlobalMem.Tr); i++ {
			x, y := GlobalMem.Tr[i].Counts()
			elem = append(elem, struct {
				Name string `json:"name"`
				Node int    `json:"node"`
				Leaf int    `json:"leaf"`
				Data int    `json:"data"`
			}{
				Name: GlobalMem.Tb[i].NameToString(),
				Node: len(GlobalMem.Tr[i].Node),
				Leaf: x,
				Data: y,
			})
		}
		b, err := json.Marshal(elem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(b)
	}
	if strings.HasPrefix(r.RequestURI, "/se") {
		content := r.URL.Query().Get("content")
		tb := r.URL.Query().Get("tb")
		if content == "" || tb == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err, content := DecodeContent(content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		var field []string
		for i := 0; i < len(content); i++ {
			x := ""
			for i < len(content) && content[i] != 32 {
				x += string(content[i])
				i += 1
			}
			field = append(field, x)
		}
		i, exist := GlobalEm.Exist(tb)
		if !exist {
			GlobalEm.Tb = append(GlobalEm.Tb, GlobalMem.NewTable(tb))
			i = len(GlobalEm.Tb) - 1
		}
		r := GlobalEm.Tb[i].Insert(field).(SingleResult)
		b, err := json.Marshal(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(b)
	}
	if strings.HasPrefix(r.RequestURI, "/ge") {
		tb := r.URL.Query().Get("tb")
		if tb == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		f := r.URL.Query().Get("f")
		t := r.URL.Query().Get("t")
		i, exist := GlobalEm.Exist(tb)
		if !exist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		r, err := GlobalEm.Tb[i].Select(f, t)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		b, err := json.Marshal(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(b)
	}
	if strings.HasPrefix(r.RequestURI, "/gt") {
		tb := r.URL.Query().Get("tb")
		s := r.URL.Query().Get("s")
		v := r.URL.Query().Get("v")
		if tb == "" || s == "" || v == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		i, exist := GlobalEm.Exist(tb)
		if !exist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		r, err := GlobalEm.Tb[i].Selector(s, v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		b, err := json.Marshal(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(b)
	}
	if strings.HasPrefix(r.RequestURI, "/up") {
		tb := r.URL.Query().Get("tb")
		p := r.URL.Query().Get("p")
		s := r.URL.Query().Get("s")
		n := r.URL.Query().Get("n")
		v := r.URL.Query().Get("v")
		if tb == "" || p == "" || s == "" || n == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		i, exist := GlobalEm.Exist(tb)
		if !exist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		r, err := GlobalEm.Tb[i].Update(p, s, n, v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		b, err := json.Marshal(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(b)
	}
	if strings.HasPrefix(r.RequestURI, "/de") {
		tb := r.URL.Query().Get("tb")
		if tb == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		p := r.URL.Query().Get("p")
		i, exist := GlobalEm.Exist(tb)
		if !exist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		r, err := GlobalEm.Tb[i].Delete(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		b, err := json.Marshal(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(b)
	}
}

// DecodeContent : Decode base decoding to string literal
func DecodeContent(base string) (error, string) {
	r, err := base64.StdEncoding.DecodeString(base)
	if err != nil {
		return err, ""
	}
	return nil, string(r)
}
