/*
	原理：使用指针确保可以传值,修改。

	只支持一维数组
	数组需要按下标可打乱放置，如果为则设置为nil


	roadmap：
	[√]支持按顺序的下标放置
*/
package mmap

import (
	"strconv"
	"strings"
)

type Mmap struct {
	globalMap map[string]interface{}
}

func (mmap *Mmap) GetMap() *map[string]interface{} {
	return &mmap.globalMap
}

//对外设置节点（未解析）
func (mmap *Mmap) SetValue(name string, value interface{}) {
	if mmap.globalMap == nil {
		mmap.globalMap = make(map[string]interface{})
	}
	splitName := strings.Split(name, ".")
	nodeLen := len(splitName)
	var nowNode *interface{} = new(interface{})
	*nowNode = &mmap.globalMap
	for i := 0; i < nodeLen; i++ {
		if i == nodeLen-1 {
			if strings.Contains(splitName[i], "[") {
				*nowNode = mmap.getNode(nowNode, splitName[i])
				mmap.setNode(nowNode, splitName[i], value, true)
			} else {
				mmap.setNode(nowNode, splitName[i], value)
			}
		} else {
			*nowNode = mmap.getNode(nowNode, splitName[i])
		}
	}
}

//最终设置节点（解析至最后）
func (mmap *Mmap) setNode(m *interface{}, name string, value interface{}, last ...interface{}) {
	if len(last) > 0 {
		(*(*m).(*interface{})) = value
		return
	}
	//如果最后一个不是数组
	if (*m) == nil {
		sub := make(map[string]interface{})
		(*m) = &sub
	}
	{
		mmap, ok := (*m).(*map[string]interface{})
		if ok {
			(*mmap)[name] = value
			return
		}
	}
	{
		mslice, ok := (*m).(*interface{})
		if ok {
			if (*mslice) == nil {
				*mslice = &map[string]interface{}{name: value}
			} else {
				mmap, ok := (*mslice).(*map[string]interface{})
				if ok {
					(*mmap)[name] = value
				}
			}
		}
	}

}

//获取节点
func (mmap *Mmap) getNode(m *interface{}, nodeName string) interface{} {
	switch {
	case strings.Contains(nodeName, "["):
		idx, err := strconv.Atoi(nodeName[strings.Index(nodeName, "[")+1 : len(nodeName)-1])
		nodeName = nodeName[0:strings.Index(nodeName, "[")]
		if err != nil {
			panic("can not convert idx:" + err.Error())
		}
		node := &map[string]interface{}{}
		ok := false
		node, ok = (*m).(*map[string]interface{})
		if ok {
			if (*node)[nodeName] == nil {
				(*node)[nodeName] = new([]interface{})
			}
		} else {
			face := (*m).(*interface{})
			if (*face) == nil {
				*face = &map[string]interface{}{}
				node = (*face).(*map[string]interface{})
				(*node)[nodeName] = new([]interface{})
			}
		}
		sliceNode := ((*node)[nodeName]).(*[]interface{})
		for {
			if len(*sliceNode) <= idx {
				(*sliceNode) = append((*sliceNode), new(interface{}))
			} else {
				break
			}
		}
		return (*sliceNode)[idx]
	default:
		node, ok := (*m).(*map[string]interface{})
		if ok {
			if strings.Contains(nodeName, "[") {
				if (*node)[nodeName] == nil {
					(*node)[nodeName] = new([]interface{})
				}
				return (*node)[nodeName]
			} else {
				if (*node)[nodeName] == nil {
					sub := make(map[string]interface{})
					(*node)[nodeName] = &sub
				}
				return (*node)[nodeName]
			}
		} else {
			face := (*m).(*interface{})
			if *face == nil {
				*face = &map[string]interface{}{}
			}
			node := (*face).(*map[string]interface{})
			if strings.Contains(nodeName, "[") {
				if (*node)[nodeName] == nil {
					(*node)[nodeName] = new([]interface{})
				}
				return (*node)[nodeName]
			} else {
				if (*node)[nodeName] == nil {
					sub := make(map[string]interface{})
					(*node)[nodeName] = &sub
				}
				return (*node)[nodeName]
			}

		}

	}
	return nil
}
