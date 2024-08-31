<!--------------------------------------
                Basics
--------------------------------------->

<!--different tags-->

`inline code`

`variable`

`sample output`

`keyboard input`

`teletype text`

* * *

<!--code tag inlined in a text-->

When `x = 3`, that means `x + 2 = 5`

A simple equation: `x` = `y` + 2

<!--------------------------------------
                Adjacent
--------------------------------------->

<!--code tags close to each other-->

before `A` middle `B` after

before`A`middle`B`after

<!--code tags next to each other-->

before `A` `B` after

before`AB`after

`ABCDE`

<!--------------------------------------
            Other Tags (outside)
--------------------------------------->

before **`inline code`** after

before *`inline code`* after

before**a`inline code`b**after

before**a`inline code`b**after

before **`inline code`** after

before **`inline code`** after

before *`inline code` and `inline code`* after

before *`inline code` and `inline code`* after

* * *

before **`inline code`** after

before *`inline code`* after

before **`inline code`** after

before *`inline code`* after

* * *

before **`<pre>`** after

<!--------------------------------------
            Other Tags (inside)
--------------------------------------->

<!--code tag with other tags inside-->

before `<img>` after

before after

before `A middle B` after

<!--code tag with pre and other tags inside-->

```

The <img> tag is used to embed an image.

The  tag is used to embed an image.
```

<!--list inside code block-->

```

    
        List Item One
        List Item Two
        List Item Three
    
```

<!--------------------------------------
    Whitespace &amp; Special Characters
--------------------------------------->

<!--empty tags-->

An inline code that is empty except spaces:

beforeafter

before after

before after

before` `after

before ` ` after

before ` ` after

before`  `after

before `  ` after

before `  ` after

beforeafter

before after

before after

```

```

```
 
```

```
  
```

```

  
```

```
Beginning of code
 
  
  


End of code
```

```
Start of many newlines






End of many newlines
```

* * *

<!--white spaces at edges-->

`inline code`

`inline code`

`inline code`

`inline code`

`inline code`

* * *

<!--code tag with backtick-->

An inline code that contains backticks:

``with ` backtick``

```with `` backticks```

`````a ``` b ```` c ` d`````

`` `starting & ending with a backtick` ``

* * *

An inline code that just contains backticks:

before``` `` ```after

before``` `` ```after

before``` `` ```after

before ``` `` ``` after

before ``` `` ``` after

before ``` `` ``` after

before ``` `` ``` after

before ``` `` ``` after

before ``` `` ``` after

* * *

<!--with fence characters-->

````
```
````

```
~~~
```

```````

Some ```
totally `````` normal
` code
```````

```

Some ~~~
totally ~~~~~~ normal
~ code
```

<!--------------------------------------
            Combinations
--------------------------------------->

before `just code` after

before

```
just pre
```

after

before

```
code inside pre
```

after

before

```
pre inside code
```

after

* * *

before `// just code // another line` after

before

```
// just pre
// another line
```

after

before

```
// code inside pre
// another line
```

after

before

```

// pre inside code
// another line
```

after

<!--------------------------------------
            Languages
--------------------------------------->

```one
content
```

```two
content
```

<!--------------------------------------
            With Highlighting
--------------------------------------->

```
Line 0
    Line 1 AB C
    Line 2 AB C
Line 3
```

* * *

```

    Line 1 AB C
    Line 2 AB C
```

* * *

```

    Line 1 AB C
    Line 2 AB C

```