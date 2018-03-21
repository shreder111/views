package models

import "time"

type Org struct {
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Status      int       `json:"status,omitempty"`
	CreatedOn   time.Time `json:"createdOn,omitempty"`
	ModifiedOn  time.Time `json:"lastModified,omitempty"`
	CreatedBy   string    `json:"createdBy,omitempty"`
	ModifiedBy  string    `json:"modifiedBy,omitempty"`
	EndCardType int       `json:"endCardType,omitempty"`

	MoatID       string  `json:"moatId"`
	MoatSampling float64 `json:"moatSampling"`

	DesktopDomainBlackList           string `json:"desktopDomainBlackList,omitempty"`
	DesktopDomainBlackListSubDomains bool   `json:"desktopDomainBlackListSubDomains,omitempty"`

	MobileWebDomainBlackList           string `json:"mobileWebDomainBlackList,omitempty"`
	MobileWebDomainBlackListSubDomains bool   `json:"mobileWebDomainBlackListSubDomains,omitempty"`

	MobileAppBundleIdBlackList string `json:"mobileAppBundleIdBlackList,omitempty"`

	TimeZone string `json:"timeZone,omitempty"`
}
