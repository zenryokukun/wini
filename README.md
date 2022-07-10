## Description
wini is a *.ini* file parser written in golang.  
It will let you read,edit,and create .ini file.

## Type of *.ini* wini supports.  
- wini assumes that *.ini* file contents are:
  - *section comment*  
  - *section*  
  - *key-value comments* 
  - *key-value data*  


- How are *comments* treated?
  - wini assumes that comments come *before* section or key-value data.  
  - Texts that start with '#' and ";" are considered as comments by default.  
  - wini **does not** support multi-character comment symbols,multi-line comment symbols, and inline-comment.  


- About *section* and *key-value data*
  - Texts that start with "[" and end with "]" are considered as *section* by default.
  - By default,wini assumes Keys and values are separated by "=" .


- Example:
```
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
```

## Usage
Let's say we have a section file like below:
```
# Name and age of the author.
[Author]
Name = ZEN
Age = 1

# Some basic info about the author.
[Info]
National = JAPAN
Home     = SAKURA-VPS 
Likes    = birds!
# Roaches in Japan are huge.
# I mean it.
Dislikes = roaches!
```

- Reading *.ini* file.   
```golang
// file is a map data that has section names as its keys.
file := wini.Load("iniFilePath.ini")

// Data method will get all the key-value data.
// Key-value data are mapped as map[string]string.
// Note that it will not get the comments.
author := file["Author"].Data()
fmt.Println(author)   // [output]:map[Age:1 Name:ZEN]

// To get section comment, use Com method.0 is the index of the comments.  
// Get method returns the text.
secCom := file["Author"].Com(0).Get()
fmt.Println(secCom)  // [output]: # Name and age of the author.

// Key method returns the key-val data. Then call Com and Get method just like
// getting section comments.
dislikeCom := file["Info"].Key("Dislikes").Com(1).Get() //[output]: # I mean it.
```
It's simple as that.
