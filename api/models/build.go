package models

type Build struct {
	Name   		string 		`json:"name"`
	Code 		string 		`json:"code,omitempty"`
	Deeplearning 	string 		`json:"deeplearning,omitempty"`
	Entrypoint 	string 		`json:"entrypoint,omitempty"`
	Runtime 	string 		`json:"runtime,omitempty"`
	Cmd 		string 		`json:"cmd,omitempty"`

}
