#mmap
- 无需初始化，内部自动进行初始化
- mmap内部使用指针确保引用传值
- 只支持一维数组
- 数组无需按下标顺序进行放置

``` go
	m := new(mmap.Mmap)
	m.SetValue("a", "a")
	m.SetValue("b.c", "c")
	m.SetValue("c[0]", "c0")
	m.SetValue("d.e.f[0].g[1].c", "ccc")
	m.SetValue("e.b[0]", "bbb")
	j, _ := json.Marshal(m.GetMap())
	fmt.Println(string(j))
```

``` json
{"a":"a","b":{"c":"c"},"c":["c0"],"d":{"e":{"f":[{"g":[null,{"c":"ccc"}]}]}},"e":{"b":["bbb"]}}
```