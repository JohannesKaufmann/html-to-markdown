<!--------------------------------------
                Basics
--------------------------------------->

<!--no src attributes-->

<!--basics-->

![](/relative_url)

![](www.example.com/absolute_url)

<!--other attributes-->

![alt text](/url)

![](/url "title text")

![alt text](/url "title text")

<!--------------------------------------
            Special Characters
--------------------------------------->

![  the  alt  attribute  ](/url)

![the alt "attribute"](/url)

![the alt 'attribute'](/url)

![the alt attribute](/url)

![the \[alt\] attribute](/url)

![the (alt) attribute](/url)

![the \](alt) attribute](/url)

* * *

![](/url "  the  title  attribute  ")

![](/url 'the title "attribute"')

![](/url "the title 'attribute'")

![](/url "the title attribute")

![](/url "the [title] attribute")

![](/url "the (title) attribute")

![](/url "the )(title) attribute")

<!--------------------------------------
                Weird URLs
--------------------------------------->

<!--image with data uri-->

![](data:image/gif;base64,abcdefghij)

![](data:image/svg+xml;utf8,%3Csvg%20xmlns='http://www.w3.org/2000/svg'%20width='1080'%20height='956'%3E%3C/svg%3E)

<!--------------------------------------
            Combinations
--------------------------------------->

<!-- link with just image -->

[*![Such Icon](/search.svg)*]() [*![Email Icon](/email.svg)*]()

[*![Such Icon](/search.svg)*]() [*![Email Icon](/email.svg)*]()

* * *

<!--image inside a link-->

[![image alt text](/image.jpg "image title text")](/page.html "link title text")  

<!--image inside an empty link-->

[![image alt text](/src)]()

<!--------------------------------------
            Picture / Figure
--------------------------------------->

![alt text](/image.jpg "title text")

* * *

![alt text](/image.jpg "title text")

caption text