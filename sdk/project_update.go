// project_update automatically sdk by SDKgen please do not edit this file manually
// @see https://sdkgen.app

package sdk

import "time"

type ProjectUpdate struct {
	Id            string      `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	BaseUrl       string      `json:"baseUrl"`
	Namespace     string      `json:"namespace"`
	Versions      []Version   `json:"versions"`
	LatestVersion interface{} `json:"latestVersion"`
	InsertDate    time.Time   `json:"insertDate"`
}
