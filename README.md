
# BFF

Writing internal tools sucks. You just want a few dumb buttons and grids of output that people can look at, and this is the best way to do it. You are a backend engineer that loves writing Go, or has a go program already, then this is a super easy to integrate framework that will give you a backend user interfaces that ALL of your organisation can easily use.

[Read more about the problems I'm solving with BFF](https://www.ersin.nz/articles/bff)



# TODO
- landing page for this tool
- make it so that inputs are responded to by the backend after the message is recieved: I.E  synchronous response from backend for a submitted message
  - Make it reload from half finished state (I.E resume after reconnection/service restart) 
- Should reconnect to the socket and continue execution
- lots of io componetns TODO see io.go
- documentation for people that want to pull this in as a library
- deploy it somewhere

# done
- more pretty styles that change input state based on submitting it or not
- blog post about this prototype
