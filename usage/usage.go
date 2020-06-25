package usage

import (
	"fmt"
	"net/http"
)

const usage = `<html><head><title>Not found</title></head><link rel="icon" href="data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>👤</text></svg>"><body>
<h1>Not found</h1>

<h2>Twitter</h2>
<p>
Give a username and get an avatar in return: <a href="https://avatars.charlyx.dev/twitter?username=charlyx">https://avatars.charlyx.dev/twitter?username=charlyx</a>
</p>

<p>
You can ask for variant sizings such as "bigger", "mini" and "original" (default size being "normal").
</p>

<p>
Example: <a href="https://avatars.charlyx.dev/twitter?username=charlyx&size=bigger">https://avatars.charlyx.dev/twitter?username=charlyx&size=bigger</a>
</p>

<h2>Gravatar</h2>

<p>
Give an email and get an avatar in return: <a href="https://avatars.charlyx.dev/gravatar?email=mon@email">https://avatars.charlyx.dev/gravatar?email=mon@email</a>
</p>

<p>
By default, images are presented at 80px by 80px if no size parameter is supplied.<br>
You may request a specific image size from 1px up to 2048px by using the s= or size= parameter and passing a single pixel dimension (since the images are square).
</p>

<p>
Example: <a href="https://avatars.charlyx.dev/gravatar?email=mon@email&s=200">https://avatars.charlyx.dev/gravatar?email=mon@email&s=200</a>
</p>
</body>
</html>`

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, usage)
}
