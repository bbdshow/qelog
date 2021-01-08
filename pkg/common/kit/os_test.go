package kit

import "testing"

func TestGetLocalIPV4(t *testing.T) {
	ip, err := GetLocalIPV4()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
}

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
}
