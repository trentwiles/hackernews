import Cookies from "js-cookie";
import { GalleryVerticalEnd } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./ui/card";
import { useNavigate } from "react-router-dom";
import { Button } from "./ui/button";

export default function Logout() {
  const navigate = useNavigate();
  Cookies.remove("token");
  Cookies.remove("username");

  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <GalleryVerticalEnd className="size-4" />
          </div>
        </a>

        <div className="flex flex-col gap-6">
          <Card>
            <CardHeader className="text-center">
              <CardTitle>Seeyaaaa</CardTitle>

              <CardDescription>
                You've been logged out.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button
                className="w-full"
                onClick={() => navigate("/")}
              >
                Return Home
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
