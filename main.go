package main

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = `I am the very model of a modern Major-General,
I've information vegetable, animal, and mineral,
I know the kings of England, and I quote the fights historical,
From Marathon to Waterloo, in order categorical.`
	text2 = `I am the very model of a cartoon individual,
My animation's comical, unusual, and whimsical,
I'm quite adept at funny gags, comedic theory I have read,
From wicked puns and stupid jokes to anvils that drop on your head.`
)

func main() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, true)

	patch := `@@ -54,5 +54,8 @@
         </div>
       </div>
     </div>
+    <button>1</button>
+    <button>2</button>
+    <button>3</button>
   </body>
 </html>`

	fmt.Println("DiffPrettyText:", dmp.DiffPrettyText(diffs))
	fmt.Println("DiffText1     :", dmp.DiffText1(diffs))
	fmt.Println("DiffText2     :", dmp.DiffText2(diffs))

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

}
