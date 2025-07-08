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
import { CalendarDays, User as UserIcon } from "lucide-react";
import CryptoJS from "crypto-js";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@radix-ui/react-tooltip";

type user = {
  email: string;
  joined: string;
  bio: string;
  birthday: string;
  full_name: string;
  username: string;
};

type basicSubmission = {
  Id: string;
  Title: string;
  Link: string;
  Created_at: string;
};

function datePrettyPrint(isoDate: string) {
  const date = new Date(isoDate);

  const day = String(date.getUTCDate()).padStart(2, "0");
  const month = String(date.getUTCMonth() + 1).padStart(2, "0");
  const year = date.getUTCFullYear();

  let hours = date.getUTCHours();
  const minutes = String(date.getUTCMinutes()).padStart(2, "0");
  const ampm = hours >= 12 ? "PM" : "AM";
  hours = hours % 12 || 12;
  const hourStr = String(hours).padStart(2, "0");

  const formatted = `${month}/${day}/${year}, ${hourStr}:${minutes} ${ampm}`;

  return formatted;
}

function getTimeAgo(dateString: string): string {
  const now = new Date();
  const past = new Date(dateString);
  const diffMs = now.getTime() - past.getTime();

  const msPerHour = 1000 * 60 * 60;
  const msPerDay = msPerHour * 24;
  const msPerMonth = msPerDay * 30.44; // average
  const msPerYear = msPerDay * 365.25; // average

  if (diffMs < msPerDay) {
    const hours = Math.floor(diffMs / msPerHour);
    return `${hours} hour${hours !== 1 ? "s" : ""} ago`;
  } else if (diffMs < msPerMonth) {
    const days = Math.floor(diffMs / msPerDay);
    return `${days} day${days !== 1 ? "s" : ""} ago`;
  } else if (diffMs < msPerYear) {
    const months = Math.floor(diffMs / msPerMonth);
    return `${months} month${months !== 1 ? "s" : ""} ago`;
  } else {
    const years = Math.floor(diffMs / msPerYear);
    return `${years} year${years !== 1 ? "s" : ""} ago`;
  }
}

export default function User() {
  const navigate = useNavigate();
  const { username } = useParams();
  const [pending, setPending] = useState<boolean>(true);
  const [error, isError] = useState<boolean>(false);
  const [currentUser, setCurrentUser] = useState<user>();
  const [submissions, setSubmissions] = useState<basicSubmission[]>();

  useEffect(() => {
    fetch("http://localhost:3000/api/v1/user?username=" + username)
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
                    </CardContent>
                  </Card>
                </div>

                {/* Right Column - Submissions (2/3) */}
                <div className="md:col-span-2">
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
                </div>
              </div>
            </div>
          </SidebarInset>
        </div>
      </SidebarProvider>
    )
  );
}
