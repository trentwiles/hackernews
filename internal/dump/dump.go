package dump

import (
	"encoding/json"
	"log"
	"os"

	"github.com/trentwiles/hackernews/internal/db"
)

func DumpForUser(user db.User) bool {
    // BEFORE RUNNING, we assume user exists and is authorized to access this data (that'll be handled via the API)
    // included in a user dump:
    // 1. user metadata
    // 2. posts
    // 3. comments
    // 4. up/downvotes
    db.UpdateSelectLimit(5000)
    var userMeta db.CompleteUser = db.SearchUser(user)
    var userSubmissions []db.BasicSubmission = db.LatestUserSubmissions(0, user) // pass 0 as offset, since we're working with a high limit
    var userComments []db.Comment = db.LatestUserComments(0, user)
    var userVotes []db.BasicSubmissionAndVote = db.GetAllUserVotes(user)
    
    exportDir := "exports/" + user.Username
    os.MkdirAll(exportDir, 0755)
    
    if !writeJSONToFile(userMeta, exportDir+"/user.json") {
        return false
    }
    if !writeJSONToFile(userSubmissions, exportDir+"/submissions.json") {
        return false
    }
    if !writeJSONToFile(userComments, exportDir+"/comments.json") {
        return false
    }
    if !writeJSONToFile(userVotes, exportDir+"/votes.json") {
        return false
    }
    
    return true
}

func writeJSONToFile(data interface{}, filepath string) bool {
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        log.Fatalln("[ERROR] Error marshaling JSON:", err)
        return false
    }
    
    err = os.WriteFile(filepath, jsonData, 0644)
    if err != nil {
        log.Fatalln("[ERROR] Error writing file:", err)
        return false
    }
    
    return true
}