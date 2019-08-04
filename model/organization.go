package model

//Endpoint services endpoints
type Endpoint struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

//Environment variables
type Environment struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

//Subscription types
type Subscription struct {
	Type  string `json:"type" bson:"type"` //Free, Silver, Gold, Platinum, Enterprise
	Limit struct {
		Endpoint uint   `json:"endpoint" bson:"endpoint"` //total registered endpoint
		Scenario uint   `json:"scenario" bson:"scenario"` //total scenario
		Cases    uint   `json:"cases" bson:"cases"`       //total cases regardless of scenario
		Data     uint64 `json:"data" bson:"data"`         //total bytes downloaded
	} `json:"limits" bson:"limit"`
}

//Organization info
type Organization struct {
	ID           string        `json:"-" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	Alias        string        `json:"alias" bson:"alias"`
	Endpoints    []Endpoint    `json:"endpoints" bson:"endpoints"`
	Environments []Environment `json:"environments" bson:"environments"`
	Subscription Subscription  `json:"-" bson:"subscription"`
}

//Validate update
func (org *Organization) Validate() error {
	/* TODO
	messages := []string{}
	for i, e := range org.Endpoints {

	}
	*/
	return nil
}
