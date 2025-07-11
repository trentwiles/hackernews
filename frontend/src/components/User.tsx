import { Link, useNavigate, useParams } from "react-router-dom";
import WebSidebar from "./fragments/WebSidebar";
import { SidebarInset, SidebarProvider } from "./ui/sidebar";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@radix-ui/react-avatar";
import { Cake, CalendarDays, Sparkle, User as UserIcon } from "lucide-react";
import CryptoJS from "crypto-js";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@radix-ui/react-tooltip";
import { datePrettyPrint, getTimeAgo, simpleDatePrettyPrint } from "@/utils";

type user = {
  email: string;
  joined: string;
  bio: string;
  birthday: string;
  full_name: string;
  username: string;
  score: number;
};

type basicSubmission = {
  Id: string;
  Title: string;
  Link: string;
  Created_at: string;
};

type votedPost = {
  Id: string;
  Title: string;
  Link: string;
  Created_at: string;
  Username: string;
  Upvoted: boolean;
};

export default function User() {
  const navigate = useNavigate();
  const { username } = useParams();
  const [pending, setPending] = useState<boolean>(true);
  const [error, isError] = useState<boolean>(false);
  const [currentUser, setCurrentUser] = useState<user>();
  const [submissions, setSubmissions] = useState<basicSubmission[]>();
  const [votedPosts, setVotedPosts] = useState<votedPost[]>();

  useEffect(() => {
    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/user?username=" + username)
      .then((res) => {
        if (res.status === 404) {
          navigate("/404");
          return;
        }
        if (!res.ok) throw new Error("Network response was not ok");
        return res.json();
      })
      .then((data) => {
        const current: user = {
          bio: data.metadata.bio,
          birthday: data.metadata.birthday,
          email: data.email,
          full_name: data.metadata.full_name,
          username: data.username,
          joined: data.joined,
          score: data.metadata.score
        };

        setCurrentUser(current);
        setPending(false);
      })
      .catch((err) => {
        toast("Error: " + err);
        setPending(false);
        isError(true);
      });
  }, [username]);

  useEffect(() => {
    setPending(true);
    fetch("http://localhost:3000/api/v1/userSubmissions?username=" + username)
      .then((res) => {
        if (res.status === 404) {
          navigate("/404");
          return;
        }
        if (!res.ok) throw new Error("Network response was not ok");
        return res.json();
      })
      .then((data) => {
        setSubmissions(data.results);
        setPending(false);
      })
      .catch((err) => {
        toast("Error: " + err);
        setPending(false);
        isError(true);
      });
  }, [username]);

  useEffect(() => {
    setPending(true);
    fetch("http://localhost:3000/api/v1/allUserVotes?username=" + username)
      .then((res) => {
        if (res.status === 404) {
          navigate("/404");
          return;
        }
        if (!res.ok) throw new Error("Network response was not ok");
        return res.json();
      })
      .then((data) => {
        const listOfVotes: votedPost[] = [];

        if (data.results === undefined || data.results === null) {
          setVotedPosts([]);
          setPending(false);
          isError(false);
          return;
        }

        data.results.map((datapoint) => {
          const current: votedPost = {
            Created_at: datapoint.Created_at,
            Id: datapoint.Id,
            Link: datapoint.Link,
            Title: datapoint.Title,
            Upvoted: datapoint.IsUpvoted,
            Username: datapoint.Username,
          };
          listOfVotes.push(current);
        });
        setVotedPosts(listOfVotes);
        setPending(false);
      })
      .catch((err) => {
        toast("Error: " + err);
        setPending(false);
        isError(true);
      });
  }, [username]);

  return (
    !pending &&
    !error &&
    currentUser != undefined && (
      <SidebarProvider>
        <div className="flex min-h-screen w-full">
          <WebSidebar />
          <SidebarInset>
            <div className="container mx-auto p-6">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Left Column - Profile Info (1/3) */}
                <div className="md:col-span-1">
                  <Card>
                    <CardHeader className="text-center">
                      <Avatar className="w-32 h-32 mx-auto mb-4">
                        <AvatarImage
                          src={
                            "https://www.gravatar.com/avatar/" +
                            CryptoJS.MD5(currentUser.username).toString() +
                            "?d=identicon"
                          }
                          alt={currentUser.username}
                        />
                        <AvatarFallback>
                          {currentUser.username.slice(0, 2).toUpperCase()}
                        </AvatarFallback>
                      </Avatar>
                      <CardTitle className="text-2xl">
                        @{currentUser.username}
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      {currentUser.full_name != "" && (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <UserIcon className="w-4 h-4" />
                          <span>{currentUser.full_name}</span>
                        </div>
                      )}
                      <div className="text-sm">
                        <p className="text-foreground">{currentUser.bio}</p>
                      </div>

                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <CalendarDays className="w-4 h-4" />
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <span>
                              Joined {datePrettyPrint(currentUser.joined)}
                            </span>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>{getTimeAgo(currentUser.joined)}</p>
                          </TooltipContent>
                        </Tooltip>
                      </div>
                      {currentUser.birthday != "" && (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <Cake className="w-4 h-4" />
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <span>
                                {simpleDatePrettyPrint(currentUser.birthday)}
                              </span>
                            </TooltipTrigger>
                            <TooltipContent>
                              <p>{getTimeAgo(currentUser.birthday)}</p>
                            </TooltipContent>
                          </Tooltip>
                        </div>
                      )}
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <Sparkle className="w-4 h-4" />
                          {currentUser.score} karma
                        </div>
                    </CardContent>
                  </Card>
                </div>

                {/* Right Column - Submissions (2/3) */}
                <div className="md:col-span-2 space-y-6">
                  <Card>
                    <CardHeader>
                      <CardTitle>Submissions</CardTitle>
                      <CardDescription>
                        All posts and contributions
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-4">
                        {submissions !== undefined &&
                          submissions.length == 0 && (
                            <p className="text-muted-foreground text-sm">
                              It's a bit empty in here, don't you think?
                            </p>
                          )}
                        {submissions !== undefined &&
                          submissions.map((submission) => (
                            <div
                              key={submission.Id}
                              className="border rounded-lg p-4 hover:bg-accent transition-colors cursor-pointer"
                            >
                              <div className="flex justify-between items-start">
                                <div className="space-y-1">
                                  <h3 className="font-semibold text-lg hover:underline">
                                    <Link to={"/submission/" + submission.Id}>
                                      {submission.Title}
                                    </Link>
                                  </h3>
                                  <div className="flex gap-3 text-sm text-muted-foreground">
                                    <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-secondary">
                                      <a href={submission.Link} target="_blank">
                                        {submission.Link}
                                      </a>
                                    </span>
                                    <span>
                                      {datePrettyPrint(submission.Created_at)}
                                    </span>
                                  </div>
                                </div>
                              </div>
                            </div>
                          ))}
                      </div>
                    </CardContent>
                  </Card>

                  {/* Voted Posts Card */}
                  <Card>
                    <CardHeader>
                      <CardTitle>Voted Posts</CardTitle>
                      <CardDescription>
                        Posts upvoted or downvoted by {currentUser.username}
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      {!pending && !error && (
                        <div className="space-y-4">
                          {votedPosts !== undefined &&
                            votedPosts.length == 0 && (
                              <p className="text-muted-foreground text-sm">
                                No voted posts yet.
                              </p>
                            )}
                          {votedPosts !== undefined &&
                            votedPosts.map((post) => (
                              <div
                                key={post.Id}
                                className="border rounded-lg p-4 hover:bg-accent transition-colors cursor-pointer"
                              >
                                <div className="flex justify-between items-start">
                                  <div className="space-y-1">
                                    <h3 className="font-semibold text-lg hover:underline">
                                      <Link to={"/submission/" + post.Id}>
                                        {post.Title}
                                      </Link>
                                    </h3>
                                    <div className="flex gap-3 text-sm text-muted-foreground">
                                      <span className="text-xs">
                                        {currentUser.username}{" "}
                                        {post.Upvoted ? "upvoted" : "downvoted"}{" "}
                                        this
                                      </span>
                                      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-secondary">
                                        <a href={post.Link} target="_blank">
                                          {post.Link}
                                        </a>
                                      </span>
                                      <span>
                                        {datePrettyPrint(post.Created_at)}
                                      </span>
                                    </div>
                                  </div>
                                </div>
                              </div>
                            ))}
                        </div>
                      )}
                    </CardContent>
                  </Card>
                </div>
              </div>
            </div>
          </SidebarInset>
        </div>
      </SidebarProvider>
    )
  );
}
