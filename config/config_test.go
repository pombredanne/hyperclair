package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/test"

	"gopkg.in/yaml.v2"
)

var loginData = []struct {
	in  string
	out int
}{
	{"", 0},
	{`{
        "docker.io": {
                "Username": "johndoe",
                "Password": "$2a$05$Qe4TTO8HMmOht"
        }
}
`, 1},
}

const defaultValues = `
clair:
  uri: http://localhost
  priority: Low
  port: 6060
  healthport: 6061
  report:
    path: reports
    format: html
auth:
  insecureskipverify: true
hyperclair:
  ip: ""
  tempfolder: /tmp/hyperclair
  port: 0
  interface: native
`

const customValues = `
clair:
  uri: http://clair
  priority: High
  port: 6061
  healthport: 6062
  report:
    path: reports/test
    format: json
auth:
  insecureskipverify: false
hyperclair:
  ip: "localhost"
  tempfolder: /tmp/hyperclair/test
  port: 64157
  interface: virtualbox
`

func TestInitDefault(t *testing.T) {
	Init("", "INFO")

	cfg := values()

	var expected config
	err := yaml.Unmarshal([]byte(defaultValues), &expected)
	if err != nil {
		t.Fatal(err)
	}

	if cfg != expected {
		t.Error("Default values are not correct")
	}
	viper.Reset()
}

func TestInitCustomLocal(t *testing.T) {
	tmpfile := test.CreateConfigFile(customValues, "hyperclair.yml", ".")
	defer os.Remove(tmpfile) // clean up
	fmt.Println(tmpfile)
	Init("", "INFO")

	cfg := values()

	var expected config
	err := yaml.Unmarshal([]byte(customValues), &expected)
	if err != nil {
		t.Fatal(err)
	}

	if cfg != expected {
		t.Error("values are not correct")
	}
	viper.Reset()
}

func TestInitCustomHome(t *testing.T) {
	tmpfile := test.CreateConfigFile(customValues, "hyperclair.yml", HyperclairHome())
	defer os.Remove(tmpfile) // clean up
	fmt.Println(tmpfile)
	Init("", "INFO")

	cfg := values()

	var expected config
	err := yaml.Unmarshal([]byte(customValues), &expected)
	if err != nil {
		t.Fatal(err)
	}

	if cfg != expected {
		t.Error("values are not correct")
	}
	viper.Reset()
}

func TestInitCustom(t *testing.T) {
	tmpfile := test.CreateConfigFile(customValues, "hyperclair.yml", "/tmp")
	defer os.Remove(tmpfile) // clean up
	fmt.Println(tmpfile)
	Init(tmpfile, "INFO")

	cfg := values()

	var expected config
	err := yaml.Unmarshal([]byte(customValues), &expected)
	if err != nil {
		t.Fatal(err)
	}

	if cfg != expected {
		t.Error("values are not correct")
	}
	viper.Reset()
}

func TestReadConfigFile(t *testing.T) {
	for _, ld := range loginData {

		tmpfile := test.CreateTmpConfigFile(ld.in)
		defer os.Remove(tmpfile) // clean up

		var logins loginMapping
		if err := readConfigFile(&logins, tmpfile); err != nil {
			t.Errorf("readConfigFile(&logins,%q) failed => %v", tmpfile, err)
		}

		if l := len(logins); l != ld.out {
			t.Errorf("readConfigFile(&logins,%q) => %v logins, want %v", tmpfile, l, ld.out)
		}
	}
}

func TestWriteConfigFile(t *testing.T) {
	logins := loginMapping{}
	logins["docker.io"] = Login{Username: "johndoe", Password: "$2a$05$Qe4TTO8HMmOht"}
	tmpfile := test.CreateTmpConfigFile("")
	defer os.Remove(tmpfile) // clean up

	if err := writeConfigFile(logins, tmpfile); err != nil {
		t.Errorf("writeConfigFile(logins,%q) failed => %v", tmpfile, err)
	}

	logins = loginMapping{}
	if err := readConfigFile(&logins, tmpfile); err != nil {
		t.Errorf("after writing: readConfigFile(&logins,%q) failed => %v", tmpfile, err)
	}

	if l := len(logins); l != 1 {
		t.Errorf("after writing: readConfigFile(&logins,%q) => %v logins, want %v", tmpfile, l, 1)
	}
}
