package auth

import (
	"fmt"
	"strings"

	"github.com/zhuoqingbin/cchome-admin-topstar/internal/lib"
)

type Tree struct {
	Icon    []string
	Nbsp    string
	PidName string
	Arr     []map[string]interface{}
}

func NewTree() *Tree {
	return &Tree{
		Icon: []string{
			"│", "├", "└",
		},
		Nbsp:    "&nbsp;",
		PidName: "pid",
	}
}

func (t *Tree) Init(arr []map[string]interface{}) {
	t.Arr = arr
}

func (t Tree) GetChild(pid interface{}) (ret []map[string]interface{}) {
	for _, v := range t.Arr {
		if v["pid"] == pid {
			ret = append(ret, v)
		}
	}
	return
}

func (t Tree) GetChildren(pid interface{}, withself bool) (ret []map[string]interface{}) {
	for _, v := range t.Arr {
		if v["pid"] == pid {
			ret = append(ret, v)
			ret = append(ret, t.GetChildren(v["id"], false)...)
		} else if withself && v["id"] == pid {
			ret = append(ret, v)
		}
	}
	return
}

func (t Tree) GetChildrenIDs(pid interface{}, withself bool) (ret []interface{}) {
	for _, v := range t.GetChildren(pid, withself) {
		ret = append(ret, v["id"])
	}
	return
}

func (t Tree) GetParent(pid interface{}) (ret []map[string]interface{}) {
	var _pid interface{}
	for _, v := range t.Arr {
		if v["id"] == pid {
			_pid = v["pid"]
			break
		}
	}
	if _pid != nil {
		for _, v := range t.Arr {
			if v["id"] == _pid {
				ret = append(ret, v)
				break
			}
		}
	}
	return
}

func (t Tree) GetParents(pid interface{}, withself bool) (ret []map[string]interface{}) {
	var _pid interface{}
	for _, v := range t.Arr {
		if v["id"] == pid {
			if withself {
				ret = append(ret, v)
			}
			_pid = v["pid"]
			break
		}
	}
	if _pid != nil {
		ret = append(t.GetParents(_pid, true), ret...)
	}
	return
}

func (t Tree) GetParentsIDs(pid interface{}, withself bool) (ret []interface{}) {
	for _, v := range t.GetParents(pid, withself) {
		ret = append(ret, v["id"])
	}
	return
}

func (t Tree) GetTreeMenu(myID interface{}, itemTpl string, selectedIDs, disabledIDs []interface{}, wrapTag string, wrapAttr string, deepLevel int) (ret string) {
	if wrapTag == "" {
		wrapTag = "ul"
	}
	childs := t.GetChild(myID)
	if len(childs) == 0 {
		return
	}
	for k, v := range childs {
		id := v["id"]
		selected := ""
		if len(selectedIDs) > 0 {
			if ok, _ := lib.Contains(id, selectedIDs); ok {
				selected = "selected"
			}
		}
		disabled := ""
		if len(disabledIDs) > 0 {
			if ok, _ := lib.Contains(id, disabledIDs); ok {
				disabled = "disabled"
			}
		}

		childs[k]["selected"] = selected
		childs[k]["disbaled"] = disabled

		newArr := make(map[string]interface{})
		for _k, _v := range childs[k] {
			newArr[fmt.Sprintf("@%s", _k)] = _v
		}
		childs[k] = newArr

		bakValue := map[string]interface{}{}
		if _url, ok := newArr["@url"]; ok {
			bakValue["@url"] = _url
		}
		if _caret, ok := newArr["@caret"]; ok {
			bakValue["@caret"] = _caret
		}
		if _class, ok := newArr["@class"]; ok {
			bakValue["@class"] = _class
		}

		for _k := range bakValue {
			delete(childs[k], _k)
		}

		nstr := ""
		{
			if length := len(childs[k]); length > 0 {
				oldnew := make([]string, length*2)
				for _o, _n := range childs[k] {
					oldnew = append(oldnew, _o, fmt.Sprintf("%v", _n))
				}
				nstr = strings.NewReplacer(oldnew...).Replace(itemTpl)
			}
		}

		for _k, _v := range bakValue {
			childs[k][_k] = _v
		}

		childData := t.GetTreeMenu(id, itemTpl, selectedIDs, disabledIDs, wrapTag, wrapAttr, deepLevel+1)
		childList := ""
		last := ""
		if len(childData) > 0 {
			childList = fmt.Sprintf(`<%s %s>%s</%s>`, wrapTag, wrapAttr, childData, wrapTag)
			last = "last"
		}
		childList = strings.Replace(childList, `@class`, last, 1)
		childs[k]["@childlist"] = childList

		url, ok := childs[k]["@url"]
		if childData != "" || !ok {
			childs[k]["@url"] = "javascript:;"
			childs[k]["@addtabs"] = ""
		} else {
			childs[k]["@url"] = url
			addtabs := "?"
			if strings.Index(childs[k]["@url"].(string), "?") >= 0 {
				addtabs = "&"
			}
			addtabs = addtabs + "ref=addtabs"
			childs[k]["@addtabs"] = addtabs
		}

		childs[k]["@caret"] = ""
		childs[k]["@badge"] = ""
		{
			childs[k]["@class"] = ""
			if selected != "" {
				childs[k]["@class"] = "active"
			}
			if disabled != "" {
				childs[k]["@class"] = childs[k]["@class"].(string) + " disabled"
			}
			if childData != "" {
				childs[k]["@cliass"] = childs[k]["@class"].(string) + " treeview"
			}
		}

		{
			if length := len(childs[k]); length > 0 {
				oldnew := make([]string, length*2)
				for _o, _n := range childs[k] {
					oldnew = append(oldnew, _o, fmt.Sprintf("%v", _n))
				}
				ret = ret + strings.NewReplacer(oldnew...).Replace(nstr)
			}
		}
	}
	return
}
