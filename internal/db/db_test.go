package db

import (
    "testing"
    "fmt"
    "github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
    CreateUser(User{Username: "james", Email: "test@example.com", Registered_ip: "127.0.0.1"})
    var res CompleteUser = SearchUser(User{Username: "james"})
    var coreUserData User = res.User
   
    assert.Equal(t, coreUserData.Username, "james", "select statement for sample user (username)")
    assert.Equal(t, coreUserData.Email, "test@example.com", "select statement for sample user (email)")
    assert.Equal(t, coreUserData.Registered_ip, "127.0.0.1", "select statement for sample user (ip address)")
    DeleteUser(User{Username: "james"})
}

func TestUserBioCreation(t *testing.T) {
    var james User = User{Username: "james", Email: "test@example.com", Registered_ip: "127.0.0.1"}
    var jamesPersonalDetails UserMetadata = UserMetadata{Username: "james", Full_name: "James Johnson", Birthdate: "02/03/2001", Bio_text: "blahblahblahblahblahblah"}
    CreateUser(james)
    UpsertUserMetadata(jamesPersonalDetails)
    var res UserMetadata = SearchUser(james).Metadata
    assert.Equal(t, res.Full_name, "James Johnson", "metadata full name")
    jamesPersonalDetails.Bio_text = "this is my new bio"
    UpsertUserMetadata(jamesPersonalDetails)
    res = SearchUser(james).Metadata
    assert.Equal(t, res.Bio_text, "this is my new bio", "metadata test bio update")
    DeleteUser(james)
    var res2 UserMetadata = SearchUser(james).Metadata
    assert.Equal(t, res2.Full_name, "", "metadata ensure postgres cascade delete")
}

func TestSubmissionCreation(t *testing.T) {
    // need a test user in the database due to FKs
    var james User = User{Username: "james", Email: "test@example.com", Registered_ip: "127.0.0.1"}
    CreateUser(james)
    var testSubmission Submission = Submission{Username: "james", Title: "Google Search", Link: "https://www.google.com", Body: "here's a search engine", Flagged: true}
    var generatedID string = CreateSubmission(testSubmission)
    var searchedSubmission Submission = SearchSubmission(Submission{Id: generatedID})
    assert.Equal(t, searchedSubmission.Body, testSubmission.Body, "submission insert, check body is the same")
    assert.Equal(t, searchedSubmission.Flagged, testSubmission.Flagged, "submission insert, check flagged is the same")
    assert.Equal(t, searchedSubmission.Link, testSubmission.Link, "submission insert, check link is the same")
    assert.Equal(t, searchedSubmission.Title, testSubmission.Title, "submission insert, check title is the same")
    DeleteSubmission(Submission{Id: generatedID})
    searchedSubmission = SearchSubmission(Submission{Id: generatedID})
    assert.Equal(t, searchedSubmission.Link, "", "submission delete, ensure link is blank")
    DeleteUser(james)
}

func TestCreateMagicLink(t *testing.T) {
    // need a test user in the database due to FKs
    var james User = User{Username: "james", Email: "test@example.com", Registered_ip: "127.0.0.1"}
    CreateUser(james)
    var insertedToken string = CreateMagicLink(james)
    fmt.Println(insertedToken)
    DeleteUser(james)
}

func TestCreateMagicLinkTwo(t *testing.T) {
    fmt.Println(CreateMagicLink(User{Email: "me@trentwil.es"}))
}

func TestCreateRandomData(t *testing.T) {
    GenerateNonsenseData(10, 300);
}