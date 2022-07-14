# Description
wini is a *.ini* file parser written in golang.  
It will let you read,edit,and create .ini file.

# Type of *.ini* wini supports.  
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
  - wini assumes Keys and// Swapping key-val data by default.

- Example:
```
# This is a "Author" section comment.
# So is this.
[Author]
Name = Zenryoku-kun
# This is a "Age" key-val data comment.
Age = 1

# "Info" section comment
[Info]
# "Country" key-val data comment.
Country = Japan

# Comments like below are not supported.
//multi-character comment symbol.
/*
multi-line comment symbol.
*/
Hobby = Fishing # inline comment
```

# Usage
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

# 1.Reading *.ini* file: 

Load ini file as *file* struct.
```golang
// file is a map data that has section names as its keys.
file := wini.Load("iniFilePath.ini")

// Data method will get all the key-value data.
// Key-value data are mapped as map[string]string.
// Note that it will not get the comments.
author := file["Author"].Data()
fmt.Println(author)   
```
[output]:
```
map[Age:1 Name:ZEN]
```

Getting section / key-val comment.
```golang
// Get section struct from file.
auth := file["Author"]

// Call Com method of section. 0 is the index of comments.
// Com returns comment struct.
secComs := auth.Com(0)

// Call Get method to get comment text.
com := secComs.Get()

fmt.Println(com) 
```
[output]:
```
# Name and age of the author.
```

```golang
// Key method returns the key-val data. 
// Then call Com on key-val data, and call Get.
dislikeCom := file["Info"].Key("Dislikes").Com(1).Get() 
```
[output]:
```
# I mean it.
```

It's simple as that.  

To Change the default key-val separator,comment symbol, and section symbol, do the following:
```golang
// Make sure to call these before Load.

// Changes key-val separator from "=" to ":".
ChangeSepSym(":")           

// Changes section symbol from "[]" to "''"
ChangeSectionSym("'", "'")  

// Texts starting with "?" will be considered as comments.
// Note that default symbols "#" and ";" are also valid.
AddCommentSym("?")         

file := wini.Load("iniFilePath.ini")
```
# 2.**Editing *.ini* file:**  

## Changing section names:

```golang
// Change section name from Author to Founder.
sec := file["Author"]
file.ChangeName("Author","Founder")

// Print all section data,which has section comments and section name.
// Check function retrieves each line of section comments and section itself as string.
secStr := wini.Check(sec)
fmt.Println(secStr)
```

```
[output]:
# Name and age of the Founder.
[Founder]
Name = ZEN
Age = 1
```

## Changing comment:  

```golang
// Change comment. Note that you need comment symbol.
secCom := sec.Com(0)
secCom.Change("# Name and age of the Founder.")

// When comment is passed to Check, it gets the specified comment as a string,
// not a whole section data.
secComStr := wini.Check(secCom)
fmt.Println(secComStr)
```
[outpupt]:
```
# Name and age of the Founder.
```

## Changing key-val data:
This example shows how to change key, but there also  *ChangeVal* and *ChangeKeyVal* methods.  

```golang
// Change key-val data
kv := file["Info"].Key("Dislikes")
kv.ChangeKey("Hates")

// Check function retrieves each line of key-val comments and key-val itself as string.
kvStr := wini.Check(kv)
fmt.Println(kvStr)
```
[output]:
```
# Roaches in Japan are huge.
# I mean it.
Hates=roaches!
```

## **Removing elements:**
Call *Pop* method on section,key-val data,or comments struct.  


### Removing section:
```golang
// To remove a section, call Pop method of file.
// All section comments and key-val data under "Info" will be removed
file.Pop("Info")
fmt.Println(wini.Check(file))
```
[output]:
```
# Name and age of the author.
[Author]
Name = ZEN
Age = 1
```

### Removing key-val data:

```golang
// To remove key-val data, call Pop method of section.
// It will also remove key-val comments.
info := file["Info"]
info.Pop("Dislikes")
fmt.Println(wini.Check(info))
```
[output]:
```
# Some basic info about the author.
[Info]
National = JAPAN
Home     = SAKURA-VPS
Likes    = birds!
```

### Removing Comment

```golang
dislikes := file["Info"].Key("Dislikes")
// Specify the comment to remove by index.
dislikes.PopCom(0)
fmt.Println(wini.Check(dislikes))
```
[output]:
```
# I mean it.
Dislikes = roaches!
```
You can also call *PopAllCom()* method of *section* or *keyval* to remove all comments.  
  

## **Creating new section and key-val data,and add to file.**
### Create new section.
```golang
// Create new section, "Employee".
sec := wini.NewSection("Employee")
// Add Comments. Note that you need a comment symbol.
sec.AddCom("# First section comment","# Second section comment")
```

### Create new key-val data
```golang
// Create two key-val data.
kv := wini.NewKeyVal("ZEN","RYOKU")
kv2 := wini.NewKeyVal("John","Adams")

// Add comment to the second data.
kv2.AddCom("# Key-val comment.")
```

### Add them to file.
```golang
// First, add key-val data to section. Then add to file.
sec.AddKeyVal(kv,kv2)
file.AddSec(sec)
fmt.Println(wini.Check(file))
```
[output]:
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

# First section comment
# Second section comment
[Employee]
ZEN=RYOKU
# Key-val comment.
John=Adams
```

## Helper function *Check*
If you want to take a look at your file,section,or keyval, use *Check* function.  
*Check* returns the data as string. Keep in mind that it ignores empty lines.  
It will also add 1 new empty line between sections for visibility when *file* is passed.

```golang
str := wini.Check(file["Info"])
fmt.Println(str)
```
[output]:
```
[Info]
National = JAPAN
Home     = SAKURA-VPS
Likes    = birds!
# Roaches in Japan are huge.
# I mean it.
Dislikes = roaches!
```


## Save changes.
There are two methods to save, *Save()* and *Savef(string,int,int,int)*.  
Both methods will save your file and creates backupfile.  
Backupfile is named *winiBK_filename.ini*, and will only be created when there is no backupfile.  
If you specified a new file name, obviously backupfile will not be created.

```golang
// Save method will simply save file struct as it is.
// If you have added a section or key-val data, empty lines
// between them might be inconsistent.
// If you want them to be consistent,use Savef method.

//Load example file, remove all comments, then save.
file := wini.Load(iniFilePath.ini)
file.PopAllCom()
file.Save(iniFilePath.ini)
```
Your *iniFilePath.ini* would now look like this.
```
[Author]
Name = ZEN
Age = 1

[Info]
National = JAPAN
Home     = SAKURA-VPS 
Likes    = birds!
Dislikes = roaches!
```
If it was your first time saving using wini, and saved as same name,  
backupfile will be created.

```
-Dir
--iniFilePath.ini
--winiBK_iniFilePath.ini
```

*Savef* method will remove all existing empty lines, and let you sepcify them. You can also specify indents on key-val data.

```golang
file := wini.Load(iniFilePath.ini)

// para/m1: name to save as.
// param2: number of new lines between sections.
// param3: number of new lines between key-vals.
// param4: number of spaces before key-val data.

file.Savef("newfile.ini",1,0,2)

```
Now your *newfile.ini* would look like this. Notice there is a single empty line between sections, and 2 spaces before key-val data.
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
  
## **Change order of  sections or key-val**

Swapping sections.

```golang
file.Swap("Author","Info")
fmt.Println(wini.Check(file))
```
[output]:
```
# Some basic info about the author.
[Info]
National = JAPAN
Home     = SAKURA-VPS
Likes    = birds!
# Roaches in Japan are huge.
# I mean it.
Dislikes = roaches!

# Name and age of the author.
[Author]
Name = ZEN
Age = 1
```
Swapping key-val data.
```golang
file["Info"].Swap("Likes", "Dislikes")
fmt.Println(wini.Check(file["Info"]))
```
[output]:  
```
# Some basic info about the author.
[Info]
National = JAPAN
Home     = SAKURA-VPS
# Roaches in Japan are huge.
# I mean it.
Dislikes = roaches!
Likes    = birds!
```
You can see that Employee section is added to file.

# 3.Creating new ini file from scratch.
Just use json or something.

```golang
// Here are the steps:
// 1.Create empty file.
// 2.Create section and key-val data.
// 4.Add key-val data to section.
// 5.Add section to file
// 6.Save

file := wini.NewFile()
sec1 := wini.NewSection("Company")
kv1 := wini.NewKeyVal("Name","Unemployed Inc.")
kv2 := wini.NewKeyVal("Code","010")
sec2 := wini.NewSection("Requirements")
kv3 := wini.NewKeyVal("Experience","Software development.")
kv4 := wini.NewKeyVal("Age","-")
kv5 := wini.NewKeyVal("Nationality","Bot")

//Add key-val data to sections.
sec1.AddKeyVal(kv1,kv2)
sec2.AddKeyVal(kv3,kv4,kv5)

//Add to file.
file.AddSec(sec1,sec2)

//Finally save.
file.Savef("scratch.ini",1,0,0)
```
There we go.
```
[Company]
Name=Unemployed Inc.
Code=010

[Requirements]
Experience=Software development.
Age=-
Nationality=Bot
```

# Report Bugs!
https://twitter.com/zenryoku_kun0