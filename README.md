# hackernews

## Feature Implementation Checklist

| Feature                | API Status     | Web Status     | Notes                                                                                                       |
| ---------------------- | -------------- | -------------- | ----------------------------------------------------------------------------------------------------------- |
| Login via Magic Link   | ✅ Complete    | ✅ Complete    | JWT-based authentication, cannot be done over `eduroam` because they block outbound email port connections. |
| Magic Link accept page | ✅ Complete    | ✅ Complete    | Page where the server validates the magic link found in the email, and adds the token to browser cookies    |
| News/Submission Feed   | 🟡 In Progress | 🟡 In Progress | Need to complete different kinds of sorts (newest, best, oldest), and requires pagination                   |
| Post Submission        | ✅ Complete    | ⬜ Not Started |                                                                                                             |
| User Profile Page      | ✅ Complete    | ⬜ Not Started |                                                                                                             |
| Comments               | ⬜ Not Started | ⬜ Not Started | SQL implementation is complete                                                                              |
| Admin Console          | ⬜ Not Started | ⬜ Not Started |                                                                                                             |

## To Do

- Modularize routes
- Improved error handling (ie. stop using `log.Fatal` and start using `fmt.Errorf`)
- Write API docs

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
