package cmd

import (
	"fmt"
	"github.com/origin-tech/quick-tricks/modules/ssrf"
	"github.com/origin-tech/quick-tricks/utils/colors"
	"github.com/spf13/cobra"
)

// ssrfCmd represents the ssrf command.
var ssrfCmd = &cobra.Command{
	Use:   "ssrf",
	Short: "Module 'ssrf' helps to check whether the target is vulnerable to SSRF or not.",
	Run: func(cmd *cobra.Command, args []string) {
		target, _ := cmd.Flags().GetString("url")
		server, _ := cmd.Flags().GetString("server")
		proxy, _ := cmd.Flags().GetString("proxy")

		ssrfSuccessed, errors := ssrf.Detect(target, server, proxy)
		if len(errors) != 0 {
			for _, e := range errors {
				colors.BAD.Println(e)
			}
			return
		}
		if len(ssrfSuccessed) != 0 {
			colors.OK.Printf("Sent SSRF requests:\n")
			for _, d := range ssrfSuccessed {
				fmt.Printf("----------------------------------------\n")
				reqTypeFiled := "Request type: POST"
				fmt.Println(reqTypeFiled)
				urlFiled := fmt.Sprintf("URL: %s", d[0])
				fmt.Println(urlFiled)
				bodyField := fmt.Sprintf("Request payload: %s", d[1])
				fmt.Println(bodyField)
				colors.NEUTRAL.Printf("Check logs of %s HTTP server to see if the request reached the destination.\n", server)

			}
		} else {
			colors.BAD.Println("No available SSRF endpoints are found.")
		}
	},
}

func init() {
	rootCmd.AddCommand(ssrfCmd)
	ssrfCmd.Flags().StringP("url", "u", "", "Target Bitrix site")
	ssrfCmd.MarkFlagRequired("url")
	ssrfCmd.Flags().StringP("server", "s", "", "External host with running HTTP server (with protocol type: http/https")
	ssrfCmd.MarkFlagRequired("server")
	ssrfCmd.Flags().String("proxy", "", "http/socks5 proxy to use. Example: socks5://IP:PORT")
}
