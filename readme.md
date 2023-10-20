<h1 align="center">Stone++</h1>
<center>

     ____    __                              __        __      
    /\  _`\ /\ \__                          /\ \      /\ \     
    \ \,\L\_\ \ ,_\   ___     ___      __   \_\ \___  \_\ \___ 
     \/_\__ \\ \ \/  / __`\ /' _ `\  /'__`\/\___  __\/\___  __\
       /\ \L\ \ \ \_/\ \L\ \/\ \/\ \/\  __/\/__/\ \_/\/__/\ \_/
       \ `\____\ \__\ \____/\ \_\ \_\ \____\   \ \_\     \ \_\ 
        \/_____/\/__/\/___/  \/_/\/_/\/____/    \/_/      \/_/ 
                                                           
</center>

# 项目介绍

Stone++是使用Go语言实现的解释器。它是《两周自制脚本语言》中Stone语言的改进版本。

## 与Stone的主要区别

- 使用Go语言实现解释器
- 支持`return`语句
- 支持带参数的构造函数

## BNF

```
elements  : expr { "," expr }
param     : Identifier
params    : param { "," param }
param_list: "(" [ params ] ")"
def       : "def" Identifier param_list block
member    : def | simple
class_body: "{" [ member ] { EOL [ member ] } "}"
defclass  : "class" Identifier [ "extends" Identifier ] class_body
args      : expr { "," expr }
postfix   : "." Identifier | "(" [ args ] ")" | "[" expr "]"
primary   : ("fun" param_list block | "[" [ elements ] "]" | "(" expr ")" | Number | Identifier | String | Boolean) { postfix }
factor    : {"-"} primary
expr      : factor { Op factor}
block     : "{" [ statement ] { EOL [ statement ] } "}"
simple    : expr [ args ] | "return" [ expr ]
statement : "if" expr block { "elif" expr block } [ "else" block ]
          | "while" expr block
          | simple
program   : [ defclass | def | statement ] EOL

```

# 如何使用

```shell
# 构建stone++解释器
go build .

# 直接进入REPL
./stonepp

# 使用 -run 命令运行代码
./stonepp -run a.spp
```

# 入门

## 变量
```
a = 1               // 整型
b = 1.0             // 浮点类型
c = [1,2.0,"3"]     // 数组
d = true            // 布尔类型
e = "hello,world"   // 字符串
```

## 控制流程
### if
```
a = 1
b = 2
if a > b {
    println("a>b")
}elif a == b {
    println("a==b")
}else{
    println("a<b")
}
```
### while
```
i = 0
while i < 10 {
    println(i)
    i = i + 1
}
```

## 数组
```
a = [1]
println(a.len())
a.append(2)
println(a)
a.insert(3,0)
println(a)
a.pop()
println(a)
a.remove(0)
println(a)
println(a.len())

//使用内置函数len
println(len(a))
```

## 函数
### 普通函数
```
def inc(x){
    return x+1
}
println(inc)
println(inc(1))
```

### 匿名函数
```
inc = fun(x){
    return x+1
}

println(inc)
println(inc(1))
```

```
def counter(c){
    return fun(){
        c = c+1
    }
}

c1 = counter(0)
c2 = counter(0)
println(c1())
println(c1())
println(c2())
```

## 类
### 基本类
```
//无构造函数
class Point{
    x = 0
    y = 0
}

p = Point()
println(p)
println(p.x,p.y)

//带构造函数
class Point{
    x = 0
    y = 0
    //和类同名的函数作为构造函数
    //参数名和类成员同名的等于直接对成员赋值
    def Point(x,z){
        y = z
    }
}

p = Point(1,2)
println(p)
println(p.x,p.y)
```

### 继承
```
class Position{
    x = 0
    y = 0
    def Position(x){}
}

//使用extends继承，只支持单继承
class Pos3D extends Position{
    z = 1
    x = 10
    def Pos3D(x){
        //使用super = super(...)调用父类的构造函数
        super = super(x)
    }
}

pp = Position(2)
p = Pos3D(3)
println(p.z)
println(p.x)
//使用super访问父类的成员
println(p.super.x)
p.x = 2
println(p.super.x)
println(p.super.y)
println(p.y)
```

# 示例
### 九九乘法表
```
i = 1
while i < 10{
    j = 1
    while j < 10{
        if i >= j{
            print(i,"*",j,"=",i*j," ")
        }
        j = j+1
    }
    println()
    i = i+1
}
```

### 阶乘
```
def fact(n){
    if n == 1{
        return 1
    }
    return n*fact(n-1)
}

res = fact(5)
println(res)
```

### 斐波那契数列
```
//记忆化搜索
fibList = [1,1]

def fib(n){
    if n < len(fibList) {
        return fibList[n]
    }
    f = fib(n-1)+fib(n-2)
    fibList.append(f)
    return f
}

fib(10)
println(fibList)
```

### 快速幂
```
def pow(a, n){
    if n == 0{
        return 1
    }
    res = pow(a, n/2)
    if n%2 == 0{
        return res*res
    }
    return res*res*a
}

res = pow(2,10)
println(res)
```