## Description
wini is a *.ini* file parser written in golang.  
It will let you read,edit,and create .ini file.

## Type of *.ini* wini supports.  
- wini assumes that *.ini* file contents are:
  - *section comment*  
  - *section*  
  - *key-value comments* 
  - *key-value data*  

It also assumes that comments come *before* section or key-value data.

- Example:



## Usage
- Reading *.ini* file.   
```golang
file := wini.Load("*iniFilePath.ini*")


```
