package goseaweed

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type Location struct {
	Url       string `json:"url,omitempty"`
	PublicUrl string `json:"publicUrl,omitempty"`
}

type Locations []Location

type LookupResult struct {
	VolumeId  string    `json:"volumeId,omitempty"`
	Locations Locations `json:"locations,omitempty"`
	Error     string    `json:"error,omitempty"`
}

func (lr *LookupResult) String() string {
	return fmt.Sprintf("VolumeId:%s, Locations:%v, Error:%s", lr.VolumeId, lr.Locations, lr.Error)
}

func (ls Locations) Head() *Location {
	return &ls[0]
}

func (ls Locations) PickForRead() *Location {
	return &ls[rand.Intn(len(ls))]
}

func (sw *Seaweed) Lookup(vid, collection string) (ret *LookupResult, err error) {
	locations, cache_err := sw.vc.Get(vid)
	if cache_err != nil {
		if ret, err = sw.doLookup(vid, collection); err == nil {
			sw.vc.Set(vid, ret.Locations, 10*time.Minute)
		}
	} else {
		ret = &LookupResult{VolumeId: vid, Locations: locations}
	}
	return
}

func (sw *Seaweed) LookupNoCache(vid, collection string) (ret *LookupResult, err error) {
	if ret, err = sw.doLookup(vid, collection); err == nil {
		sw.vc.Set(vid, ret.Locations, 10*time.Minute)
	}
	return
}

func (sw *Seaweed) doLookup(vid, collection string) (*LookupResult, error) {
	values := make(url.Values)
	values.Add("volumeId", vid)
	if collection != "" {
		values.Set("collection", collection)
	}
	jsonBlob, err := sw.HC.Post(sw.Master, "/dir/lookup", values)
	if err != nil {
		return nil, err
	}
	var ret LookupResult
	err = json.Unmarshal(jsonBlob, &ret)
	if err != nil {
		return nil, err
	}
	if ret.Error != "" {
		return nil, errors.New(ret.Error)
	}
	return &ret, nil
}

func (sw *Seaweed) LookupServerByFid(fileId, collection string, readonly bool) (server string, e error) {
	var parts []string
	if strings.Contains(fileId, ",") {
		parts = strings.Split(fileId, ",")
	} else {
		parts = strings.Split(fileId, "/")
	}

	if len(parts) != 2 {
		return "", errors.New("Invalid fileId " + fileId)
	}
	lookup, lookupError := sw.Lookup(parts[0], collection)
	if lookupError != nil {
		return "", lookupError
	}
	if len(lookup.Locations) == 0 {
		return "", errors.New("File Not Found")
	}
	var u string
	if readonly {
		u = lookup.Locations.PickForRead().Url
	} else {
		u = lookup.Locations.Head().Url
	}
	return u, nil
}

func (sw *Seaweed) LookupFileId(fileId, collection string, readonly bool) (fullUrl string, err error) {
	u, e := sw.LookupServerByFid(fileId, collection, readonly)
	if e != nil {
		return "", e
	}
	return MkUrl(u, fileId, nil), nil
}

// LookupVolumeIds find volume locations by cache and actual lookup
func (sw *Seaweed) LookupVolumeIds(vids []string) (map[string]LookupResult, error) {
	ret := make(map[string]LookupResult)
	var unknown_vids []string

	//check vid cache first
	for _, vid := range vids {
		locations, cache_err := sw.vc.Get(vid)
		if cache_err == nil {
			ret[vid] = LookupResult{VolumeId: vid, Locations: locations}
		} else {
			unknown_vids = append(unknown_vids, vid)
		}
	}
	//return success if all volume ids are known
	if len(unknown_vids) == 0 {
		return ret, nil
	}

	//only query unknown_vids
	values := make(url.Values)
	for _, vid := range unknown_vids {
		values.Add("volumeId", vid)
	}
	jsonBlob, err := sw.HC.Post(sw.Master, "/vol/lookup", values)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBlob, &ret)
	if err != nil {
		return nil, errors.New(err.Error() + " " + string(jsonBlob))
	}

	//set newly checked vids to cache
	for _, vid := range unknown_vids {
		locations := ret[vid].Locations
		sw.vc.Set(vid, locations, 10*time.Minute)
	}

	return ret, nil
}
