import WebSidebar, { SidebarBreadcrumbHeader } from "./fragments/WebSidebar";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Button } from "./ui/button";
import { useEffect, useState } from "react";
import Cookies from "js-cookie";
import { useSearchParams } from "react-router-dom";

export default function UrlCheck() {
  const [buttonText, setButtonText] = useState<string>("Continue");
  const [buttonEnabled, setButtonEnabled] = useState<boolean>(true);
  const [searchParams] = useSearchParams();

  const string: string = searchParams.get("q") || "/";
  useEffect(() => {
    if (string === "/") {
      setButtonEnabled(false);
      setButtonText("Invalid URL");
      return;
    }
    setButtonText("Please Wait...");
    setButtonEnabled(false);

    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/urlCheck?q=" + string, {
      headers: {
        Authorization: "Bearer " + Cookies.get("token"),
      },
      method: "GET",
    }) // replace with your API URL
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((j) => {
        if (j.passed) {
          window.location.href = string;
          return;
        } else {
          setButtonText("Malicious Link Detected");
          setButtonEnabled(false);
        }
      })
      .catch((err) => {
        console.error(err);
        setButtonText("Malicious Link Detected");
        setButtonEnabled(false);
      });
  }, [string, searchParams]);

  const breadcrumbs = [{ label: "URL Redirect", isCurrentPage: true }];

  return (
    <SidebarProvider>
      <WebSidebar />
      <SidebarInset>
        <SidebarBreadcrumbHeader breadcrumbs={breadcrumbs} />
        <div className="flex flex-1 flex-col gap-4 p-4">
          {/* BEGIN FORM */}
          <div className="max-w-5xl mx-auto p-4 w-full">
            <Card>
              <CardHeader>
                <CardTitle>Attention</CardTitle>
              </CardHeader>

              <CardContent className="space-y-6">
                {/* Calendar Field */}
                <div className="space-y-2">
                  <Button disabled={!buttonEnabled}>{buttonText}</Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
