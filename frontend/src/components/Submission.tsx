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
import { ArrowUpRight, Clock, User as UserIcon } from "lucide-react";
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

  useEffect(() => {
    if (sid === undefined || sid == "") {
      setPending(false);
      setError(true);
      return;
    }

    fetch("http://localhost:3000/api/v1/submission?id=" + sid)
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
  }, [sid]);

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
                        <Button variant="outline" size="sm">
                          <span className="font-medium">
                            {s.totalScore} upvotes
                          </span>
                        </Button>
                      </div>

                      <div className="flex items-center gap-2">
                        <Button variant="ghost" size="sm">
                          Share
                        </Button>
                        <Button variant="ghost" size="sm">
                          Save
                        </Button>
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
