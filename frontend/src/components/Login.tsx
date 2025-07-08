import { GalleryVerticalEnd } from "lucide-react";
import Cookies from "js-cookie";
import { LoginForm } from "@/components/login-form";
import { useNavigate } from "react-router-dom";
import { useEffect } from "react";

type LoginProps = {
  serviceName: string;
};

export default function Login(props: LoginProps) {
  const navigate = useNavigate();

  useEffect(() => {
    if (Cookies.get("token") != undefined) {
      navigate("/");
      return;
    }
  }, []);

  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="/" className="flex items-center gap-2 self-center font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <GalleryVerticalEnd className="size-4" />
          </div>
          {props.serviceName}
        </a>
        <LoginForm />
      </div>
    </div>
  );
}
