package main

import (
	"encoding/xml"
)

type serviceResponse struct {
	XMLName               xml.Name              `xml:"serviceResponse"`
	AuthenticationSuccess authenticationSuccess `xml:"authenticationSuccess"`
}

type authenticationSuccess struct {
	XMLName    xml.Name       `xml:"authenticationSuccess"`
	User       string         `xml:"user"`
	Attributes UserAttributes `xml:"attributes"`
}

//UserAttributes is a struct for defining the attributes that a user who has successfully logged in has.
type UserAttributes struct {
	Cn             string `xml:"cn"`
	Mail           string `xml:"mail"`
	Sn             string `xml:"sn"`
	Ou             string `xml:"ou"`
	ItbStatus      string `xml:"itbStatus"`
	ItbNIM         string `xml:"itbNIM"`
	ItbEmailNonITB string `xml:"itbEmailNonITB"`
}

//ParseResponseXML is a function that returns the attributes of an authenticated user
func ParseResponseXML(b []byte) UserAttributes {
	var resp serviceResponse
	xml.Unmarshal(b, &resp)
	return resp.AuthenticationSuccess.Attributes
}
