package cmd

import (
	"fmt"
	s3 "github.com/dicksontung/cloudflare-dns-backup/awsS3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func rootC(cmd *cobra.Command, args []string) {
	zoneIDs := viper.GetStringSlice("zone")
	for _, z := range zoneIDs {
		if !download(z) {
			return
		}
	}
	uploadAll()
}

func download(zoneID string) bool {
	cfURL := strings.ReplaceAll(cloudflareURL, ":zone_id", zoneID)
	fmt.Printf("Exporting from: %+v \n", cfURL)
	req, err := http.NewRequest(http.MethodGet, cfURL, nil)
	if err != nil {
		fmt.Printf("unable to create request: %+v \n", err)
		return false
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("token")))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("unable to perform request: %+v \n", err)
		return false
	}
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("error code: %d, %s \n", resp.StatusCode, b)
		return false
	}

	fileName := "out/" + viper.GetString("prefix") + zoneID + ".txt"
	os.Mkdir("out/", os.ModePerm)
	out, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("unable to create file %q: %+v \n", fileName, err)
		return false
	}
	defer func() {
		if err := out.Close(); err != nil {
			fmt.Printf("error when closing: %+v \n", err)
		}
	}()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Unable to write file: %+v \n", err)
		return false
	}
	fmt.Printf("Write to file %q, %d bytes \n", fileName, n)
	return true
}

func uploadAll() {
	files, err := ioutil.ReadDir("out")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		s3.Upload("out/" + f.Name())
	}
}