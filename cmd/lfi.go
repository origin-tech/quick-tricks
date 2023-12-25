package cmd

import (
	"fmt"
	"github.com/origin-tech/quick-tricks/modules/lfi"
	"github.com/origin-tech/quick-tricks/utils/colors"

	"github.com/spf13/cobra"
)

// lfiCmd represents the spoofing command
var lfiCmd = &cobra.Command{
	Use:   "lfi",
	Short: "Module 'lfi' checks if there are endpoints vulnerable to Local File Inclusion.",
	Run: func(cmd *cobra.Command, args []string) {
		target, _ := cmd.Flags().GetString("url")
		proxy, _ := cmd.Flags().GetString("proxy")

		lfiUrls, err := lfi.Detect(target, proxy)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(lfiUrls) != 0 {
			colors.OK.Printf("Path to LFI:\n")
			for _, v := range lfiUrls {
				fmt.Println(v)
			}
		} else {
			colors.BAD.Println("There is no path to LFI.")
		}
	},
}

func init() {
	rootCmd.AddCommand(lfiCmd)
	lfiCmd.Flags().StringP("url", "u", "", "Target Bitrix site")
	lfiCmd.MarkFlagRequired("url")
	lfiCmd.Flags().String("proxy", "", "http/socks5 proxy to use. Example: socks5://IP:PORT")
}
