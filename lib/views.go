package lib

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
)

type ViewsCreator struct {
	bqClient *bigquery.Client
}

func InitViewsCreator() (*ViewsCreator, error) {
	bqClient, err := bigquery.NewClient(context.Background(), "streamrail")
	if err != nil {
		return nil, err
	}
	return &ViewsCreator{
		bqClient: bqClient,
	}, nil
}

func (vc *ViewsCreator) Start() error {
	orgIds, err := LoadOrgIDs()
	if err != nil {
		return err
	}
	err = vc.updateRawView(orgIds, "yesterday")
	if err != nil {
		return err
	}
	err = vc.updateRawView(orgIds, "today")
	if err != nil {
		return err
	}
	err = vc.updateRawView(orgIds, "last hour")
	if err != nil {
		return err
	}
	err = vc.updateRawView(orgIds, "last 3 hours")
	if err != nil {
		return err
	}
	err = vc.updateRawView(orgIds, "last 30 minutes")
	if err != nil {
		return err
	}
	err = vc.updateAggView(orgIds, "yesterday")
	if err != nil {
		return err
	}
	err = vc.updateAggView(orgIds, "today")
	if err != nil {
		return err
	}

	return nil
}

func (vc *ViewsCreator) updateRawView(orgs []string, offset string) error {
	var viewName string
	var days int
	var hours float64
	var useDecorator bool
	switch offset {
	case "yesterday":
		viewName = "RAW_all_orgs_yesterday"
		days = 1
		useDecorator = false
	case "today":
		viewName = "RAW_all_orgs_today"
		days = 0
		useDecorator = false
	case "last hour":
		viewName = "RAW_all_orgs_last_hour"
		useDecorator = true
		hours = 1
	case "last 3 hours":
		viewName = "RAW_all_orgs_last_3_hours"
		useDecorator = true
		hours = 3
	case "last 30 minutes":
		viewName = "RAW_all_orgs_last_30_minutes"
		useDecorator = true
		hours = 0.5
	default:
		return fmt.Errorf("Cannot update view without specifiying a correct time offset")
	}

	var buff bytes.Buffer
	queryFile, err := ioutil.ReadFile("/etc/raw_view.sql")
	if err != nil {
		return err
	}
	if !useDecorator {
		for _, orgId := range orgs {
			setRawTemplate(orgId, days, &buff)
		}
	} else {
		for _, orgId := range orgs {
			setRawTemplateWithDecorator(orgId, hours, &buff)
		}
	}
	q := strings.Replace(string(queryFile), "{{TABLES}}", buff.String(), 1)
	return vc.updateView(viewName, q)
}

func (vc *ViewsCreator) updateAggView(orgs []string, offset string) error {
	var viewName string
	var days int
	switch offset {
	case "yesterday":
		viewName = "AGG_BIG_all_orgs_yesterday"
		days = 1
	case "today":
		viewName = "AGG_BIG_all_orgs_today"
		days = 0
	default:
		return fmt.Errorf("Cannot update view without specifiying a correct time offset")
	}

	var buff bytes.Buffer
	queryFile, err := ioutil.ReadFile("/etc/agg_view.sql")
	if err != nil {
		return err
	}
	for _, orgId := range orgs {
		setAggTemplate(orgId, &buff)
	}
	q := strings.Replace(string(queryFile), "{{TABLES}}", buff.String(), 1)

	startDate := time.Now().UTC().AddDate(0, 0, -1*days).Format("2006-01-02 00:00:00")
	endDate := time.Now().UTC().AddDate(0, 0, -1*days+1).Format("2006-01-02 00:00:00")
	q2 := strings.Replace(q, "{{START_TIMESTAMP}}", startDate, 1)
	q3 := strings.Replace(q2, "{{END_TIMESTAMP}}", endDate, 1)

	return vc.updateView(viewName, q3)
}

func setRawTemplate(org string, days int, buff *bytes.Buffer) {
	daysString := strconv.Itoa(days)
	buff.WriteString("TABLE_DATE_RANGE([")
	buff.WriteString(org)
	buff.WriteString(".adsmanager_],DATE_ADD(CURRENT_TIMESTAMP(),-")
	buff.WriteString(daysString)
	buff.WriteString(",'DAY'), DATE_ADD(CURRENT_TIMESTAMP(),-")
	buff.WriteString(daysString)
	buff.WriteString(",'DAY') ), ")
}

func setRawTemplateWithDecorator(org string, hours float64, buff *bytes.Buffer) {
	decoratorString := strconv.Itoa(int(hours * 3600000))
	todayDate := time.Now().UTC().Format("20060102")
	buff.WriteString("[")
	buff.WriteString(org)
	buff.WriteString(".adsmanager_")
	buff.WriteString(todayDate)
	buff.WriteString("@-")
	buff.WriteString(decoratorString)
	buff.WriteString("-],")
}

func setAggTemplate(org string, buff *bytes.Buffer) {
	buff.WriteString("[")
	buff.WriteString(org)
	buff.WriteString(".adsmanager_agg_big],")
}

func (vc *ViewsCreator) updateView(viewName, query string) error {
	log.Printf("Query is %s", query)
	tableMetadataToUpdate := &bigquery.TableMetadataToUpdate{
		Schema:    nil,
		ViewQuery: query,
	}
	_, err := vc.bqClient.Dataset("views").Table(viewName).Update(context.Background(), *tableMetadataToUpdate, "")
	if err != nil {
		return err
	}
	return nil
}
