package aci

import (
	"bytes"
	"fmt"
)

func rnL3Dom(dom string) string {
	return "l3dom-" + dom
}

// ExternalRoutedDomainAdd creates a new L3 External Domain.
func (c *Client) ExternalRoutedDomainAdd(dom, nameAlias string) error {

	me := "ExternalRoutedDomainAdd"

	rn := rnL3Dom(dom)

	api := "/api/node/mo/uni/" + rn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"l3extDomP":{"attributes":{"dn":"uni/%s","name":"%s","nameAlias":"%s","rn":"%s","status":"created"}}}`,
		rn, dom, nameAlias, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ExternalRoutedDomainDel deletes an existing L3 External Domain.
func (c *Client) ExternalRoutedDomainDel(dom string) error {

	me := "ExternalRoutedDomainDel"

	rn := rnL3Dom(dom)

	api := "/api/node/mo/uni.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"polUni":{"attributes":{"dn":"uni","status":"modified"},"children":[{"l3extDomP":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ExternalRoutedDomainList retrieves the list of L3 External Domains.
func (c *Client) ExternalRoutedDomainList() ([]map[string]interface{}, error) {

	me := "ExternalRoutedDomainList"

	key := "l3extDomP"

	api := "/api/node/mo/uni.json?query-target=subtree&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// ExternalRoutedDomainVlanPoolGet retrieves the VLAN pool for the l3 domain.
func (c *Client) ExternalRoutedDomainVlanPoolGet(name string) (string, error) {

	key := "infraRsVlanNs"

	rn := rnL3Dom(name)

	api := "/api/node/mo/uni/" + rn + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("ExternalRoutedDomainVlanPoolGet: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return "", errGet
	}

	c.debugf("ExternalRoutedDomainVlanPoolGet: reply: %s", string(body))

	attrs, errAttr := jsonImdataAttributes(c, body, key, "ExternalRoutedDomainVlanPoolGet")
	if errAttr != nil {
		return "", errAttr
	}

	if len(attrs) < 1 {
		return "", fmt.Errorf("empty list of vlanpool")
	}

	attr := attrs[0]

	pool, found := attr["tDn"]
	if !found {
		return "", fmt.Errorf("vlanpool not found")
	}

	poolName, isStr := pool.(string)
	if !isStr {
		return "", fmt.Errorf("vlanpool is not a string")
	}

	return poolName, nil
}

// Set vlanpool on l3domain
// method: POST
//url: https://apic.labb01.idcn.se/api/node/mo/uni/l3dom-l3_domain001/rsvlanNs.json
//payload{"infraRsVlanNs":{"attributes":{"tDn":"uni/infra/vlanns-[vp002]-static"},"children":[]}}
//response: {"totalCount":"0","imdata":[]}
//
//timestamp: 15:43:31 DEBUG
