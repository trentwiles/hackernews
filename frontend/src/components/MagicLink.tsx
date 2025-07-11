import { useNavigate } from 'react-router-dom';
import { Button } from "@/components/ui/button";
import { GalleryVerticalEnd } from "lucide-react";
import { useLocation } from "react-router-dom";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useEffect, useState } from "react";

type MagicLinkProps = {
  serviceName: string;
};

export default function MagicLink(props: MagicLinkProps) {
  const [username, setUsername] = useState<string>("");
  const [canConfirm, setCanConfirm] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);
  const [errorMessage, setErrorMessage] = useState<string>("");

  const location = useLocation();
  const navigate = useNavigate();

  const searchParams = new URLSearchParams(location.search);
  const token = searchParams.get("token"); // "123123"

  useEffect(() => {
    if (!token || token == "" || token == null) {
      setIsError(true);
      setErrorMessage("No Token Provided");
      return
    }

    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/magic?token=" + token)
      .then((response) => {
        // Check if the response is ok (status 200-299)
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        return response.json();
      })
      .then((data) => {
        // if we have made it to this point, we know there has been a 200 OK
        // which means there will be a valid token!
        const token = data.token
        setUsername(data.username)

        document.cookie = `token=${token}; path=/; max-age=3600`;
        // DANGER: this username cookie can obviously be edited by a client and forged
        //         it should never be used to validate that someone is logged in, the token
        //         above should always be passed to the server to preform this job
        document.cookie = `username=${data.username}; path=/; max-age=3600`;
        setIsError(false)
        setCanConfirm(true)
      })
      .catch((error) => {
        // Handle errors
        console.error("Login error:", error);
        setIsError(true);
        setErrorMessage("Invalid Magic Link")
      });
  }, [token]);

  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <GalleryVerticalEnd className="size-4" />
          </div>
          {props.serviceName}
        </a>

        <div className="flex flex-col gap-6">
          <Card>
            <CardHeader className="text-center">
              {username && (
                <CardTitle className="text-xl">Hey {username}!</CardTitle>
              )}

              <CardDescription>
                {(!username && !isError) && <p>Validating your magic link...</p>}
                {isError && (
                  <span style={{ color: "red" }}>{errorMessage}</span>
                )}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button className="w-full" disabled={!canConfirm} onClick={() => navigate('/')}>
                Continue
              </Button>
            </CardContent>
          </Card>
          <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
            By authenticating, you agree to our <a href="#">Terms of Service</a>{" "}
            and <a href="#">Privacy Policy</a>.
          </div>
        </div>
      </div>
    </div>
  );
}
