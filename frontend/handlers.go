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
	// fmt.Println("this is printing")
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

	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte

	writer.Write([]byte(string(resp.Status) + "\n"))
	writer.Write([]byte(string(body)))
}

func RenderGETPackage(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("template", "getPackage.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleGETPackage(writer http.ResponseWriter, request *http.Request) {
	// Send Delete request

	request.ParseForm()

	var given_xAuth string
	var id string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	}
	if request.Form["id"] != nil {
		id = request.Form["id"][0]
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8080/package/"+id, nil)
	if err != nil {
		fmt.Println("request error")
		return
	}

	req.Header.Add("X-Authorization", given_xAuth)
	req.Header.Add("id", id)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body) // response body is []byte

	writer.Write([]byte(string(resp.Status) + "\n"))
	writer.Write([]byte(string(resp_body)))
}

func RenderPackage(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("template", "postPackage.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandlePackage(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	var given_xAuth string
	var url string
	var content string
	var jsprogram string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	}
	if request.Form["URL"] != nil {
		url = request.Form["URL"][0]
	}
	if request.Form["Content"] != nil {
		content = request.Form["Content"][0]
	}
	if request.Form["JSProgram"] != nil {
		jsprogram = request.Form["JSProgram"][0]
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:8080/package/", nil)
	if err != nil {
		fmt.Println("request error")
		return
	}

	req.Header.Add("X-Authorization", given_xAuth)
	req.Header.Add("Url", url)
	req.Header.Add("Content", content)
	req.Header.Add("Jsprogram", jsprogram)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body) // response body is []byte

	writer.Write([]byte(string(resp.Status) + "\n"))
	writer.Write([]byte(string(resp_body)))

}

func RenderDELETEPackage(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("template", "deletePackage.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleDELETEPackage(writer http.ResponseWriter, request *http.Request) {
	// Send Delete request

	request.ParseForm()

	var given_xAuth string
	var id string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	}
	if request.Form["id"] != nil {
		id = request.Form["id"][0]
	}

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/package/"+id, nil)
	if err != nil {
		fmt.Println("request error")
		return
	}

	req.Header.Add("X-Authorization", given_xAuth)
	req.Header.Add("id", id)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body) // response body is []byte

	writer.Write([]byte(string(resp.Status) + "\n"))
	writer.Write([]byte(string(resp_body)))
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

	req, err := http.NewRequest("PUT", "https://tomr-g17-mdljbaftcq-uc.a.run.app/package/"+id, strings.NewReader(body))
	if err != nil {
		fmt.Println("request error")
		return
	}

	req.Header.Add("X-Authorization", given_xAuth)
	req.Header.Add("id", id)
	req.Header.Add("Content", body)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body) // response body is []byte

	writer.Write([]byte(string(resp.Status) + "\n"))
	writer.Write([]byte(string(resp_body)))
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

func RenderPackages(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("template", "packages.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandlePackages(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	var given_xAuth string
	var full_body string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	}
	if request.Form["bodyfull"] != nil {
		full_body = request.Form["bodyfull"][0]
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/packages", strings.NewReader(full_body))
	if err != nil {
		fmt.Println("request error")
		return
	}

	req.Header.Add("X-Authorization", given_xAuth)
	// req.Header.Add("Fullbody", full_body)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte

	writer.Write([]byte(string(resp.Status) + "\n"))
	writer.Write([]byte(string(body)))
}

func RenderRegex(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("template", "regex.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleRegex(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	var given_xAuth string
	var regex_str string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	}
	if request.Form["Regex"] != nil {
		regex_str = request.Form["Regex"][0]
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/package/byRegEx", nil)
	if err != nil {
		fmt.Println("request error")
		return
	}

	req.Header.Add("X-Authorization", given_xAuth)
	req.Header.Add("Regex", regex_str)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte

	writer.Write([]byte(string(resp.Status) + "\n"))
	writer.Write([]byte(string(body)))
}
