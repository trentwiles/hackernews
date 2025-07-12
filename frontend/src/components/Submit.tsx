import WebSidebar from "./fragments/WebSidebar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SidebarInset, SidebarProvider } from "./ui/sidebar";
import { useForm } from "react-hook-form";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { useGoogleReCaptcha } from "react-google-recaptcha-v3";

export default function Submit() {
  console.log(Cookies.get("token"))
  const { executeRecaptcha } = useGoogleReCaptcha();
  const navigate = useNavigate();

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm({
    defaultValues: {
      link: "",
      title: "",
      body: "",
      agreeToTerms: false,
    },
  });

  const linkValue = watch("link");

  useEffect(() => {
    if (linkValue == "") {
      return;
    }
    const timeoutId = setTimeout(() => {
      fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/fetchWebsiteTitle?url=" + linkValue) // Replace with your actual endpoint
        .then((res) => {
          if (res.status !== 200) {
            throw new Error("non-200 status: " + res.status)
          }
          return res.json()
        })
        .then((json) => {
          // assuming the JSON response is 200 OK, we should have a response in the form:
          // {"title":"..."}
          setValue("title", json.title)
        })
        .catch((err) => console.error("Error fetching data:", err));
    }, 500);
    // ^^^ 500 milliseconds pause between when they stop typing and when we make the HTTP request

    return () => clearTimeout(timeoutId);
  }, [linkValue]);

  type formData = {
    link: string;
    title: string;
    body: string;
    agreeToTerms: boolean;
  };

  type submitSuccessResponse = {
    id: string;
  };

  const onSubmit = async (data: formData) => {
    console.log("Form submitted:", data);

    if (!executeRecaptcha) {
      console.log("Execute recaptcha not yet available");
      return;
    }

    const token = await executeRecaptcha("login");

    // req.CaptchaToken == "" || req.Link == "" || req.Title == ""

    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + Cookies.get("token"),
      },
      body: JSON.stringify({
        CaptchaToken: token,
        Link: data.link,
        Title: data.title,
        Body: data.body,
      }),
    })
      .then((res) => res.json())
      .then((data: submitSuccessResponse) => {
        navigate("/submission/" + data.id + "?new=true");
      })
      .catch((error) => {
        console.error(error);
      });
  };

  const [isChecked, setIsChecked] = useState<boolean>(false);

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          {/* BEGIN FORM */}
          <div className="max-w-5xl mx-auto p-4 w-full">
            <Card>
              <CardHeader>
                <CardTitle>Submit Link</CardTitle>
                <CardDescription>
                  Share a link with the community
                </CardDescription>
              </CardHeader>

              <CardContent className="space-y-6">
                {/* Link Field */}
                <div className="space-y-2">
                  <Label htmlFor="link">Link *</Label>
                  <Input
                    id="link"
                    type="url"
                    placeholder="https://example.com"
                    {...register("link", {
                      required: "Link is required",
                      pattern: {
                        value:
                          /^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_.~#?&//=]*)$/,
                        message: "Please enter a valid URL",
                      },
                    })}
                  />
                  {errors.link && (
                    <p className="text-sm text-red-500">
                      {errors.link.message}
                    </p>
                  )}
                </div>

                {/* Title Field */}
                <div className="space-y-2">
                  <Label htmlFor="title">Title *</Label>
                  <Input
                    id="title"
                    placeholder="Enter a descriptive title"
                    {...register("title", {
                      required: "Title is required",
                      minLength: {
                        value: 3,
                        message: "Title must be at least 3 characters",
                      },
                      maxLength: {
                        value: 200,
                        message: "Title must be less than 200 characters",
                      },
                    })}
                  />
                  {errors.title && (
                    <p className="text-sm text-red-500">
                      {errors.title.message}
                    </p>
                  )}
                </div>

                {/* Body Field */}
                <div className="space-y-2">
                  <Label htmlFor="body">Body (Optional)</Label>
                  <Textarea
                    id="body"
                    placeholder="Add any additional context or description..."
                    className="min-h-[120px]"
                    {...register("body", {
                      maxLength: {
                        value: 1000,
                        message: "Body must be less than 1000 characters",
                      },
                    })}
                  />
                  {errors.body && (
                    <p className="text-sm text-red-500">
                      {errors.body.message}
                    </p>
                  )}
                </div>

                {/* Terms Checkbox */}
                <div className="space-y-2">
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="agreeToTerms"
                      onCheckedChange={(checked: boolean) => {
                        setValue("agreeToTerms", checked);
                        setIsChecked(checked);
                      }}
                      {...register("agreeToTerms", {
                        required: "You must agree to the terms of service",
                      })}
                    />
                    <Label htmlFor="agreeToTerms" className="font-normal">
                      I agree to the terms of service *
                    </Label>
                  </div>
                  {errors.agreeToTerms && (
                    <p className="text-sm text-red-500">
                      {errors.agreeToTerms.message}
                    </p>
                  )}
                </div>
              </CardContent>

              <CardFooter className="flex justify-between">
                <Link to="/">
                  <Button type="button" variant="outline">
                    Cancel
                  </Button>
                </Link>
                <Button
                  onClick={handleSubmit(onSubmit)}
                  disabled={isSubmitting || !isChecked}
                >
                  {isSubmitting ? "Submitting..." : "Submit"}
                </Button>
              </CardFooter>
            </Card>
          </div>
          {/* END FORM */}
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
