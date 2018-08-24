package main

import (
	"fmt"
)

func (DigestWriter DigestWriter) headline(Unit string, Digest Digest) string {
	return fmt.Sprintf(`<br></br> <h2>This %v's %v Posts from /r/%v </h2>`, Unit, Digest.NumPosts, Digest.Subreddit)
}

func (DigestWriter DigestWriter) writeHeader() {
	DigestWriter.Email += `<h1>Email Digest for you!</H1> <br></br>`
}

func (DigestWriter DigestWriter) writeFooter() {
	DigestWriter.Email += "<br></br><br></br>Stay cool <br></br> -RedditDigest Bot"
}

func (DigestWriter DigestWriter) writeSpacer() {
	DigestWriter.Email += "<br></br><br><hr></hr><hr></hr>"
}
