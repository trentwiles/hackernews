# hackernews

[Periodic token deletion](https://chatgpt.com/s/t_68601c424ef8819187259e6afde1c01e)

* note: in the future, on the vote endpoint, return the total number of votes

## To Do
* Modularize routes
* Improved error handling (ie. stop using `log.Fatal` and start using `fmt.Errorf`)
* Write API docs

## Route Modularization File Structure

```
|--main.go
|--routes/
|    |--user.go
|    |--login.go
|--handlers/
|    |--user_handler.go
|    |--login_handler.go
```