The MJML tool provides a CLI you can use to transform MJML into HTML:

```

bash

$ mjml input.mjml -o output.html
```

Curiously, all text elements (paragraphs and headings) use the same tag,
`<mj-text>`. You can
create headings by applying cosmetic styles as inline attributes, like:

```

html

<mj-text
  align="center"
  font-size="32px"
  font-weight="bold"
  color="#FF0000"
>
```

* * *

```
func HasPrefix(s, prefix []byte) bool
```

* * *

#### HTML

```
<p>Using CSS to change the font color is easy.</p>
<pre>
body {
  color: red;
}
</pre>

```

Copy to Clipboard

Using CSS to change the font color is easy.

```
body {
  color: red;
}

```

* * *

```lang-go
 复制代码const Pi = 3.14159
```

 `const Pi = 3.14159`

```lang-go
 复制代码var n intf(n + 5) // 无类型的数字型常量 “5” 它的类型在这里变成了 int
```

* * *

```
mkdir hugo-contrib
cd hugo-contrib
git clone git@github.com:gohugoio/hugo.git
```

Save the changes and run your first bob command:

```
bob build
```

* * *

To confirm that everything went fine, just run:

```codeBlockLines_39YC
bob --version

```

Copy


* * *

For instance, MetaMask displays your recent transactions by making an API call
to etherscan:

```
GET https://api.etherscan.io/api?module=account&address=0x0208376c899fdaEbA530570c008C4323803AA9E8&offset=40&order=desc&action=txlist&tag=latest&page=1 HTTP/2.0

```

…displays your account balance by making an API call to Infura:

```
POST https://mainnet.infura.io/v3/d039103314584a379e33c21fbe89b6cb HTTP/2.0

{
    "id": 2628746552039525,
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": [
        "0x0208376c899fdaEbA530570c008C4323803AA9E8",
        "latest"
    ]
}

```

* * *

After it was working, I started adding as many different programming languages
as I could. As you can see from this excerpt of my project’s version history,
I got a little overexcited:

```
                                                                   languages
                                                                   ---------
2020-06-05 df9ba38 Initial commit                                          0
2020-06-05 5e3a4a4 Install some packages into a Docker image               0
2020-06-05 e937c8f Simple Express server with "Hello world"                0
2020-06-06 0961498 Embed terminal on frontend app                          0
2020-06-06 c66cf63 Embed Monaco editor on frontend                         0
2020-06-06 27ab1f7 Add "run" button                                        0
2020-06-06 f417858 You can run Python code now                             1
2020-06-07 d543081 You can run many languages now                          8
2020-06-07 e2a3e71 All languages 17 working now                           17
2020-06-07 473c50c ALL THE LANGUAGES                                      25
2020-06-08 3718315 even more languages                                    33
2020-06-08 548c1c1 repl.it superiority!!                                  38
2020-06-08 1ae424f More languages, we need all the languages              48
2020-06-09 c34ccf2 A lot more languages                                   77
2020-06-09 846caf2 At this point the number of languages is absurd        79

```

* * *

Disassembling the now-less-obfuscated dump reveals that one of the addresses
has a label pulled out of somewhere! It’s `strlen`? Going down the
call stack the next one is labeled `vscan_fn` and after that the
labels end, tho I’m fairly confident it’s
[`sscanf`](https://github.com/chakra-core/ChakraCore/blob/master/pal/src/safecrt/sscanf.c#L47).

```
1
2
3
4
5
6
7
8
9
10
11

```

```
...,
{
    "key": "WP_WCT_TINT_21_t2_v9_n2",
    "price": 45000,
    "statName": "CHAR_KIT_FM_PURCHASE20",
    "storageType": "BITFIELD",
    "bitShift": 7,
    "bitSize": 1,
    "category": ["CATEGORY_WEAPON_MOD"]
},
...

```

* * *

```hljs xml
<p>This is the <code>Panel</code> constructor:</p>
<pre><code>function Panel(element, canClose, closeHandler) {
      this.element = element;
      this.canClose = canClose;
      this.closeHandler = function () { if (closeHandler) closeHandler() };
    }</code></pre>
```

Run code snippet
Hide results


Expand snippet

Problem with `<pre>` is it modifies whitespace processing
as well: all spaces are preserved, and wrapping is switched off. Unless
there's a way to switch this off?

**Use the `<textarea>` element to share code**,
like so:

```hljs xml
<textarea class="code" contenteditable="true" spellcheck="false" aria-label='Code Sample'>
  My Sample Bookmark:
  <a href="#bookmark1" id="b1" title="View my bookmark" target="_blank" rel="noreferrer nofollow noopener" accesskey="a" tabindex="0" aria-label="Bookmark">Got to My Bookmark</a>
</textarea>

```

* * *

* * *

```
linux-nfs.vger.kernel.org archive mirror
 help / color / mirror / Atom feed
```

```
From: Leon Romanovsky <leon@kernel.org>
To: "J. Bruce Fields" <bfields@fieldses.org>
Cc: Greg KH <gregkh@linuxfoundation.org>,
	Aditya Pakki <pakki001@umn.edu>,
	Chuck Lever <chuck.lever@oracle.com>,
	Trond Myklebust <trond.myklebust@hammerspace.com>,
	Anna Schumaker <anna.schumaker@netapp.com>,
	"David S. Miller" <davem@davemloft.net>,
	Jakub Kicinski <kuba@kernel.org>,
	Dave Wysochanski <dwysocha@redhat.com>,
	linux-nfs@vger.kernel.org, netdev@vger.kernel.org,
	linux-kernel@vger.kernel.org
Subject: Re: [PATCH] SUNRPC: Add a check for gss_release_msg
Date: Wed, 21 Apr 2021 08:10:25 +0300	[thread overview]
Message-ID: <YH+zwQgBBGUJdiVK@unreal> (raw)
In-Reply-To: <20210420171008.GB4017@fieldses.org>

On Tue, Apr 20, 2021 at 01:10:08PM -0400, J. Bruce Fields wrote:
> On Tue, Apr 20, 2021 at 09:15:23AM +0200, Greg KH wrote:
> > If you look at the code, this is impossible to have happen.
> >
> > Please stop submitting known-invalid patches.  Your professor is playing
> > around with the review process in order to achieve a paper in some
> > strange and bizarre way.
> >
> > This is not ok, it is wasting our time, and we will have to report this,
> > AGAIN, to your university...
>
> What's the story here?

Those commits are part of the following research:
https://github.com/QiushiWu/QiushiWu.github.io/blob/main/papers/OpenSourceInsecurity.pdf

They introduce kernel bugs on purpose. Yesterday, I took a look on 4
accepted patches from Aditya and 3 of them added various severity security
"holes".

Thanks

>
> --b.

```

* * *

```
next prev parent reply	other threads:[~2021-04-21  5:10 UTC|newest]

Thread overview: 49+ messages / expand[flat|nested]  mbox.gz  Atom feed  top
2021-04-07  0:16 Aditya Pakki
2021-04-07 15:34 ` J. Bruce Fields
2021-04-08 15:01 ` Trond Myklebust
2021-04-08 15:24   ` Olga Kornievskaia
2021-04-08 16:02     ` Trond Myklebust
2021-04-20  7:15 ` Greg KH
2021-04-20 17:10   ` J. Bruce Fields
2021-04-21  5:10     ` Leon Romanovsky [this message]
2021-04-21  5:43       ` Greg KH
2021-04-21  6:08         ` Leon Romanovsky
     [not found]         ` <CA+EnHHSw4X+ubOUNYP2zXNpu70G74NN1Sct2Zin6pRgq--TqhA@mail.gmail.com>
2021-04-21  8:15           ` Greg KH
2021-04-21 10:07         ` Sudip Mukherjee
2021-04-21 10:21           ` Greg KH
2021-04-21 11:58             ` Shelat, Abhi
2021-04-21 12:08               ` Greg KH
2021-04-21 12:19               ` Leon Romanovsky
2021-04-21 13:11                 ` Trond Myklebust
2021-04-21 13:20                   ` Leon Romanovsky
2021-04-21 13:42                     ` Steven Rostedt
2021-04-21 13:21                   ` gregkh
2021-04-21 13:34                     ` Leon Romanovsky
2021-04-21 13:50                       ` gregkh
2021-04-21 14:12                         ` Leon Romanovsky
2021-04-21 18:50                         ` Alexander Grund
2021-04-21 13:37               ` J. Bruce Fields
2021-04-21 13:49                 ` Leon Romanovsky
2021-04-21 13:56                   ` J. Bruce Fields
2021-04-22 19:39                     ` J. Bruce Fields
2021-04-23 17:25                       ` Leon Romanovsky
2021-04-23 18:07                         ` J. Bruce Fields
2021-04-23 19:29                           ` Leon Romanovsky
2021-04-23 21:48                             ` J. Bruce Fields
2021-04-24  7:21                               ` Leon Romanovsky
2021-04-24 18:34                               ` Al Viro
2021-04-24 21:34                                 ` J. Bruce Fields
2021-04-25  0:41                                   ` Theodore Ts'o
2021-04-25  6:29                                     ` Greg KH
     [not found]                                       ` <20210426133605.GD21222@fieldses.org>
2021-04-26 13:47                                         ` J. Bruce Fields
2021-04-22  8:10             ` Sudip Mukherjee
2021-04-22  8:27               ` Greg KH
2021-04-21 12:51       ` Anna Schumaker
2021-04-21 14:15         ` Leon Romanovsky
2021-04-21 15:48           ` Theodore Ts'o
2021-04-21 17:34             ` Mike Rapoport
2021-04-22  3:57               ` Leon Romanovsky
2021-04-21 22:52 ` Guenter Roeck
     [not found] <CAHr+ZK-ayy2vku9ovuSB4egtOxrPEKxCdVQN3nFqMK07+K5_8g@mail.gmail.com>
2021-04-21 19:49 ` Theodore Ts'o
2021-04-22  7:50   ` Eric Biggers
2021-04-21 20:27 Weikeng Chen

```

* * *

```
Reply instructions:

You may reply publicly to this message via plain-text email
using any one of the following methods:

* Save the following mbox file, import it into your mail client,
  and reply-to-all from there: mbox

  Avoid top-posting and favor interleaved quoting:
  https://en.wikipedia.org/wiki/Posting_style#Interleaved_style

* Reply using the --to, --cc, and --in-reply-to
  switches of git-send-email(1):

  git send-email \
    --in-reply-to=YH+zwQgBBGUJdiVK@unreal \
    --to=leon@kernel.org \
    --cc=anna.schumaker@netapp.com \
    --cc=bfields@fieldses.org \
    --cc=chuck.lever@oracle.com \
    --cc=davem@davemloft.net \
    --cc=dwysocha@redhat.com \
    --cc=gregkh@linuxfoundation.org \
    --cc=kuba@kernel.org \
    --cc=linux-kernel@vger.kernel.org \
    --cc=linux-nfs@vger.kernel.org \
    --cc=netdev@vger.kernel.org \
    --cc=pakki001@umn.edu \
    --cc=trond.myklebust@hammerspace.com \
    --subject='Re: [PATCH] SUNRPC: Add a check for gss_release_msg' \
    /path/to/YOUR_REPLY

  https://kernel.org/pub/software/scm/git/docs/git-send-email.html

* If your mail client supports setting the In-Reply-To header
  via mailto: links, try the mailto: link

```

* * *

```
This is a public inbox, see mirroring instructions
for how to clone and mirror all data and code used for this inbox;
as well as URLs for NNTP newsgroup(s).
```

* * *

Here's an example of an XML-RPC request:

```xml
POST /RPC2 HTTP/1.0
User-Agent: Frontier/5.1.2 (WinNT)
Host: betty.userland.com
Content-Type: text/xml
Content-length: 181

<?xml version="1.0"?>
<methodCall>
    <methodName>examples.getStateName</methodName>
    <params>
        <param>
            <value><i4>41</i4></value>
            </param>
        </params>
    </methodCall>

```

* * *

For example, running:

```
octosql "SELECT email, COUNT(*) as invoice_count
         FROM invoices.csv JOIN mydb.customers ON invoices.customer_id = customers.id
         WHERE first_name <= 'D'
         GROUP BY email
         ORDER BY invoice_count DESC" --explain 1
```

* * *

If you really want this functionality, you can implement it by using the
following function:

```
function sleepSync(timeout) {
  const sab = new SharedArrayBuffer(1024);
  const int32 = new Int32Array(sab);
  Atomics.wait(int32, 0, 0, timeout);
}
```

* * *

It doesn’t get much simpler than this. You’ll notice we only define a main()
function followed by a println to stdout.

``` rust
fn main() {
    println!("Hello, world!");
}

```

* * *

```ruby
# Ruby knows what you
# mean, even if you
# want to do math on
# an entire Array
cities  = %w[ London
              Oslo
              Paris
              Amsterdam
              Berlin ]
visited = %w[Berlin Oslo]

puts "I still need " +
     "to visit the " +
     "following cities:",
     cities - visited
```

* * *

The equivalent portion is just these two lines:

```
  var btn = new Gtk.Button.with_label("Hello World");
  btn.click.connect(win.close);
```

* * *

To preserve linebreaks inside a `div` using CSS:

```css hljs
div.code {
  white-space: pre;
}

```

* * *

```css
@font-face {
	src: url(http://lea.verou.me/logo.otf);
	font-family: 'LeaVerou';
}
```

* * *

```none
This raw text
is not highlighted
but it still has
line numbers
```

* * *

Live JSX Editor JSX?

```
class HelloMessage extends React.Component {
  render() {
    return <div>Hello {this.props.name}</div>;
  }
}

root.render(<HelloMessage name="Taylor" />);

```

Result

Hello Taylor