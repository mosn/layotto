package integrate

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSayHello(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:2045", nil)
	name := "Layotto"
	req.Header.Add("name",name)
	resp,_ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), "Hi, " + name)
}
