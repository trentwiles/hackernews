import WebSidebar from "./fragments/WebSidebar";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";
import {
  ArrowDown,
  ArrowUp,
  ArrowUpRight,
  Clock,
  Share,
  Trash,
  User as UserIcon,
  Send,
} from "lucide-react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { Button } from "./ui/button";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { datePrettyPrint, getTimeAgo } from "@/utils";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@radix-ui/react-tooltip";
import Cookies from "js-cookie";
import { Textarea } from "./ui/textarea";
import { useGoogleReCaptcha } from "react-google-recaptcha-v3";

type submission = {
  username: string;
  title: string;
  link: string;
  body: string;
  created_at: string;
  upvotes: number;
  downvotes: number;
  totalScore: number;
};

type comment = {
  id: string;
  author: string;
  content: string;
  flagged: boolean;
  created_at: string;
  parent?: string;
  upvotes: number;
  downvotes: number;
  isUpvoted: boolean;
  isDownvoted: boolean;
};

export default function Submission() {
  const { executeRecaptcha } = useGoogleReCaptcha();
  const { sid } = useParams();
  const navigate = useNavigate();
  const [s, setS] = useState<submission>();
  const [pending, setPending] = useState<boolean>(true);
  const [error, setError] = useState<boolean>(false);
  const [currentUser, setCurrentUser] = useState<string>();

  const [upvoteEnabled, setUpvoteEnabled] = useState<boolean>(true);
  const [downvoteEnabled, setDownvoteEnabled] = useState<boolean>(true);
  const [canVote, setCanVote] = useState<boolean>(false);

  const [shareButtonText, setShareButtonText] = useState<string>("Share");

  const [comments, setComments] = useState<comment[]>([]);
  const [newComment, setNewComment] = useState<string>("");
  const [submittingComment, setSubmittingComment] = useState<boolean>(false);

  useEffect(() => {
    if (sid === undefined || sid == "") {
      setPending(false);
      setError(true);
      return;
    }

    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/submission?id=" + sid)
      .then((res) => {
        if (res.status === 404) {
          navigate("/404");
          return;
        }
        if (!res.ok) throw new Error("Network response was not ok");
        return res.json();
      })
      .then((data) => {
        const current: submission = {
          body: data.metadata.body,
          created_at: data.metadata.createdAt,
          link: data.metadata.link,
          title: data.metadata.title,
          username: data.metadata.author,
          downvotes: data.votes.downvotes,
          totalScore: data.votes.total,
          upvotes: data.votes.upvotes,
        };

        setS(current);
        setPending(false);
        setError(false);
      })
      .catch((err) => {
        toast("Error: " + err);
        setPending(false);
        setError(true);
        console.log(err);
      });
  }, [sid, upvoteEnabled, downvoteEnabled]);
  //       ^^^^           ^^^^
  // whenever the user votes or downvotes, we refresh the total votes

  useEffect(() => {
    setCurrentUser(Cookies.get("username"));
  }, []);

  // set the vote button, if the user voted
  useEffect(() => {
    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/vote?id=" + sid, {
      headers: {
        Authorization: "Bearer " + Cookies.get("token"),
      },
    })
      .then((res) => {
        if (res.status === 401) {
          setCanVote(false);
          throw new Error("Unauthorized"); // This will skip to catch
        }

        if (!res.ok) throw new Error("Network response was not ok");
        return res.json();
      })
      .then((data) => {
        setCanVote(true);
        if (!data.didVote) {
          setUpvoteEnabled(true);
          setDownvoteEnabled(true);
          return;
        }
        if (data.didUpvote) {
          setUpvoteEnabled(false);
          setDownvoteEnabled(true);
        } else {
          setUpvoteEnabled(true);
          setDownvoteEnabled(false);
        }
      })
      .catch((err) => {
        // whatever
        console.error(err);
        return;
      });
  }, [sid, upvoteEnabled, downvoteEnabled]); // <--- fixes the upvote/downvote malfunction

  // fetch the comments
  const fetchComments = () => {
    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/comments?id=" + sid + "&username=" + (Cookies.get("username") || ""), {
      headers: {
        Authorization: "Bearer " + Cookies.get("token"),
      },
    })
      .then((res) => {
        if (res.status != 200) {
          throw new Error("unable to fetch comments"); // This will skip to catch
        }

        if (!res.ok) throw new Error("Network response was not ok");
        return res.json();
      })
      .then((data) => {
        if (data.comments == null || data.comments.length == 0) {
          return;
        } else {
          const res: comment[] = []
          data.comments.forEach((item) => {
            const tmp: comment = {
              author: item.Author,
              content: item.Content, 
              created_at: item.CreatedAt,
              flagged: item.Flagged,
              id: item.Id,
              parent: item.ParentComment,
              upvotes: item.Upvotes,
              downvotes: item.Downvotes,
              isUpvoted: item.HasUpvoted,
              isDownvoted: item.HasDownvoted
            }
            res.push(tmp)
          })

          setComments(res)
          setPending(false)
        }
      })
      .catch((err) => {
        // whatever
        console.error(err);
        setPending(false)
        return;
      });
  };

  useEffect(() => {
    fetchComments();
  }, [sid]);

  function deletePost() {
    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/submission", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + Cookies.get("token"),
      },
      body: JSON.stringify({
        Id: sid,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error, status: ${response.status}`);
        }

        return response.json();
      })
      .then((data) => {
        console.log(data);

        navigate("/?deleted=" + sid);
        return;
      })
      .catch((error) => {
        console.log(error);
      });
  }

  function vote(intent: boolean) {
    if (intent) {
      setUpvoteEnabled(false);
      setDownvoteEnabled(true);
    } else {
      setDownvoteEnabled(false);
      setUpvoteEnabled(true);
    }
    // Id: req.Id}, req.Upvote
    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/vote", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + Cookies.get("token"),
      },
      body: JSON.stringify({
        Id: sid,
        Upvote: intent,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error, status: ${response.status}`);
        }

        return response.json();
      })
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.log(error);
      });
  }

  // Add comment voting function
  function voteComment(commentId: string, intent: boolean) {
    // Optimistically update UI
    setComments(prevComments => 
      prevComments.map(comment => {
        if (comment.id === commentId) {
          const newComment = { ...comment };
          
          // Reset vote counts based on previous state
          if (comment.isUpvoted) {
            newComment.upvotes--;
          }
          if (comment.isDownvoted) {
            newComment.downvotes--;
          }
          
          // Apply new vote
          if (intent) {
            newComment.upvotes++;
            newComment.isUpvoted = true;
            newComment.isDownvoted = false;
          } else {
            newComment.downvotes++;
            newComment.isUpvoted = false;
            newComment.isDownvoted = true;
          }
          
          return newComment;
        }
        return comment;
      })
    );

    // Make API call
    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/comment/vote", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + Cookies.get("token"),
      },
      body: JSON.stringify({
        CommentId: commentId,
        Upvote: intent,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error, status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        // Refresh comments to get accurate data
        fetchComments();
      })
      .catch((error) => {
        console.error("Error voting on comment:", error);
        toast.error("Failed to vote on comment");
        // Revert optimistic update on error
        fetchComments();
      });
  }

  async function submitComment() {
    if (!newComment.trim()) {
      toast.error("Comment cannot be empty");
      return;
    }

    if (!currentUser) {
      toast.error("You must be logged in to comment");
      return;
    }

    if (!executeRecaptcha) {
      toast.error("Execute recaptcha not yet available");
      return;
    }

    const token = await executeRecaptcha("login");

    setSubmittingComment(true);

    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/comment", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + Cookies.get("token"),
      },
      body: JSON.stringify({
        InResponseTo: sid,
        Content: newComment,
        CaptchaToken: token,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error, status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        toast.success("Comment posted successfully");
        setNewComment("");

        console.log(data)

        // then, refresh comments
        fetchComments();
        setSubmittingComment(false);
      })
      .catch((error) => {
        console.error(error);
        toast.error("Failed to post comment");
        setSubmittingComment(false);
      });
  }

  console.log(canVote + " <-- can vote???");

  return (
    !pending &&
    !error &&
    s !== undefined && (
      <SidebarProvider>
        <div className="flex min-h-screen w-full">
          <WebSidebar />
          <SidebarInset>
            <div className="min-h-screen bg-gray-50 py-8 px-4">
              <div className="max-w-3xl mx-auto space-y-6">
                <Card className="shadow-sm">
                  <CardHeader className="space-y-1">
                    <CardTitle className="text-2xl font-bold leading-tight">
                      {s.title}
                    </CardTitle>
                    <CardDescription className="flex items-center gap-2 text-base">
                      <ArrowUpRight className="h-4 w-4" />
                      <a
                        href={s.link}
                        className="text-blue-600 hover:text-blue-800 hover:underline break-all"
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {s.link}
                      </a>
                    </CardDescription>
                  </CardHeader>

                  <CardContent className="space-y-4">
                    <div className="prose prose-gray max-w-none">
                      <p className="text-gray-700 leading-relaxed">{s.body}</p>
                    </div>

                    <div className="flex items-center gap-4 pt-4 border-t">
                      <div className="flex items-center gap-2 text-sm text-gray-600">
                        <UserIcon className="h-4 w-4" />
                        <span>
                          Posted by&nbsp;
                          <Link
                            to={"/u/" + s.username}
                            className="font-medium text-gray-900 hover:text-blue-600 hover:underline"
                          >
                            {s.username}
                          </Link>
                        </span>
                      </div>

                      <div className="flex items-center gap-2 text-sm text-gray-600">
                        <Clock className="h-4 w-4" />
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <p>{getTimeAgo(s.created_at)}</p>
                          </TooltipTrigger>
                          <TooltipContent>
                            <span>{datePrettyPrint(s.created_at)}</span>
                          </TooltipContent>
                        </Tooltip>
                      </div>
                    </div>
                  </CardContent>

                  <CardFooter className="bg-gray-50 border-t">
                    <div className="flex items-center justify-between w-full">
                      <div className="flex items-center gap-4">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={async () => {
                            await navigator.clipboard.writeText(
                              window.location.href
                            );
                            setShareButtonText("Copied To Clipboard");
                            await new Promise((r) => setTimeout(r, 2000));
                            setShareButtonText("Share");
                          }}
                        >
                          <Share /> {shareButtonText}
                        </Button>
                        {currentUser == s.username && (
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => deletePost()}
                          >
                            <Trash />
                            <span className="font-medium">Delete</span>
                          </Button>
                        )}
                      </div>

                      <div className="flex items-center gap-2">
                        {canVote && (
                          <Button
                            variant={!upvoteEnabled ? "destructive" : "outline"}
                            size="sm"
                            disabled={!upvoteEnabled}
                            onClick={() => vote(true)}
                          >
                            <ArrowUp />
                          </Button>
                        )}
                        <Button variant="outline" size="sm">
                          <span className="font-medium">
                            {s.totalScore} upvotes
                          </span>
                        </Button>
                        {canVote && (
                          <Button
                            variant={
                              !downvoteEnabled ? "destructive" : "outline"
                            }
                            size="sm"
                            disabled={!downvoteEnabled}
                            onClick={() => vote(false)}
                          >
                            <ArrowDown />
                          </Button>
                        )}
                      </div>
                    </div>
                  </CardFooter>
                </Card>

                {/* "add a comment" is a seperate card from the actual comments */}
                {currentUser ? (
                  <Card className="shadow-sm">
                    <CardHeader>
                      <CardTitle className="text-lg font-semibold">
                        Add a Comment
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-3">
                      <Textarea
                        placeholder="Write a comment..."
                        value={newComment}
                        onChange={(e) => setNewComment(e.target.value)}
                        className="min-h-[100px] resize-none"
                        disabled={submittingComment}
                      />
                      <div className="flex justify-end">
                        <Button
                          onClick={submitComment}
                          disabled={submittingComment || !newComment.trim()}
                          size="sm"
                        >
                          <Send className="h-4 w-4 mr-2" />
                          {submittingComment ? "Posting..." : "Post Comment"}
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                ) : (
                  <Card className="shadow-sm">
                    <CardContent className="py-6">
                      <div className="text-center text-gray-600">
                        <p>
                          <Link to="/login" className="text-blue-600 hover:underline">
                            Log in
                          </Link>{" "}
                          to post a comment
                        </p>
                      </div>
                    </CardContent>
                  </Card>
                )}

                {/* list of comments; future: limit to X comments and paginate */}
                <Card className="shadow-sm">
                  <CardHeader>
                    <CardTitle className="text-xl font-semibold">
                      Comments ({comments.length})
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      {comments.length === 0 ? (
                        <p className="text-center text-gray-500 py-8">
                          No comments yet 🙏
                        </p>
                      ) : (
                        comments.map((comment) => (
                          <div
                            key={comment.id}
                            className="border-b last:border-0 pb-4 last:pb-0"
                          >
                            <div className="flex items-start gap-3">
                              <div className="flex-shrink-0">
                                <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center">
                                  <UserIcon className="h-4 w-4 text-gray-600" />
                                </div>
                              </div>
                              <div className="flex-1 min-w-0 space-y-1">
                                <div className="flex items-center gap-2 text-sm">
                                  <Link
                                    to={`/u/${comment.author}`}
                                    className="font-medium text-gray-900 hover:text-blue-600 hover:underline"
                                  >
                                    {comment.author}
                                  </Link>
                                  <span className="text-gray-500">•</span>
                                  <Tooltip>
                                    <TooltipTrigger asChild>
                                      <span className="text-gray-500">
                                        {getTimeAgo(comment.created_at)}
                                      </span>
                                    </TooltipTrigger>
                                    <TooltipContent>
                                      <span>{datePrettyPrint(comment.created_at)}</span>
                                    </TooltipContent>
                                  </Tooltip>
                                </div>
                                <p className="text-gray-700 leading-relaxed break-words overflow-wrap-anywhere">
                                  {comment.content}
                                </p>
                                {comment.flagged && (
                                  <p className="text-xs text-red-600 mt-1">
                                    This comment has been flagged
                                  </p>
                                )}
                                
                                {/* Comment voting buttons */}
                                <div className="flex items-center gap-2 mt-2">
                                  {currentUser && (
                                    <Button
                                      variant={comment.isUpvoted ? "destructive" : "outline"}
                                      size="sm"
                                      className="h-7 px-2"
                                      disabled={comment.isUpvoted}
                                      onClick={() => voteComment(comment.id, true)}
                                    >
                                      <ArrowUp className="h-3 w-3" />
                                    </Button>
                                  )}
                                  <Button variant="outline" size="sm" className="h-7 px-3">
                                    <span className="text-xs font-medium">
                                      {comment.upvotes - comment.downvotes} upvotes
                                    </span>
                                  </Button>
                                  {currentUser && (
                                    <Button
                                      variant={comment.isDownvoted ? "destructive" : "outline"}
                                      size="sm"
                                      className="h-7 px-2"
                                      disabled={comment.isDownvoted}
                                      onClick={() => voteComment(comment.id, false)}
                                    >
                                      <ArrowDown className="h-3 w-3" />
                                    </Button>
                                  )}
                                </div>
                              </div>
                            </div>
                          </div>
                        ))
                      )}
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </SidebarInset>
        </div>
      </SidebarProvider>
    )
  );
}