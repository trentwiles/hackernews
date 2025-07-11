# hackernews

# Versioning Goals

## Alpha (MVP, v0.1.0)

### Features

- Magic Link login system (password free, managed by JWT tokens)
- Submission Posting (with full Google Captcha v3 support)
- Submission Voting (upvotes/downvotes, options to sort by votes)
- User bio/birthday/full name customization on "account settings" page
- Command line logging (HTTP request logging + logging from business logic files)

### In Progress
- Basic admin page (basic website metrics)
- Name change!

## Beta ("version two", v0.2.0)

### In Progress

- Comments (and comment threads)
- Ability to report a post (and maybe a user?)
- `.env` usage on the frontend
- Fully working admin page (view/delete flagged posts, delete users, view full website metrics)
- Clean up error handling (mostly using Log.fatal now, not great for production as this kills the program)

## Production ("version three", v1.0.0)

- Logging to a file, which is then synced with a S3 Bucket

<hr>

## Feature Implementation Checklist

| Feature                | API Status     | Web Status     | Notes                                                                                                       |
| ---------------------- | -------------- | -------------- | ----------------------------------------------------------------------------------------------------------- |
| Login via Magic Link   | ✅ Complete    | ✅ Complete    | JWT-based authentication, cannot be done over `eduroam` because they block outbound email port connections. |
| Magic Link accept page | ✅ Complete    | ✅ Complete    | Page where the server validates the magic link found in the email, and adds the token to browser cookies    |
| News/Submission Feed   | 🟡 In Progress | 🟡 In Progress | Need to complete different kinds of sorts (newest, best, oldest), and requires pagination                   |
| Post Submission        | ✅ Complete    | ✅ Complete    |                                                                                                             |
| User Profile Page      | ✅ Complete    | ✅ Complete    |                                                                                                             |
| Comments               | ⬜ Not Started | ⬜ Not Started | SQL implementation is complete                                                                              |
| Admin Console          | 🟡 In Progress | ⬜ Not Started | SQL implementation completed; added isAdmin to user API                                                     |

## To Do

- Modularize routes
- Improved error handling (ie. stop using `log.Fatal` and start using `fmt.Errorf`)
- New name

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

## Logging

HTTP request logging is done by a built in middleware for Fiber. Other logs are placed in the database and backend "core" logic files.

| Type      | Description                                                      |
| --------- | ---------------------------------------------------------------- |
| `[INFO]`  | Basic event, not very important                                  |
| `[WARN]`  | Some sort of failure took place, wasn't fatal                    |
| `[FATAL]` | Error that killed the program took place, pay attention to these |

Future: toggle logging via `.env` (`PRODUCTION=true`), see this [ChatGPT conversion](https://chatgpt.com/share/68701ac1-5848-800a-9325-aea6489727a8) for more on that.