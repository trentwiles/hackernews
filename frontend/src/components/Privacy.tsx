import WebSidebar from "./fragments/WebSidebar";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Button } from "./ui/button";
import { useState } from "react";
import Cookies from "js-cookie";

export default function Privacy() {
  const [buttonText, setButtonText] = useState<string>("Export Your Data");
  const [buttonEnabled, setButtonEnabled] = useState<boolean>(true);
  const [exportComplete, setExportComplete] = useState<boolean>(false);

  function exportData() {
    setButtonText("Processing Request...");
    setButtonEnabled(false);

    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/dump", {
      headers: {
        Authorization: "Bearer " + Cookies.get("token"),
      },
      method: "POST",
    }) // replace with your API URL
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then(() => {
        setButtonText("Request Submitted");
        setExportComplete(true);
      })
      .catch((err) => {
        console.error(err);
        setButtonText("Issue Sending Request.. Try Again Later");
      });
  }

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          {/* BEGIN FORM */}
          <div className="max-w-5xl mx-auto p-4 w-full">
            <Card>
              <CardHeader>
                <CardTitle>Privacy Settings</CardTitle>
              </CardHeader>

              <CardContent className="space-y-6">
                {/* Calendar Field */}
                <div className="space-y-2">
                  <Button
                    onClick={() => exportData()}
                    disabled={!buttonEnabled}
                  >
                    {buttonText}
                  </Button>
                </div>
                <p>
                  You'll be emailed a copy of the your data, including account
                  details, submissions, comments, and voting history.
                </p>
                <p>
                  Depending on how much information you have, this process could
                  take a few minutes to a few hours.
                </p>
                {exportComplete && (
                  <>
                    <p>
                      <span style={{ color: `green` }}>Export complete.</span>
                      <a href={import.meta.env.VITE_API_ENDPOINT + "/api/v1/dump?authToken=" + Cookies.get("token")} target="_blank">&nbsp;Click here to download.</a>
                    </p>
                  </>
                )}
              </CardContent>
            </Card>
          </div>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
