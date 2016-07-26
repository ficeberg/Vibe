package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

// NOTE : struct doesn't work here due to the nature
// of the JSON data from country.io
// it has no field identifiers

type Country struct {
	Map map[string]interface{}
	Log Logger
}

func (c *Country) Init() error {
	_, filename, _, _ := runtime.Caller(1)
	meta := map[string]interface{}{
		"from":    filename,
		"section": "GetCountry",
		"time":    time.Now(),
	}

	resp, err := http.Get("http://country.io/names.json")
	if err != nil {
		c.Log.Error(meta, "Failed to fetch from country.io")
		return err
	}
	defer resp.Body.Close()

	jsonCountriesData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		meta["section"] = "ReadRespBody"
		c.Log.Error(meta, err.Error())
		return err
	}

	//fmt.Println(string(jsonData))
	// Decode JSON into our map
	err = json.Unmarshal([]byte(jsonCountriesData), &c.Map)
	if err != nil {
		meta["section"] = "ParseFromJson"
		c.Log.Error(meta, err.Error())
		return err
	}

	/*for iso2, name := range c.Map {
		fmt.Println("ISO2 code:", iso2, "Country name :", name)
	}*/
	return nil
}

func (c *Country) Iso2Country(iso string) string {
	if err := c.Init(); err == nil {
		return c.Map[iso].(string)
	}
	return ""
}

func (c *Country) Country2Iso(country string) string {
	if err := c.Init(); err == nil {
		for iso2, name := range c.Map {
			if name == country {
				return iso2
			}
		}
	}
	return ""
}
