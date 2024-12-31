<!--------------------------------------
                Basics
--------------------------------------->

<!--no href attributes-->

[no href]()

[no href]()

[no href]()

* * *

<!--no content-->

[](/no_content)

[](/no_content)

[](/no_content)

[](/no_content)

[](/no_content)

<!--no content but fallback-->

[link title](/no_content "link title")

* * *

[relative link](/page.html)

[absolute link](http://simple.org/)

[query params](/page?b=1&a=2)

[fragment heading](#heading)

[fragment](#)

Wir freuen uns über eine [Mail](mailto:hi@example.com?body=Hello%0AJohannes)!

<!--link with broken href-->

[broken link](/page)

[broken link](/page%0A%0A.html)

[with whitespace around](example.com)

[with space inside](http://Open%20Demo)

<!--------------------------------------
            Attributes
--------------------------------------->

<!--link with title-->

[content](/ "link title")

[content](/ "  link title  ")

<!--link with multiline title-->

[content](/ " link  title ")

[content](/ '"link title"')

[content](/ "'link title'")

[content](/ '"link title"')

<!--------------------------------------
            Escaping
--------------------------------------->

<!--list with link-->

- [a(b)\[c\]](/page.html)
  
  [a\]](/page.html)

<!--TODO: list with paragraph-->

<!--link-->

[a(b)\[c\]](/page.html)

[a\]](/page.html)

<!--paragraph-->

a(b)\[c]

\[a]

[a

a]

(a)

(a

a)

<!--------------------------------------
            Adjacent
--------------------------------------->

[A](/)[B](/)

[A](/) [B](/)

before[A](/)middle[B](/)after

before [A](/) middle [B](/) after

before [A](/) middle [B](/) after

<!--------------------------------------
        Content and Combinations
--------------------------------------->

<!--link with space-->

before [content](/) after

before [content](/) after

before [content](/) after

* * *

<!--link with inline styles-->

[**bold** and *italic* text](/)

**bold [and *italic*](/) text**

<!--link with one br-->

[A  
B](/) 

<!--link with two br-->

[A  
\
B](/) 

<!--link with three br-->

[A  
\
B](/)

<!--link with multiple div-->

[A
\
B
\
C](/)

<!--multiline link with too many newlines-->

[Start Line
\
End Line](/)

<!--newlines inside link-->

[newlines around the link content](/)

<!--multiline link inside a list item-->

- [first text
  \
  second text](/)

<!--link with image-->

[![](/image.jpg)](/page.html)

<!--multiline link-->

[first text
\
![](/image.jpg)
\
second text](/page.html)

<!--link with headings-->

[**Heading A**  
**Heading B**](/page.html)  

<!--link with an svg-->

[title](/ "title")

<!-- link and strong inside word -->

before [**a inside strong**](/) after

before[**a inside strong**](/)after

before [**strong inside a**](/) after

before[**strong inside a**](/)after

before [**middle**](/) after

before [**middle**](/) after

before[**middle**](/)after

before [**middle**](/) after

before [**middle**](/) after

before [**middle**](/) after

before**[with empty span](/)**after

before **[with empty span](/)** after

before **[with empty span](/)** after

* * *

before**[a](/) b**after

before**[a](/)b**after

before**[a](/) b [c](/)**after

* * *

before[*a inside italic*](/)after

before[*italic inside a*](/)after

before[**a inside b**](/)after

before[**b inside a**](/)after

before[**already bold**](/)after

* * *

before**[middle](/)**after

before**[*inside bold &amp; italic*](/)**after

before***[inside bold &amp; italic](/)a*b**after

before**[inside bold &amp; italic](/)**after

before**a*b[c](/)d*e**after

* * *

before***italic*[link](/)strong**after

<!--------------------------------------
                Nesting
--------------------------------------->

[before
\
another link
\
after](/a)