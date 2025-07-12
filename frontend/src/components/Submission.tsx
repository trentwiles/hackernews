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

export default function Submission() {
  const { sid } = useParams();
  const navigate = useNavigate();
  const [s, setS] = useState<submission>();
  const [pending, setPending] = useState<boolean>(true);
  const [error, setError] = useState<boolean>(false);
  const [currentUser, setCurrentUser] = useState<string>();

  const [upvoteEnabled, setUpvoteEnabled] = useState<boolean>(true);
  const [downvoteEnabled, setDownvoteEnabled] = useState<boolean>(true);
  const [canVote, setCanVote] = useState<boolean>(false);

  const [shareButtonText, setShareButtonText] = useState<string>("Share")

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
          throw new Error("Unauthorized");  // This will skip to catch
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

  console.log(canVote + " <-- can vote???")

  return (
    !pending &&
    !error &&
    s !== undefined && (
      <SidebarProvider>
        <div className="flex min-h-screen w-full">
          <WebSidebar />
          <SidebarInset>
            <div className="min-h-screen bg-gray-50 py-8 px-4">
              <div className="max-w-3xl mx-auto">
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
                        <Button variant="outline" size="sm" onClick={async () => {
                          await navigator.clipboard.writeText(window.location.href);
                          setShareButtonText("Copied To Clipboard")
                          await new Promise((r) => setTimeout(r, 2000))
                          setShareButtonText("Share")
                        }}>
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
                              variant={
                                !upvoteEnabled ? "destructive" : "outline"
                              }
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
              </div>
            </div>
          </SidebarInset>
        </div>
      </SidebarProvider>
    )
  );
}
