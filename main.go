package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"fmt"
	"strings"
	"io/ioutil"
	"strconv"
	"log"
	"./model"
	"github.com/gorilla/mux"
)

const PORT = "8000" 

var (
	Work map[string]int32
	Auth map[string]int64
)

func main() {
	fmt.Println("start main\n")
	r := mux.NewRouter()

	r.HandleFunc("/", mainPage)
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/changepass", changePass).Methods("POST")
	r.HandleFunc("/dowork", doWork).Methods("POST")

	http.Handle("/", r)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":" + PORT, r))
	
	fmt.Println("end main\n")
}

func init() {
	Work = make(map[string]int32, 2)
	Auth = make(map[string]int64, 2)
	Work["-1"] = -1
	Auth["-1"] = -1
	model.GormInit();
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
		<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta name="theme-color" content="#375EAB">
		
			<title>main page</title>
		</head>
		<body>
			Page body and some more content
		</body>
		</html>`))
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
 // Create return string
 var request []string
 // Add the request string
 url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
 request = append(request, url)
 // Add the host
 request = append(request, fmt.Sprintf("Host: %v", r.Host))
 // Loop through headers
 for name, headers := range r.Header {
   name = strings.ToLower(name)
   for _, h := range headers {
     request = append(request, fmt.Sprintf("%v: %v", name, h))
   }
 }
 
 // If this is a POST, add post data
 if r.Method == "POST" {
 	//fmt.Println("STAGE POST START\n")
    r.ParseForm()
    request = append(request, "\n")
    request = append(request, r.Form.Encode())
 } 
  // Return the request as a string
  return strings.Join(request, "\n")
}

func check_login(login string,w http.ResponseWriter) bool {
	//return true
	_, login_key_found_ := Auth[login]
	if !login_key_found_ || login == "-1"  {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DENIED"))
		return false
	} else {
		return true
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	pass := r.FormValue("pass")
	
	if ((login != "") && (pass != "")) {
		
		user := &model.Usr{}
		if user != nil {
			err := user.Get(login, pass)
			if err == nil {			
				Auth[login] = user.ID
				Work[login] = user.WorkNumber
			} else {
				if err.Error() == "record not found" {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("DENIED"))
					return
				}
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func changePass(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	newPass := r.FormValue("newPass")
	if ((login != "") && (newPass != "")) {
		if !check_login(login,w) {
			return
		}
		user := &model.Usr{}
		if user != nil {
			err := user.Save(Auth[login],newPass)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				delete(Auth, login)
				delete(Work, login)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("PASSWORD IS CHANGED"))
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

type DTO struct {
	bigNumber int64
	Number int32
	Name      string
}
//+
func doWork(w http.ResponseWriter, r *http.Request) {
	var value DTO
	login := r.FormValue("login")
	if (login == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
			if !check_login(login,w) {
				return
			}
	}
	body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
		return
    }

    if &value != nil {
		err = json.Unmarshal([]byte(body), &value)
		if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Println(err)
		return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
        w.WriteHeader(http.StatusBadRequest)
		return
    }
    data := []byte("{")
    w.Header().Set("Content-Type", "application/json")
	v := reflect.ValueOf(&value).Elem()
	for i := 0; i < v.NumField(); i++ {
		field_v := v.Field(i)
		if field_v.IsValid() {
			tmp := reverse(field_v)
			data = append(data,tmp...)
			data = append(data,([]byte(","))...)
		}
	}
	data[len(data)-1] = '\x7D'
	w.Write(data)
}

func reverse(val reflect.Value) []byte {
	switch val.Kind().String() {
	case "int64":
		//fallthrough
		key_ := "\"bigNumber\":"
		return []byte(key_ + strconv.FormatUint(uint64(9223372036854775807-uint64(val.Int())), 10))
	case "int32":
		key_ := "\"Number\":"
		return []byte(key_ + strconv.FormatUint(uint64(2147483647-uint32(val.Int())), 10))
	case "string":
		key_ := "\"Name\":\""
		var result string
		str_tmp := val.String()
		if str_tmp != "" {
			for i := len(str_tmp); i >= 1; i-- {
				result += string(str_tmp[i-1])
			}
		}
		result = key_ + result + "\""
		return []byte(result)
	}
	return nil
}
