package lib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/streamrail/views/models"
)

const (
	dsExportBucketAWSOhio = "https://s3.us-east-2.amazonaws.com/half-pipe-bq-ohio.streamrail.com/"
	orgsFileName          = "entities/orgs.json.gz"
)

func DownloadFile(url string) ([]byte, error) {
	client := &http.Client{}
	resp, err := getWithRetry(client, url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if strings.Contains(string(data), "Access Denied") {
		return nil, fmt.Errorf("Failed downloading file %v: Access denied.", url)
	}
	log.Printf("DataLoader: Successfully downloaded file: %s", url)
	return data, nil
}

func getWithRetry(client *http.Client, url string) (resp *http.Response, err error) {
	for i := 1; ; i++ {
		resp, err = client.Get(url)
		if err == nil {
			return
		}
		log.Printf("DataLoader: Fail to download %s, attempt %d. Error: %v", url, i, err)
		if i >= 10 {
			return
		}
		time.Sleep(2 * time.Second)
	}
}

func Unzip(data io.Reader) (unzippedData []byte, err error) {
	unzippedReader, err := gzip.NewReader(data)
	if err != nil {
		return nil, err
	}
	unzippedData, err = ioutil.ReadAll(unzippedReader)
	if err != nil {
		return nil, err
	}
	return unzippedData, nil
}

func LoadOrgIDs() ([]string, error) {
	url := dsExportBucketAWSOhio + orgsFileName
	gzippedOrgs, err := DownloadFile(url)
	if err != nil {
		return nil, fmt.Errorf("Error downloading orgIDs json file: Url: %v, error: %v", url, err)
	}
	zippedReader := bytes.NewReader(gzippedOrgs)
	orgsJson, err := Unzip(zippedReader)
	if err != nil {
		return nil, fmt.Errorf("Error unzipping orgIDs json file: %v", err)
	}
	var orgs []models.Org
	json.Unmarshal(orgsJson, &orgs)
	if err != nil {
		return nil, fmt.Errorf("Error in unmarshal orgIDs json file: %v", err)
	}
	orgIDs := make([]string, 0, len(orgs))
	for _, orgEntity := range orgs {
		if orgEntity.Id == "" {
			continue
		}
		orgIDs = append(orgIDs, orgEntity.Id)
	}
	return orgIDs, nil
}
