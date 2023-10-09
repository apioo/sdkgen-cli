// CommonMessageException automatically sdk by SDKgen please do not edit this file manually
// @see https://sdkgen.app

package sdk

import (
	"encoding/json"
	"fmt"
)

type CommonMessageException struct {
	Payload CommonMessage
}

func (e *CommonMessageException) Error() string {
	raw, err := json.Marshal(e.Payload)
	if err != nil {
		return "could not marshal provided JSON data"
	}

	return fmt.Sprintf("The server returned an error: %s", raw)
}
