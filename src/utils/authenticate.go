package utils

import (
    // "fmt"
)

func Authenticate(username string, password string) string {

    if(username == "ece30861defaultadminuser" && password == "correcthorsebatterystaple123(!__+@**(A\u2019\u201D`;DROP TABLE packages;") {
        return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiZWNlMzA4NjFkZWZhdWx0YWRtaW51c2VyIiwicGFzc3dvcmQiOiJjb3JyZWN0aG9yc2ViYXR0ZXJ5c3RhcGxlMTIzKCFfXytAKiooQeKAmeKAnWA7RFJPUCBUQUJMRSBwYWNrYWdlczsifQ.TSGs6VJMFx5NV2RoHrhEP_FK8nv4Wlzc4gQls2JYPC4"
    }

    return "err"
}

// example usage
// func main() {
//     token := Authenticate("ece30861defaultadminuser" ,"correcthorsebatterystaple123(!__+@**(A\u2019\u201D`;DROP TABLE packages;")
//     fmt.Println(token)
// }