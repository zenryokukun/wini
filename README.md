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
  - wini **does not** support multi-character comment symbols,multi-line comment symbols, and inline-comment.  


- About *section* and *key-value data*
  - Texts that start with "[" and end with "]" are considered as *section* by default.
  - wini assumes Keys and values are separated by "=" by default.


- Example:
"""
# This is a "Author" section comment.
# So is.
[Author]
Name = Zenryoku-kun
# This is a "Age" key-val data comment.
Age = 1

# "Info" section comment
[Info]
# "Country" key-val data comment.
Country = Japan

#Comments like below are not supported.
//multi-character comment symbol.
/*
multi-line comment symbol.
*/
Hobby = Fishing # inline comment
"""

## Usage
- Reading *.ini* file.   
```golang
file := wini.Load("*iniFilePath.ini*")


```
