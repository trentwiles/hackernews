import { cn } from "@/lib/utils";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { useGoogleReCaptcha } from "react-google-recaptcha-v3";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const DEFAULT_LOGIN_TEXT = "Log In";

  const { executeRecaptcha } = useGoogleReCaptcha();
  const [username, setUsername] = useState<string>();
  const [email, setEmail] = useState<string>();
  const [canSubmit, setCanSubmit] = useState<boolean>(true);
  const [loginText, setLoginText] = useState<string>(DEFAULT_LOGIN_TEXT);

  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault(); // prevent page reload

    if (!executeRecaptcha) {
      console.log("Execute recaptcha not yet available");
      return;
    }

    const currentUsername = username;
    const currentEmail = email;

    setUsername("");
    setEmail("");

    setCanSubmit(false);

    setLoginText("Please wait...");

    const token = await executeRecaptcha("login");

    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: currentUsername,
        email: currentEmail,
        captchaToken: token,
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
        setCanSubmit(true);

        navigate("/login-thanks?email=" + email);
        return;
      })
      .catch((error) => {
        console.log(error);
        setCanSubmit(true);
      });

    setLoginText("Log In");
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome back</CardTitle>
          <CardDescription>
            Enter your username and email to recieve a magic link.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div>
              <div className="grid gap-6">
                <div className="grid gap-3">
                  <div className="flex items-center">
                    <Label htmlFor="password">Username</Label>
                    {/* When you click here, include a popup that tells you how to
          make an accont */}
                    <a
                      href="#"
                      className="ml-auto text-sm underline-offset-4 hover:underline"
                    >
                      No account?
                    </a>
                  </div>
                  <Input
                    id="username"
                    type="text"
                    placeholder="Alice"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    pattern="^\S+$"
                    title="Username cannot contain spaces"
                    required
                  />
                </div>
                <div className="grid gap-3">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="alice@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                  />
                </div>
                <Button type="submit" className="w-full" disabled={!canSubmit}>
                  {loginText}
                </Button>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
      <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
        By authenticating, you agree to our <a href="#">Terms of Service</a> and{" "}
        <a href="#">Privacy Policy</a>.
      </div>
    </div>
  );
}
