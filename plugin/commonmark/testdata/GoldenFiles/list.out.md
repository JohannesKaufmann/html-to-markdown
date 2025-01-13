A paragraph

- 1
- 2
- - 3.1
  - 3.2
- 4 Before
  
  - 4.1
  - 4.2
- - 5.1
  
  5 After
- 6 Before  
  6 also Before
  
  - 6A.1
  
  6 Between
  
  - 6B.1
  
  6 After
  
  6 also After
- 7

* * *

And also other lists...

- First
- Someone once said:
  
  > My famous quote
  
  \- someone

<!--THE END-->

09. Nine
10. Ten
11. 111. Eleven.A
    112. Eleven.B
12. Someone once said:
    
    > My famous quote
    
    \- someone
13. Thirteen

List Item without Container

* * *

<!-- list with blockquote that contains break -->

1. > Line A  
   > Line B

* * *

<!-- parsing the number fails -->

1. one
2. two

* * *

<!-- the max is one character: "9" -->

8. a
9. b

<!--THE END-->

<!-- the max is two characters: "10" -->

09. a
10. b

* * *

- Before text after
- Before [text](/page) after

* * *

- A double `**` [can open strong emphasis](/page)

* * *

- List 1

<!--THE END-->

- List 2

<!--THE END-->

<!--THE END-->

- List 3

<!--THE END-->

- List 4

text between

- List 5

<!--THE END-->

- List 6

<!--THE END-->

- List 7

* * *

- - List 1
  
  <!--THE END-->
  
  - List 2
  
  <!--THE END-->
  
  - List 3

<!--THE END-->

<!-- nesting -->

1. 1. 1. 1. 1. lots of list containers

* * *

1. 1. 1. lots of list items

<!--THE END-->

<!-- with other elements inside the list -->

1. A 1 (div)
   
   A 2 (#text)
2. A 3 (li) A 4 (#text)
   
   1. B 1 (li)
      
      1. C 1 (li)
         
         C 2 (div)
         
         C 3 (div)
      
      B 2 (div)
   2. B 3 (li)

<!--THE END-->

<!-- with breaks -->

- Start Line
  
  End Line
- Start Line
  
  End Line

* * *

<!-- with code block in item -->

- item:
  
  ```
  line 1
  line 2
  ```
- item 2

<!--THE END-->

<!-- with code block in nested item -->

- item 1:
  
  - nested item 1:
    
    ```
    line 1
    line 2
    ```
  - nested item 2
- item 2

* * *

<!--------------------------------------
            Special Characters
--------------------------------------->

1\.

\-

\+

\*

* * *

1\. not a list

\- not a list

\+ not a list

\* not a list