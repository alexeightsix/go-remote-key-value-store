simple key/value store written in GO to aid with learning the language

support for SET, GET, DEL

TCP with custom protocol is used

support for different store drivers (array and map)

set and del events get written to disk

store will re-hydrate on init by reading the file

tests were written to ensure functionality works as expected
