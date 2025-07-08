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
import { Link, useParams } from "react-router-dom";
import { Button } from "./ui/button";
import { useEffect } from "react";

export default function Submission() {
  const { sid } = useParams();

  useEffect(() => {
    console.log(sid);
  }, [sid]);
  
  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          <div className="min-h-screen bg-gray-50 py-8 px-4">
            <div className="max-w-3xl mx-auto">
              <Card className="shadow-sm">
                <CardHeader className="space-y-1">
                  <CardTitle className="text-2xl font-bold leading-tight">
                    Test Title
                  </CardTitle>
                  <CardDescription className="flex items-center gap-2 text-base">
                    <ArrowUpRight className="h-4 w-4" />
                    <a
                      href={
                        "https://www.google.com/example/page/another/title.html"
                      }
                      className="text-blue-600 hover:text-blue-800 hover:underline break-all"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      https://www.google.com/example/page/another/title.html
                    </a>
                  </CardDescription>
                </CardHeader>

                <CardContent className="space-y-4">
                  <div className="prose prose-gray max-w-none">
                    <p className="text-gray-700 leading-relaxed">
                      This is an example of text. This is an example of text.
                      This is an example of text. This is an example of text.
                      This is an example of text. This is an example of text.
                      This is an example of text. This is an example of text.
                    </p>
                  </div>

                  <div className="flex items-center gap-4 pt-4 border-t">
                    <div className="flex items-center gap-2 text-sm text-gray-600">
                      <UserIcon className="h-4 w-4" />
                      <span>Posted by</span>
                      <Link
                        to={"/u/trent"}
                        className="font-medium text-gray-900 hover:text-blue-600 hover:underline"
                      >
                        trent
                      </Link>
                    </div>

                    <div className="flex items-center gap-2 text-sm text-gray-600">
                      <Clock className="h-4 w-4" />
                      <span>01/01/2001</span>
                    </div>
                  </div>
                </CardContent>

                <CardFooter className="bg-gray-50 border-t">
                  <div className="flex items-center justify-between w-full">
                    <div className="flex items-center gap-4">
                      <Button variant="outline" size="sm">
                        <span className="font-medium">12</span>
                        <span className="ml-1">upvotes</span>
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
  );
}
