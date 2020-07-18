package md

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var update = flag.Bool("update", false, "update .golden files")

func TestFromString(t *testing.T) {
	var tests = []struct {
		name string

		domain  string
		html    string
		options *Options
	}{
		{
			name: "p tag",
			html: `<p>Some Text</p>`,
		},
		{
			name: "two p tags",
			html: `
			<div>
				<p>Text</p>
				<p>Some Text</p>
			</div>
			`,
		},
		{
			name: "span in p tag",
			html: "<p>Some <span>Text</span></p>",
		},
		{
			name: "strong in p tag",
			html: "<p>Some <strong>Text</strong></p>",
		},
		{
			name: "strong in p tag with whitespace",
			html: "<p> Some <strong> Text </strong></p>",
		},
		{
			name: "strong in p tag with whitespace inside",
			html: "<p>Some<strong> Text. </strong></p>",
		},
		{
			name: "strong in p tag whithout whitespace",
			html: "<p>Some<strong>Text</strong></p>",
		},
		{
			name: "strong in p tag with spans",
			html: "<p><span>Some</span><strong>Text</strong>Content</p>",
		},
		{
			name: "strong in p tag with punctuation",
			html: "<p><span>Some</span><strong>Text.</strong></p>",
		},
		{
			name: "i inside an em ",
			html: `<em><i>Double</i>Italic</em>`,
		},
		{
			name: "em in a p",
			html: `<p>首付<em class="red"><i>19,8</i>万</em> 月供</p>`,
		},
		{
			name: "two em in a h4",
			html: `<h4>首付<em class="red"><i>19,8</i>万</em> /
			月供<em class="red">6339元X24</em></h4>`,
		},
		{
			name: "em in p tag",
			html: "<p>Some <em>Text</em></p>",
		},
		{
			name: "em in p tag with whitespace",
			html: "<p> Some <em> Text </em></p>",
		},
		{
			name: "h1",
			html: "<h1>Header</h1>",
		},
		{
			name: "h2",
			html: "<h2>Header</h2>",
		},
		{
			name: "h6",
			html: "<h6>Header</h6>",
		},
		{
			name: "escape h1",
			html: "<h1>#hashtag</h1>",
		},
		{
			name: "heading in link",
			html: `
<div class="p-1 w-full sm:w-1/2 md:w-1/3 lg:w-1/4">
    <a class="no-underline text-black" href="https://code.design/events/zuhause_jul20">
        <div class="flex flex-col justify-between h-full rounded-lg shadow hover:shadow-lg bg-white">
            <div class="flex flex-row px-4 py-3 rounded-t-lg shadow bg-white">
                <h3 class="mr-auto text-lg font-lemonmilklight font-semibold">#zuhause_jul20</h3>
                <span class="text-sm rounded-full shadow-lg bg-primary-200 text-primary-800 font-semibold px-3" style="padding-top:.1rem;padding-bottom:.1rem;">
    Remote
</span>
            </div>
            <div class="px-4 py-3">
                <div class="-mx-2">
                    <div class="flex flex-row items-center p-2 w-full">
    <div class="mr-2"><svg class="fill-current w-auto text-primary-600" xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 27 27">
    <path d="M24.1877 9.33745H3.2627C2.9252 9.33745 2.7002 9.11245 2.7002 8.88745C2.7002 8.54995 2.9252 8.32495 3.2627 8.32495H24.1877C24.5252 8.32495 24.7502 8.54995 24.7502 8.88745C24.7502 9.11245 24.5252 9.33745 24.1877 9.33745Z"></path>
    <path d="M9.675 6.18755H7.3125C6.975 6.18755 6.75 5.96255 6.75 5.73755V1.80005C6.75 1.46255 6.975 1.23755 7.3125 1.23755H9.675C9.9 1.23755 10.2375 1.46255 10.2375 1.80005V5.62505C10.2375 5.96255 10.0125 6.18755 9.675 6.18755ZM7.7625 5.17505H9.1125V2.36255H7.7625V5.17505Z"></path>
    <path d="M19.9123 6.18755H17.5498C17.2123 6.18755 16.9873 5.96255 16.9873 5.62505V1.80005C16.9873 1.46255 17.2123 1.23755 17.5498 1.23755H19.9123C20.1373 1.23755 20.3623 1.46255 20.3623 1.80005V5.62505C20.3623 5.96255 20.1373 6.18755 19.9123 6.18755ZM17.9998 5.17505H19.3498V2.36255H17.9998V5.17505Z"></path>
    <path d="M16.9876 25.5375H5.2876C3.7126 25.5375 2.4751 24.3 2.4751 22.725V6.07495C2.4751 4.49995 3.7126 3.26245 5.2876 3.26245H7.3126C7.5376 3.26245 7.8751 3.48745 7.8751 3.82495V5.28745H9.2251V3.82495C9.2251 3.48745 9.4501 3.26245 9.7876 3.26245H17.5501C17.7751 3.26245 18.1126 3.48745 18.1126 3.82495V5.28745H19.4626V3.82495C19.4626 3.48745 19.6876 3.26245 20.0251 3.26245H22.0501C23.5126 3.26245 24.7501 4.49995 24.7501 6.07495V17.775C24.7501 17.8875 24.6376 18 24.6376 18.1125L17.4376 25.3125C17.2126 25.5375 17.1001 25.5375 16.9876 25.5375ZM5.2876 4.27495C4.2751 4.27495 3.4876 5.06245 3.4876 6.07495V22.725C3.4876 23.7375 4.3876 24.525 5.2876 24.525H16.7626L23.7376 17.55V6.07495C23.7376 5.06245 22.9501 4.27495 21.9376 4.27495H20.4751V5.62495C20.4751 5.96245 20.2501 6.18745 20.0251 6.18745H17.6626C17.3251 6.18745 17.1001 5.96245 17.1001 5.62495V4.27495H10.3501V5.62495C10.3501 5.96245 10.1251 6.18745 9.7876 6.18745H7.3126C6.9751 6.18745 6.7501 5.96245 6.7501 5.73745V4.27495H5.2876Z"></path>
    <path d="M16.9878 25.5375C16.8753 25.5375 16.8753 25.5375 16.7628 25.5375C16.6503 25.425 16.4253 25.3125 16.4253 25.0875V20.1375C16.4253 18.5625 17.6628 17.325 19.2378 17.325H24.1878C24.4128 17.325 24.5253 17.4375 24.6378 17.6625C24.7503 17.8875 24.6378 18.1125 24.5253 18.225L17.3253 25.425C17.2128 25.5375 17.1003 25.5375 16.9878 25.5375ZM19.3503 18.3375C18.3378 18.3375 17.5503 19.125 17.5503 20.1375V23.85L23.0628 18.3375H19.3503Z"></path>
    <path d="M13.9501 12.0376H13.1626C12.9376 12.0376 12.6001 11.8126 12.6001 11.4751C12.6001 11.1376 12.8251 10.9126 13.1626 10.9126H13.9501C14.2876 10.9126 14.4001 11.1376 14.4001 11.4751C14.4001 11.8126 14.2876 12.0376 13.9501 12.0376Z"></path>
    <path d="M17.5502 12.0376H16.7627C16.5377 12.0376 16.2002 11.8126 16.2002 11.4751C16.2002 11.1376 16.4252 10.9126 16.7627 10.9126H17.5502C17.7752 10.9126 18.1127 11.1376 18.1127 11.4751C18.1127 11.8126 17.7752 12.0376 17.5502 12.0376Z"></path>
    <path d="M21.1503 12.0376H20.3628C20.0253 12.0376 19.8003 11.8126 19.8003 11.4751C19.8003 11.1376 20.0253 10.9126 20.3628 10.9126H21.1503C21.3753 10.9126 21.7128 11.1376 21.7128 11.4751C21.7128 11.8126 21.3753 12.0376 21.1503 12.0376Z"></path>
    <path d="M13.9501 14.625H13.1626C12.9376 14.625 12.6001 14.4 12.6001 14.175C12.6001 13.8375 12.8251 13.6125 13.1626 13.6125H13.9501C14.2876 13.6125 14.4001 13.8375 14.4001 14.175C14.5126 14.4 14.2876 14.625 13.9501 14.625Z"></path>
    <path d="M17.5502 14.625H16.7627C16.5377 14.625 16.2002 14.4 16.2002 14.175C16.2002 13.8375 16.4252 13.6125 16.7627 13.6125H17.5502C17.7752 13.6125 18.1127 13.8375 18.1127 14.175C18.1127 14.4 17.7752 14.625 17.5502 14.625Z"></path>
    <path d="M21.1503 14.625H20.3628C20.0253 14.625 19.8003 14.4 19.8003 14.175C19.8003 13.8375 20.0253 13.6125 20.3628 13.6125H21.1503C21.3753 13.6125 21.7128 13.8375 21.7128 14.175C21.6003 14.4 21.3753 14.625 21.1503 14.625Z"></path>
    <path d="M13.9501 17.2126H13.1626C12.9376 17.2126 12.6001 16.9876 12.6001 16.6501C12.6001 16.3126 12.8251 16.0876 13.1626 16.0876H13.9501C14.2876 16.0876 14.4001 16.3126 14.4001 16.6501C14.4001 16.9876 14.2876 17.2126 13.9501 17.2126Z"></path>
    <path d="M6.8627 19.9126H6.0752C5.7377 19.9126 5.5127 19.6876 5.5127 19.3501C5.5127 19.0126 5.7377 18.7876 6.0752 18.7876H6.8627C7.0877 18.7876 7.4252 19.0126 7.4252 19.3501C7.3127 19.6876 7.0877 19.9126 6.8627 19.9126Z"></path>
    <path d="M10.3503 19.9126H9.67529C9.33779 19.9126 9.11279 19.6876 9.11279 19.3501C9.11279 19.0126 9.33779 18.7876 9.67529 18.7876H10.4628C10.8003 18.7876 11.0253 19.0126 11.0253 19.3501C10.9128 19.6876 10.6878 19.9126 10.3503 19.9126Z"></path>
    <path d="M13.9501 19.9126H13.1626C12.9376 19.9126 12.6001 19.6876 12.6001 19.3501C12.6001 19.0126 12.8251 18.7876 13.1626 18.7876H13.9501C14.2876 18.7876 14.4001 19.0126 14.4001 19.3501C14.5126 19.6876 14.2876 19.9126 13.9501 19.9126Z"></path>
    <path d="M6.8625 22.5H6.075C5.7375 22.5 5.625 22.275 5.625 21.9375C5.625 21.7125 5.85 21.4875 6.1875 21.4875H6.975C7.2 21.4875 7.5375 21.7125 7.5375 21.9375C7.3125 22.275 7.0875 22.5 6.8625 22.5Z"></path>
    <path d="M10.3498 22.5H9.6748C9.3373 22.5 9.1123 22.275 9.1123 21.9375C9.1123 21.7125 9.3373 21.4875 9.6748 21.4875H10.4623C10.7998 21.4875 11.0248 21.7125 11.0248 21.9375C10.9123 22.275 10.6873 22.5 10.3498 22.5Z"></path>
    <path d="M13.9501 22.5H13.1626C12.9376 22.5 12.6001 22.275 12.6001 21.9375C12.6001 21.7125 12.8251 21.4875 13.1626 21.4875H13.9501C14.2876 21.4875 14.4001 21.7125 14.4001 21.9375C14.5126 22.275 14.2876 22.5 13.9501 22.5Z"></path>
    <path d="M9.1125 17.5501H7.0875C6.75 17.5501 6.525 17.3251 6.525 16.9876V15.5251H5.0625C4.725 15.5251 4.5 15.3001 4.5 15.0751V13.0501C4.5 12.8251 4.725 12.6001 5.0625 12.6001H6.525V11.0251C6.525 10.6876 6.75 10.4626 7.0875 10.4626H9.1125C9.3375 10.4626 9.675 10.6876 9.675 11.0251V12.4876H11.1375C11.3625 12.4876 11.7 12.7126 11.7 12.9376V14.9626C11.7 15.1876 11.475 15.4126 11.1375 15.4126H9.5625V16.8751C9.5625 17.3251 9.3375 17.5501 9.1125 17.5501ZM7.5375 16.5376H8.55V15.0751C8.55 14.8501 8.775 14.6251 9.1125 14.6251H10.575V13.5001H9.1125C8.8875 13.5001 8.55 13.2751 8.55 12.9376V11.4751H7.5375V12.9376C7.5375 13.2751 7.3125 13.5001 7.0875 13.5001H5.5125V14.6251H6.975C7.2 14.6251 7.425 14.8501 7.425 15.0751V16.5376H7.5375Z"></path>
</svg></div>
    <div class="leading-tight">
        <span class="text-gray-600 font-medium">Datum</span><br>
        31.07. - 02.08.20
    </div>
</div>
                    <div class="flex flex-row items-center p-2 w-full">
    <div class="mr-2"><svg class="fill-current w-auto text-primary-600" xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 35 35">
    <path d="M17.5002 34.4167C8.16683 34.4167 0.583496 26.8333 0.583496 17.5C0.583496 8.16668 8.16683 0.583344 17.5002 0.583344C26.8335 0.583344 34.4168 8.16668 34.4168 17.5C34.4168 26.8333 26.8335 34.4167 17.5002 34.4167ZM17.5002 2.04168C8.896 2.04168 2.04183 8.89584 2.04183 17.5C2.04183 26.1042 9.04183 32.9583 17.5002 32.9583C26.1043 32.9583 32.9585 25.9583 32.9585 17.5C32.9585 8.89584 26.1043 2.04168 17.5002 2.04168Z"></path>
    <path d="M17.5003 26.5417C12.542 26.5417 8.60449 22.4584 8.60449 17.6459C8.60449 12.8334 12.542 8.60419 17.5003 8.60419C22.4587 8.60419 26.3962 12.6875 26.3962 17.5C26.3962 22.3125 22.4587 26.5417 17.5003 26.5417ZM17.5003 10.0625C13.417 10.0625 10.0628 13.4167 10.0628 17.5C10.0628 21.5834 13.417 24.9375 17.5003 24.9375C21.5837 24.9375 24.9378 21.5834 24.9378 17.5C24.9378 13.4167 21.5837 10.0625 17.5003 10.0625Z"></path>
    <path d="M17.5002 18.2292C17.0627 18.2292 16.771 17.9375 16.771 17.5V12.25C16.771 11.8125 17.0627 11.5208 17.5002 11.5208C17.9377 11.5208 18.2293 11.8125 18.2293 12.25V17.5C18.2293 17.9375 17.9377 18.2292 17.5002 18.2292Z"></path>
    <path d="M15.4588 20.2709C15.313 20.2709 15.0213 20.2709 14.8755 20.125C14.5838 19.8334 14.5838 19.3959 14.8755 19.1042L16.9172 17.0625C17.2088 16.7709 17.6463 16.7709 17.938 17.0625C18.2297 17.3542 18.2297 17.7917 17.938 18.0834L15.8963 20.125C15.7505 20.2709 15.6047 20.2709 15.4588 20.2709Z"></path>
</svg></div>
    <div class="leading-tight">
        <span class="text-gray-600 font-medium">Uhrzeit</span><br>
        Fr 15 Uhr - So 19 Uhr
    </div>
</div>
                </div>
            </div>
                            <div class="px-4 py-3 border-t border-gray-300 font-semibold text-primary-500">
                    Details anzeigen →
                </div>
                    </div>
    </a>
</div>
			`,
		},
		{
			name: "setext h1",
			html: "<h1>Header</h1>",
			options: &Options{
				HeadingStyle: "setext",
			},
		},
		{
			name: "setext h2",
			html: "<h2>Header</h2>",
			options: &Options{
				HeadingStyle: "setext",
			},
		},
		{
			name: "setext h3",
			html: "<h3>Header</h3>",
			options: &Options{
				HeadingStyle: "setext",
			},
		},
		{
			name: "ul",
			html: `
			<ul>
				<li>Some Thing</li>
				<li>Another Thing</li>
			</ul>
			`,
		},
		{
			name: "ol",
			html: `
			<ol>
				<li>First Thing</li>
				<li>Second Thing</li>
			</ol>
			`,
		},
		{
			name: "indent content in li",
			html: `
			<ul>
				<li>
					Indent First Thing
					<p>Second Thing</p>
				</li>
				<li>Third Thing</li>
			</ul>
			`,
		},
		{
			name: "nested list",
			html: `
			<ul>
				<li>foo
					<ul>
						<li>bar
							<ul>
								<li>baz
									<ul>
										<li>boo</li>
									</ul>
								</li>
							</ul>
						</li>
					</ul>
				</li>
			</ul>
			`,
		},

		{
			name: "nested list real world",
			html: `
<ul class="primary">
	<li><a href="/" title="Startseite">Startseite</a></li>
	<li>
		<a href="/die-gruppe/unsere-unternehmen/" title="Die Gruppe">Die Gruppe</a>
		<ul>
			<li><a href="/die-gruppe/unsere-unternehmen/" title="Unsere Unternehmen">Unsere Unternehmen</a></li>
			<li><a href="/die-gruppe/unternehmenshistorie/" title="Unternehmenshistorie">Unternehmenshistorie</a></li>
			<li><a href="/die-gruppe/standortportraits/" title="Standortportraits">Standortportraits</a></li>
			<li><a href="/die-gruppe/unsere-marken/" title="Unsere Marken">Unsere Marken</a></li>
			<li><a href="/die-gruppe/kontakt/" title="Kontakt">Kontakt</a></li>
		</ul>
	</li>
	<li class="active">
		<a href="/medien/aktuelle-meldungen/" title="Medien">Medien</a>
		<ul>
			<li><a href="/medien/aktuelle-meldungen/" title="Aktuelle Meldungen">Aktuelle Meldungen</a></li>
			<li><a href="/medien/pressearchiv/" title="Pressearchiv">Pressearchiv</a></li>
			<li><a href="/medien/pressekontakt/" title="Pressekontakt">Pressekontakt</a></li>
			<li class="active"><a href="/medien/einblicke/" title="Einblicke">Einblicke</a></li>
		</ul>
	</li>
	<li>
		<a href="/karriere/" title="Karriere">Karriere</a>
		<ul>
			<li><a href="/karriere/video-einblicke/" title="Video-Einblicke">Video-Einblicke</a></li>
			<li><a href="https://career5.successfactors.eu/career?company=mllerservi&amp;site=VjItNGdGZlNGSEJEYTVJSVRUaXp4N1E4Zz09" target="_blank" title="Stellenangebote">Stellenangebote</a></li>
			<li><a href="/karriere/erlebnisberichte/" title="Erlebnisberichte">Erlebnisberichte</a></li>
			<li><a href="/karriere/traineeprogramm/" title="Traineeprogramm">Traineeprogramm</a></li>
			<li><a href="/karriere/termine/" title="Termine">Termine</a></li>
			<li><a href="/karriere/news/" title="News">News</a></li>
		</ul>
	</li>
</ul>
		 
				 `,
		},
		{
			name: "italic with no space after",
			html: `
<p><em>Content </em>and no space afterward.</p>	
			`,
		},
		{
			name: "ul in ol",
			html: `
			<ol>
				<li>
					<p>First Thing</p>
					<ul>
						<li>Some Thing</li>
						<li>Another Thing</li>
					</ul>
				</li>
				<li>Second Thing</li>
			</ol>
			`,
		},
		{
			name: "empty list item",
			html: `
			<ul>
				<li>foo</li>
				<li></li>
				<li>bar</li>
			</ul>
			`,
		},
		{
			name: "list items ending with a space",
			html: `
<ul>
	<li>List items </li>
	<li>Ending with </li>
	<li>A space </li>
</ul>
			`,
		},
		{
			name: "sup element",
			html: `
			<p>One of the most common equations in all of physics is
			<var>E</var>=<var>m</var><var>c</var><sup>2</sup>.<p>
			`,
		},
		{
			name: "sup element in list",
			html: `
			<p>The ordinal number "fifth" can be abbreviated in various languages as follows:</p>
			<ul>
				<li>English: 5<sup>th</sup></li>
				<li>French: 5<sup>ème</sup></li>
			</ul>
			`,
		},
		{
			name: "image",
			html: `<img alt="website favicon" src="http://commonmark.org/help/images/favicon.png" />`,
		},
		{
			name: "image with alt tag",
			html: `<img alt='website "favicon"' src="http://commonmark.org/help/images/favicon.png" />`,
		},
		{
			name: "image inside an empty link",
			html: `	<a href="" title="title">
						<img alt="website favicon" src="http://commonmark.org/help/images/favicon.png" />
					</a>`,
		},
		{
			name: "link",
			html: `<a href="http://commonmark.org/">Link</a>`,
		},
		{
			name: "link with title",
			html: `<a href="http://commonmark.org/" title="Some Text">Link</a>`,
		},
		{
			name: "multiline link",
			html: `
<a href="http://commonmark.org/">
	<p>First Text</p>
	<img src="xxx">
	<p>Second Text</p>
</a>`,
		},
		{
			name: "multiline link inside a list item",
			html: `
<ul>
	<li>
		<a href="http://commonmark.org/">
			<p>First Text</p>
			<br />
			<br />
			<br />
			<br />
			<br />
			<br />
			<br />
			<br />
			<p>Second Text</p>
		</a>
	</li>
</ul>
			`,
		},
		{
			name: "link with svg inlined",
			html: `
<a aria-label="Homepage" title="GitHub" href="https://github.com">
	<svg height="24" viewBox="0 0 16 16" version="1.1" width="24" aria-hidden="true"><path fill-rule="evenodd" d="..."></path></svg>
</a>
			`,
			options: &Options{
				LinkStyle: "inlined",
			},
		},
		{
			name: "link with svg reference link full",
			html: `
<a aria-label="Homepage" title="GitHub" href="https://github.com">
	<svg height="24" viewBox="0 0 16 16" version="1.1" width="24" aria-hidden="true"><path fill-rule="evenodd" d="..."></path></svg>
</a>
			`,
			options: &Options{
				LinkStyle:          "referenced",
				LinkReferenceStyle: "full",
			},
		},
		{
			name: "link with svg reference link collapsed",
			html: `
<a aria-label="Homepage" title="GitHub" href="https://github.com">
	<svg height="24" viewBox="0 0 16 16" version="1.1" width="24" aria-hidden="true"><path fill-rule="evenodd" d="..."></path></svg>
</a>
			`,
			options: &Options{
				LinkStyle:          "referenced",
				LinkReferenceStyle: "collapsed",
			},
		},
		{
			name: "link with svg reference link shortcut",
			html: `
<a aria-label="Homepage" title="GitHub" href="https://github.com">
	<svg height="24" viewBox="0 0 16 16" version="1.1" width="24" aria-hidden="true"><path fill-rule="evenodd" d="..."></path></svg>
</a>
			`,
			options: &Options{
				LinkStyle:          "referenced",
				LinkReferenceStyle: "shortcut",
			},
		},
		{
			name: "tweet",
			html: `
<div class="tweet" data-attrs="{&quot;url&quot;:&quot;https://twitter.com/kroger/status/1271516803756425218&quot;,&quot;full_text&quot;:&quot;As a company, it’s our responsibility to better support our Black associates, customers and allies. We know there is more work to do and will keep you updated on our progress, this is only the beginning. Black Lives Matter. &quot;,&quot;username&quot;:&quot;kroger&quot;,&quot;name&quot;:&quot;Kroger&quot;,&quot;date&quot;:&quot;Fri Jun 12 18:56:44 +0000 2020&quot;,&quot;photos&quot;:[{&quot;img_url&quot;:&quot;https://pbs.substack.com/media/EaVVy4aXsAglkCk.jpg&quot;,&quot;link_url&quot;:&quot;https://t.co/DxScre83q4&quot;}],&quot;quoted_tweet&quot;:{},&quot;retweet_count&quot;:17,&quot;like_count&quot;:93,&quot;expanded_url&quot;:{}}">
<a href="https://twitter.com/kroger/status/1271516803756425218" target="_blank"><div class="tweet-header"><img class="tweet-user-avatar" src="https://cdn.substack.com/image/twitter_name/w_36/kroger.jpg"><span class="tweet-author-name">Kroger </span><span class="tweet-author">@kroger</span></div>
<p>As a company, it’s our responsibility to better support our Black associates, customers and allies. We know there is more work to do and will keep you updated on our progress, this is only the beginning. Black Lives Matter. </p>
<img class="tweet-photo" src="https://cdn.substack.com/image/fetch/w_600,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fpbs.substack.com%2Fmedia%2FEaVVy4aXsAglkCk.jpg"><div class="tweet-footer"><p class="tweet-date">June 12th 2020</p><span class="retweets"><span class="rt-count">17</span> Retweets</span><span class="likes"><span class="like-count">93</span> Likes</span></div></a></div>
			`,
		},
		{
			name: "reference link full",
			html: `
<a href="http://commonmark.org/first">First Link</a>

<a href="http://commonmark.org/second">Second Link</a>
`,
			options: &Options{
				LinkStyle:          "referenced",
				LinkReferenceStyle: "full",
			},
		},
		{
			name: "reference link collapsed",
			html: `<a href="http://commonmark.org/">Link</a>`,
			options: &Options{
				LinkStyle:          "referenced",
				LinkReferenceStyle: "collapsed",
			},
		},
		{
			name: "reference link shortcut",
			html: `<a href="http://commonmark.org/">Link</a>`,
			options: &Options{
				LinkStyle:          "referenced",
				LinkReferenceStyle: "shortcut",
			},
		},
		{
			name: "escape strong",
			html: `<p>**Not Strong**
			**Still Not
			Strong**</p>`,
		},
		{
			name: "escape italic",
			html: `<p>_Not Italic_</p>`,
		},
		{
			name: "escape ordered list",
			html: `<p>1. Not List 1. Not List
			1. Not List</p>`,
		},
		{
			name: "escape unordered list",
			html: `<p>- Not List</p>`,
		},
		{
			name: "pre tag",
			html: `
			<div>
				<p>Who ate the most donuts this week?</p>
				<pre><code class="language-foo+bar">Jeff  15
Sam   11
Robin  6</code></pre>
			</div>
			`,
		},
		{
			name: "code tag inside p",
			html: `
			<p>When <code>x = 3</code>, that means <code>x + 2 = 5</code></p>
			`,
		},
		{
			name: "code tag",
			html: `
			<code>last_30_days</code>
			`,
		},
		{
			name: "hr",
			html: `
			<p>Some Content</p>
			<hr>
			</p>Other Content</p>
			`,
		},
		{
			name: "blockquote",
			html: `
<blockquote>
Some Quote
Next Line
</blockquote>
			`,
		},
		{
			name: "large blockquote",
			html: `
			<blockquote>
				<p>Allowing an unimportant mistake to pass without comment is a wonderful social grace.</p>
				<p>Ideological differences are no excuse for rudeness.</p>
			</blockquote>
			`,
		},
		{
			name: "empty blockquote",
			html: `
<blockquote></blockquote>
			`,
		},
		{
			name: "turndown demo",
			html: `
			<h1>Turndown Demo</h1>

			<p>This demonstrates <a href="https://github.com/domchristie/turndown">turndown</a> – an HTML to Markdown converter in JavaScript.</p>

			<h2>Usage</h2>

			<pre><code class="language-js">var turndownService = new TurndownService()
console.log(
  turndownService.turndown('&lt;h1&gt;Hello world&lt;/h1&gt;')
)</code></pre>

			<hr />

			<p>It aims to be <a href="http://commonmark.org/">CommonMark</a>
			 compliant, and includes options to style the output. These options include:</p>

			<ul>
			  <li>headingStyle (setext or atx)</li>
			  <li>horizontalRule (*, -, or _)</li>
			  <li>bullet (*, -, or +)</li>
			  <li>codeBlockStyle (indented or fenced)</li>
			  <li>fence</li>
			  <li>emDelimiter (_ or *)</li>
			  <li>strongDelimiter (** or __)</li>
			  <li>linkStyle (inlined or referenced)</li>
			  <li>linkReferenceStyle (full, collapsed, or shortcut)</li>
			</ul>
						`,
		},
		{
			name: "keep tag",
			html: `<keep-tag><p>Content</p></keep-tag>`,
		},
		{
			name: "remove tag",
			html: `<remove-tag><p>Content</p></remove-tag>`,
		},
		{
			/*
				When a header (eg. <h3>) contains any new lines in its body, it will split the header contents
				over multiple lines, breaking the header in Markdown (because in Markdown, a header just
				starts with #'s and anything on the next line is not part of the header). Since in HTML
				and Markdown all white space is treated the same, I chose to replace line endings with spaces.
				-> https://github.com/lunny/html2md/pull/6
			*/
			name: "strip newlines from header",
			html: `
<h3>

Header
Containing

Newlines

</h3>
			`,
		},
		{
			name: "text with whitespace",
			html: `
						<div id="sport_single_post-2" class="widget sport_single_post">
			<h1 class="widget-title">Aktuelles</h1>
			
			<!-- featured image -->
			<div class="mosaic-block fade"><a href="http://www.bonnerruderverein.de/wp-content/uploads/2015/09/BRV-abend.jpg" class="mosaic-overlay fancybox" title="BRV-abend"></a><div class="mosaic-backdrop"><div class="corner-date">25 Mai</div><img src="http://www.bonnerruderverein.de/wp-content/uploads/2015/09/BRV-abend.jpg" alt="" /></div></div>
			<!-- title -->
			<h3 class="title"><a href="http://www.bonnerruderverein.de/bonner-nachtlauf/">9. Bonner Nachtlauf - Einschränkungen am Bootshaus</a></h3>

            <!-- excerpt -->
            am Mittwoch, dem 30. Mai 2018 findet am Bonner Rheinufer der 9. ...
            <a href="http://www.bonnerruderverein.de/bonner-nachtlauf/" class="more">More</a>



			</div>

			<hr />
			
		<div>
			<h1 class="widget-title">Aktuelles</h1>
			<h3 class="title"><a href="some_url">Title</a></h3>

						<!-- excerpt -->
						Fusce dapibus, tellus ac cursus commodo, tortor mauris condimentum nibh, ut fermentum massa justo sit amet risus. Vestibulum id ligula porta felis euismod semper.
						<a href="other_url" class="more">More</a>

		</div>
`,
		},
		{
			name: "pre tag without code tag",
			html: `
<div class="code"><pre>// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {</pre></div>
`,
		},
		{
			name: "escape pipe characters because of the use in tables",
			html: `<p>With | Character<p>`,
		},
		{
			name: "br adds new line break",
			html: `<p>1. xxx <br/>2. xxxx<br/>3. xxx</p><p><span class="img-wrap"><img src="xxx"></span><br>4. golang<br>a. xx<br>b. xx</p>`,
		},
		{
			name: "br does not add new line inside header",
			html: `<h1>Heading<br/> <br/>One</h1>`,
		},
		{
			name: "dont escape too much",
			html: `jmap –histo[:live]`,
		},
		/*
					{ // TODO: not working yet
						name: "p tag with lots of whitespace",
						html: `
			<p>
				Sometimes a struct field, function, type, or even a whole package becomes


				redundant or unnecessary, but must be kept for compatibility with existing


				programs.


				To signal that an identifier should not be used, add a paragraph to its doc


				comment that begins with "Deprecated:" followed by some information about the


				deprecation.


				There are a few examples <a href="https://golang.org/search?q=Deprecated:" target="_blank">in the standard library</a>.
			</p>
			`,
					},
		*/
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			converter := NewConverter(test.domain, true, test.options)
			converter.Keep("keep-tag").Remove("remove-tag")

			markdown, err := converter.ConvertString(test.html)
			if err != nil {
				t.Error(err)
			}
			data := []byte(markdown)

			name := strings.Replace(test.name, " ", "_", -1)
			if url.QueryEscape(name) != name {
				fmt.Printf("'%s' is not a safe path", name)
			}
			name = url.QueryEscape(name) // for safety
			name = filepath.Join("TestFromString", name)

			gp := filepath.Join("testdata", name+".golden")
			if *update {
				t.Log("update golden file")
				if err := ioutil.WriteFile(gp, data, 0644); err != nil {
					t.Fatalf("failed to update golden file: %s", err)
				}
			}

			g, err := ioutil.ReadFile(gp)
			if err != nil {
				t.Logf("Result:\n'%s'\n", markdown)
				t.Fatalf("failed reading .golden: %s", err)
			}

			if !bytes.Equal([]byte(markdown), g) {
				dmp := diffmatchpatch.New()

				diffs := dmp.DiffMain(string(g), markdown, false)
				diffs = dmp.DiffCleanupSemantic(diffs)

				fmt.Println(dmp.DiffToDelta(diffs))

				t.Errorf("written json does not match .golden file: %+v \n", dmp.DiffPrettyText(diffs))
			}
		})
	}
}

func BenchmarkFromString(b *testing.B) {
	converter := NewConverter("www.google.com", true, nil)

	strongRule := Rule{
		Filter: []string{"strong"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			return nil
		},
	}

	var wg sync.WaitGroup
	convert := func(html string) {
		defer wg.Done()
		_, err := converter.ConvertString(html)
		if err != nil {
			b.Error(err)
		}
	}
	add := func() {
		defer wg.Done()
		converter.AddRules(strongRule)
	}

	for n := 0; n < b.N; n++ {
		wg.Add(2)
		go add()
		go convert("<strong>Bold</strong>")
	}

	wg.Wait()
}

func TestAddRules_ChangeContent(t *testing.T) {
	expected := "Some other Content"

	var wasCalled bool
	rule := Rule{
		Filter: []string{"p"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			wasCalled = true

			if content != "Some Content" {
				t.Errorf("got wrong `content`: '%s'", content)
			}
			if !selec.Is("p") {
				t.Error("selec is not p")
			}
			return String(expected)
		},
	}

	conv := NewConverter("", true, nil)
	conv.AddRules(rule)
	md, err := conv.ConvertString(`<p>Some Content</p>`)
	if err != nil {
		t.Error(err)
	}

	if md != expected {
		t.Errorf("wanted '%s' but got '%s'", expected, md)
	}
	if !wasCalled {
		t.Error("rule was not called")
	}
}

func TestAddRules_Fallback(t *testing.T) {
	// firstExpected := "Some other Content"
	expected := "Totally different Content"

	var firstWasCalled bool
	var secondWasCalled bool
	firstRule := Rule{
		Filter: []string{"p"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			firstWasCalled = true
			if secondWasCalled {
				t.Error("expected first rule to be called before second rule. second is already called")
			}

			if content != "Some Content" {
				t.Errorf("got wrong `content`: '%s'", content)
			}
			if !selec.Is("p") {
				t.Error("selec is not p")
			}

			return nil
		},
	}
	secondRule := Rule{
		Filter: []string{"p"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			secondWasCalled = true
			if !firstWasCalled {
				t.Error("expected first rule to be called before second rule. first is not called yet")
			}

			if content != "Some Content" {
				t.Errorf("got wrong `content`: '%s'", content)
			}
			if !selec.Is("p") {
				t.Error("selec is not p")
			}
			return String(expected)
		},
	}

	conv := NewConverter("", true, nil)
	conv.AddRules(secondRule, firstRule)
	md, err := conv.ConvertString(`<p>Some Content</p>`)
	if err != nil {
		t.Error(err)
	}

	if md != expected {
		t.Errorf("wanted '%s' but got '%s'", expected, md)
	}
	if !firstWasCalled {
		t.Error("first rule was not called")
	}
	if !secondWasCalled {
		t.Error("second rule was not called")
	}
}
func TestWholeSite(t *testing.T) {
	var tests = []struct {
		name   string
		domain string

		file string
	}{
		{
			name: "golang.org",

			domain: "golang.org",
		},
		{
			name:   "bonnerruderverein.de",
			domain: "bonnerruderverein.de",
		},
		{
			name:   "blog.golang.org",
			domain: "blog.golang.org",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			converter := NewConverter(test.domain, true, nil)

			htmlData, err := ioutil.ReadFile(
				filepath.Join("testdata", t.Name()+".html"),
			)
			if err != nil {
				t.Error(err)
			}

			markdownData, err := converter.ConvertBytes(htmlData)
			if err != nil {
				t.Error(err)
			}

			// output := blackfriday.Run(data)
			// fmt.Println(string(output))

			gp := filepath.Join("testdata", t.Name()+".md")
			if *update {
				t.Log("update golden file")
				if err := ioutil.WriteFile(gp, markdownData, 0644); err != nil {
					t.Fatalf("failed to update golden file: %s", err)
				}
			}

			g, err := ioutil.ReadFile(gp)
			if err != nil {
				t.Logf("Result:\n'%s'\n", string(markdownData))
				t.Fatalf("failed reading .golden: %s", err)
			}

			if !bytes.Equal(markdownData, g) {
				t.Errorf("written json does not match .golden file \nexpected:\n'%s'\nbut got:\n'%s'", string(g), string(markdownData))
			}
		})
	}
}
