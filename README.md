gotoken  
=
一个轻量级token包
-
> * 支持多端登陆
> * 支持单端登陆

# 安装
```
go get github.com/mosalut/gotoken
```

# 使用
可以参看token_test.go文件中注释，并尝试按照《程序结构》种描述，相应地放开注释，修改端号执行，看效果。  

# 程序结构
本程序支持多端令牌，和单端令牌两种模式，static.go 文件中 singleMode 表示当前模式。  
本程序的多端登陆共支持4个端，分别对应常量  
* TOKEN_WEB: 0  
* TOKEN_APP: 1  
* TOKEN_PC: 2  
* TOKEN_OTHERS: 3  
其内部结构是一个[4]\*token数组  
将其设为长度为4的数组而非切片的原因是防止攻击，同一账号不停地开辟token端空间。  
其实 TOKEN_WEB 是不是真正对应web端，可以由用户自己决定。以上其他常量也是。  
如要使用单端登陆也就是：一端登陆，另一端令牌自动失效的效果，则对应常量TOKEN_SINGLE。  
* 创建一个多端登陆的令牌
```
token, err := gotoken.New("user", 20, TOKEN_WEB)
if err != nil {
	errorHandler()
	return
}
```
* 创建一个单端登陆的令牌
```
token, err := gotoken.New("user", 300, TOKEN_SINGLE)
if err != nil {
	log.Println(err.Error())
	return
}
```
* 第一个参数为用户唯一标志  
* 第二个参数为生效时间，单位为秒  
* 第三个参数为令牌模式，如果模式不为 TOKEN_SINGLE，则表示端号  
* 此时在 static.go 中定义的 tokens map 就会根据传入的用户唯一标志（string），添加一个token。tokens结构设为 map[string]interface{} 是因为不确定为 \*token 还是 [4]\*token。所以说首次创建token之后，模式已经确定，并且不可更改。
* 如果你只是使用这个包，可以省略这一段落
> 包变量 singleMode 默认为false，如果是单端令牌模式 gotoken.New 内部会调用包私有方法 newSingle，并将包变量 singleMode 设为true。其实每次创建令牌时，都会配合singleMode和当前数组长度是否为0，来判断是否第一次创建，和当前模式。也就是说只有第一次创建，不关心是模式不匹配，之后的都会关心，但是创建的最后都会根据实际情况设置singleMode。  
> 其实在代码中创建多端令牌部分之判断了singleMode，因为 singleMode 默认为 false，并且只能设置一次。所以没有必要在判断 tokens 长度，如果在创建多端令牌时，singleMode 为 true，则表示一定已经被设过一次。所以 tokens 长度一定大于 0，而即使singleMode 被第一次创建时设为 false，也无所谓，因为对应的模式本来就是多端，所以流程直接走下去。这一段是为了阅读源码方便准备的，源码为了性能的提高和代码量缩短，少了判断。但是比较晦涩。  
包函数 GetCurrentToken 是用来通过用户唯一标识和端号，获取当前令牌，如果端号为 TOKEN_SINGLE，则根据单端模式规则获取，否则根据多端。
```
	token, err = gotoken.GetCurrentToken("user", TOKEN_PC) // 多端令牌模式
	if err != nil {
		log.Println(err.Error())
		return
	}

	if token == nil {
		log.Println("无此令牌")
		return
	}
```
```
	token, err = gotoken.GetCurrentToken("user", TOKEN_SINGLE) // 单端令牌模式
	if err != nil {
		t.Error(err.Error())
		return
	}

	if token == nil {
		log.Println("无此令牌")
		return
	}
```
* 第一个参数是用户唯一标志。  
* 第二个参数是端号。  
* 如果创建时是在WEB端，使用TOKEN_WEB，而获取时是在APP端，也会获取不到——“无此令牌”。  
* 得到当前令牌后则进入验证环节，调用 gotoken.token.Validation 函数  
```
	ok := token.Validation(token.Code)
```
如果 ok == true 则表示验证通过，并且刷新了该令牌的createTimestamp，否则表示失败。  
校验通过后可以调用 gotoken.token.update 函数，刷新令牌的code和Code，也可以不刷新，仍然使用原来的code和Code。  
但是要注意，如果刷新，则要将刷新后的 Code 传给前端，以供其下次调用接口时传来。  
```
token.Update("user")
```
其实此函数中可以增加对这个用户是否存在的判断。但是最终我决定不加，因为可能出现代码传错或输错的时候，正好tokens中有此输错的用户令牌，那也会通过。如：
> 一个 user = "mosalut" 的用户登陆了。  
> 这时一个 user = "Linda" 用户更新时，由于使用者写业务逻辑的时候发生错误，最终让user 正好 = "mosalut"，那么调用 token.Update(user) 的结果将会更新 mosalut 的令牌。  
> 所以请用户自行保证外层代码正确，此处没有错误提示。  


* 另外无需使map线程安全，因为对于map中某个用户下的令牌几乎不会发生同时读取和修改的情况，即使出现了一边修改，一边读取，也没有什么实际的影响，不影响令牌业务。
