## Description
wini is a *.ini* file parser written in golang.  
It will let you read,edit,and create .ini file.

## Type of *.ini* wini supports.  
wini assumes that *.ini* file contains a *section*,*section comment*,*key-value data*,*key-value comments*.  

## Usage
- Reading *.ini* file.   
```golang
file := wini.Load("*iniFilePath.ini*")


```
