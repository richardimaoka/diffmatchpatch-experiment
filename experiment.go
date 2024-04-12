package main

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = `I am the very model of a modern Major-General,
I've information vegetable, animal, and mineral,
I know the kings of England, and I quote the fights historical,
From Marathon to Waterloo, in order categorical.`
	text2 = `I am the very model of a cartoon individual,
My animation's comical, animal, and whimsical,
I'm quite adept at funny gags, comedic theory I have read,
From wicked puns and stupid jokes to anvils that drop on your head.`
)

func experiment() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, true)
	currentLine := 1
	for i, d := range diffs {
		nCount := strings.Count(d.Text, "\n")
		currentLine += nCount
		fmt.Printf("[%d], %6s, slash-n count = %d, current = %d: %s\n", i, d.Type, nCount, currentLine, strings.ReplaceAll(d.Text, "\n", "\\n"))
	}

	patch := `@@ -54,5 +54,8 @@
         </div>
       </div>
     </div>
+    <button>1</button>
+    <button>2</button>
+    <button>3</button>
   </body>
 </html>`

	fmt.Printf("\\n count = %d", strings.Count(diffs[1].Text, "\n"))

	fmt.Println("DiffPrettyText:", dmp.DiffPrettyText(diffs))
	fmt.Println("DiffPrettyHtml:", dmp.DiffPrettyHtml(diffs))

	patches, err := dmp.PatchFromText(patch)
	if err != nil {
		fmt.Println(err)
	}
	for _, p := range patches {
		fmt.Println("p:", p.String())
	}

	textA := `FROM node:18-alpine
WORKDIR /app
ADD package.json /app 
RUN npm i --silent
ADD . /app 
CMD npm run start`

	textB := `
FROM node:18-alpine
WORKDIR /src
ADD package.json /src 
RUN npm i --silent
ADD . /src 
CMD npm start`
	diffs = dmp.DiffMain(textA, textB, true)
	fmt.Println("DiffPrettyText:", dmp.DiffPrettyText(diffs))

	s1, s2, lines := dmp.DiffLinesToChars(text1, text2)
	fmt.Println("-----------------------")
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(lines)

}
