import { GalleryVerticalEnd } from "lucide-react";
import { useLocation } from "react-router-dom";

import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function LoginThanks() {
  const location = useLocation();

  const searchParams = new URLSearchParams(location.search);
  let email = searchParams.get("email"); // "123123"

  if (email == "" || email == null) {
    email = "your email";
  }

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
              <CardTitle>Check Your Email</CardTitle>

              <CardDescription>
                Check {email} for a magic link to log in. Not seeing anything?
                Check your spam.
              </CardDescription>
            </CardHeader>
          </Card>
        </div>
      </div>
    </div>
  );
}
