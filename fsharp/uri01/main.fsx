open System

// Get a parent path from a URL
// Based on
// - https://stackoverflow.com/questions/510240/getting-the-parent-name-of-a-uri-url-from-absolute-name-c-sharp#comment14991159_5025018

let url = "https://syohex.org/aaa/bbb/ccc/ddd.jpg"
let uri = new Uri(url)
let parent = new Uri(uri, ".")
printfn "orig=%s, parent=%s" (uri.ToString()) (parent.ToString())
