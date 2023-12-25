// Package lp provides function to check possible Bitrix login pages.
package lp

import (
	"fmt"
	"github.com/origin-tech/quick-tricks/utils/netclient"
)

const (
	endpoint1 = "/bitrix/admin"
	endpoint2 = "/bitrix/components/bitrix/desktop/admin_settings.php"
	endpoint3 = "/bitrix/components/bitrix/map.yandex.search/settings/settings.php"
	endpoint4 = "/bitrix/components/bitrix/player/player_playlist_edit.php"
	endpoint5 = "/bitrix/tools/autosave.php"
	endpoint6 = "/bitrix/tools/get_catalog_menu.php"
	endpoint7 = "/bitrix/tools/upload.php"
	endpoint8 = "/?SEF_APPLICATION_CUR_PAGE_URL=/bitrix/admin/"
)

func Detect(target, proxy string) ([][]string, error) {
	var pages [][]string
	endpoints := []string{endpoint1, endpoint2, endpoint3, endpoint4, endpoint5, endpoint6, endpoint7, endpoint8}

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
			var page []string
			page = append(page, url)
			page = append(page, resp.Status)
			pages = append(pages, page)
		}
	}

	return pages, nil
}
