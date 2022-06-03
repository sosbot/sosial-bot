package main

import (
	"fmt"
	_ "github.com/lib/pq"
)

type ViewData struct {
	Title   string
	Message string
}

var tplInputTemplate = `<div class="col-xs-12" ><label for="%s">%s</label> <input type="text" class="form-control form-control-lg" id="%s" name="%s" required size="%s" placeholder="%s" minlength="%s" maxlength="%s" title="%s"></div><br>`
var tplDateTemplate = `<div class="col-xs-12"><label for="%s">%s</label><input type="date"  class="form-control form-control-lg" id="%s" name="%s" value="" min="1900-01-01" max="2030-12-31"></div><br>`
var tplSelectTemplate = `<div class="col-xs-12"><label for="%s">%s</label><select class="form-control-lg"  name="%s" id="%s">%s</select></div><br>`
var tplCheckboxTemplate = `<div class="col-xs-12"><input  class="form-control-lg""  type="checkbox" id="%s" name="%s" value=%s><label for="%s">%s</label></div><br>`

type inputForm struct {
	Fields []string
}

type inputField struct {
	template    string
	Id          string
	Name        string
	Label       string
	Placeholder string
	ReqSize     string
	MinLength   string
	MaxLength   string
	Pattern     string
	ErrMsg      string
}

type selectField struct {
	Template string
	Label    string
	Name     string
	Id       string
	Values   string
}

type checkboxField struct {
	Template string
	Id       string
	Name     string
	Value    string
	Label    string
}

type Field struct {
	Value string
}

type RepoData struct {
	service_name           string
	order_num              string
	component_description  string
	component_type         string
	data_driven            string
	componentId            string
	component_id           string
	component_name         string
	component_value        string
	component_label        string
	component_requiredsize string
	component_placeholder  string
	component_minlength    string
	component_maxlength    string
	component_title        string
	component_mindate      string
	component_maxdate      string
	data_value             string
}

type componentData struct {
	ComponentId    string
	ComponentValue string
}
type componentDatas struct {
	ComponentDatas []componentData
}

func (m *inputField) appendText() string {
	res := fmt.Sprintf(m.template, m.Id, m.Label, m.Id, m.Name, m.ReqSize, m.Placeholder, m.MinLength, m.MaxLength, m.ErrMsg)
	return res
}

func (m *inputField) appendDate() string {
	res := fmt.Sprintf(m.template, m.Id, m.Label, m.Id, m.Name)
	return res
}

func (m *selectField) appendSelect() string {
	res := fmt.Sprintf(m.Template, m.Id, m.Label, m.Name, m.Id, m.Values)
	return res
}

func (m *checkboxField) appendCheckbox() string {
	res := fmt.Sprintf(m.Template, m.Id, m.Name, m.Value, m.Id, m.Label)
	return res
}

func (m *inputForm) fieldsToString() string {
	var res string
	for _, v := range m.Fields {
		res = res + v
	}
	return res
}
