package main

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = "Lorem ipsum dolor."
	text2 = "Lorem dolor sit amet."
)

func main() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, false)

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

}
