## Description
wini is a *.ini* file parser written in golang.  
It will let you read,edit,and create .ini file.

## Type of *.ini* wini supports.  
- wini assumes that *.ini* file contents are:
  - *section comment*  
  - *section*  
  - *key-value comments* 
  - *key-value data*  

- How are *omments* treated?
  - wini assumes that comments come *before* section or key-value data.  
  - Texts that start with '#' and ";" are considered as comments by default.  
  - wini **does not** support multi-character comment,multi-line comment, and inline-comment.  

- About *section* and *key-value data*
  - Texts that start with "[" and end with "]" are considered as *section*.
- Example:



## Usage
- Reading *.ini* file.   
```golang
file := wini.Load("*iniFilePath.ini*")


```
