package frontend

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

func RenderHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is printing")
	fp := path.Join("template", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderReset(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("template", "reset.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleReset(writer http.ResponseWriter, request *http.Request) {
	// Send Delete request

	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	}

	fmt.Println(given_xAuth)

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/reset", nil)
	if err != nil {
		fmt.Println("request error")
		return
	}
	req.Header.Add("X-Authorization", given_xAuth)
	resp, err := client.Do(req)

	fmt.Print(resp.Status)

	writer.Write([]byte(string(resp.Status)))
}

func RenderPUTPackage(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("template", "putPackage.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandlePUTPackage(writer http.ResponseWriter, request *http.Request) {
	// Send Delete request

	request.ParseForm()

	var given_xAuth string
	var id string
	var body string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	}
	if request.Form["id"] != nil {
		id = request.Form["id"][0]
	}
	if request.Form["Content"] != nil {
		body = request.Form["Content"][0]
	}

	client := &http.Client{}

	req, err := http.NewRequest("PUT", "http://localhost:8080/package/"+id, strings.NewReader(body))
	if err != nil {
		fmt.Println("request error")
		return
	}

	req.Header.Add("X-Authorization", given_xAuth)
	req.Header.Add("id", id)
	req.Header.Add("Content", body)

	resp, err := client.Do(req)

	fmt.Print(resp.Status)
	writer.Write([]byte(string(resp.Status)))
}

func RenderAuthenticatePackage(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("template", "authenticate.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleAuthenticatePackage(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	var username string
	var password string

	fmt.Print("am i here\n")

	if request.Form["Username"] != nil {
		username = request.Form["Username"][0]
	}
	if request.Form["Password"] != nil {
		password = request.Form["Password"][0]
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", "https://tomr-g17-mdljbaftcq-uc.a.run.app/authenticate", nil)
	if err != nil {
		fmt.Println("request error")
		return
	}

	req.Header.Add("Username", username)
	req.Header.Add("Password", password)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte

	writer.Write([]byte(string(resp.Status) + "\n"))
	writer.Write([]byte(string(body)))

}
