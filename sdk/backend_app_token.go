// backend_app_token automatically generated by SDKgen please do not edit this file manually
// @see https://sdkgen.app

package sdk

import "time"

type BackendAppToken struct {
	Id     int       `json:"id"`
	Status int       `json:"status"`
	Token  string    `json:"token"`
	Scope  []string  `json:"scope"`
	Ip     string    `json:"ip"`
	Expire time.Time `json:"expire"`
	Date   time.Time `json:"date"`
}
