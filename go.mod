module github.com/eolso/chat

go 1.18

require (
	github.com/eolso/memcache v0.0.0-20220501062927-10a730408c04
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/render v1.0.1
	github.com/rs/xid v1.3.0
	github.com/rs/zerolog v1.26.1
)

require github.com/eolso/threadsafe v0.0.0-20220426062731-233f2996097c // indirect

replace github.com/eolso/threadsafe => ../threadsafe
replace github.com/eolso/memcache => ../memcache
