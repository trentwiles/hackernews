import WebSidebar from "./fragments/WebSidebar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { format } from "date-fns";
import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SidebarInset, SidebarProvider } from "./ui/sidebar";
import { Controller, useForm } from "react-hook-form";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Cookies from "js-cookie";
import { Calendar as CalendarIcon } from "lucide-react";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

type FormData = {
  dateOfBirth: Date | undefined;
  fullName: string;
  body: string;
  agreeToTerms: boolean;
};

// converts the Date object provided by the shadcn calendar
// into a string format (MM-DD-YYYY) that the backend (aka postgres)
// can understand
function formatDateForBackend(date: Date | undefined): string {
  if (!date) return '';
  
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  
  return `${month}-${day}-${year}`;
}

// postgres returns an ISO string by default
// this code converts that iso string into a format that the ShadCN calendar
// can read and then display to the user
// (also provides support for other non-ISO formats for backwards compatability)
function parseBackendDate(dateStr: string | undefined): Date | undefined {
  if (!dateStr) return undefined;
  
  if (dateStr.includes('T')) {
    const datePart = dateStr.split('T')[0];
    const [year, month, day] = datePart.split('-').map(n => parseInt(n, 10));
    
    return new Date(year, month - 1, day, 12, 0, 0);
  }
  
  if (dateStr.match(/^\d{4}-\d{2}-\d{2}$/)) {
    const [year, month, day] = dateStr.split('-').map(n => parseInt(n, 10));
    return new Date(year, month - 1, day, 12, 0, 0);
  }
  
  if (dateStr.match(/^\d{2}-\d{2}-\d{4}$/)) {
    const [month, day, year] = dateStr.split('-').map(n => parseInt(n, 10));
    return new Date(year, month - 1, day, 12, 0, 0);
  }
  
  console.error("Unable to parse date:", dateStr);
  return undefined;
}

export default function Settings() {
  const [reloads, setReloads] = useState<number>(0);
  const [submitButtonText, setSubmitButtonText] = useState<string>("Submit")

  const {
    register,
    control,
    handleSubmit,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    defaultValues: {
      dateOfBirth: undefined,
      fullName: "",
      body: "",
      agreeToTerms: false, // legacy as I copied from Submission.tsx, should remove at some point
    },
  });

  const onSubmit = async (data) => {
    console.log("Form submitted:", data);

    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/bio", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + Cookies.get("token"),
      },
      body: JSON.stringify({
        fullName: data.fullName,
        birthdate: formatDateForBackend(data.dateOfBirth),
        bioText: data.body,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to fetch user data");
        }
        return response.json();
      })
      .then(() => {
        console.log("success")
        setReloads((prev) => prev + 1)
      })
      .catch((err) => {
        console.error("Fetch error:", err);
      });
  };

  useEffect(() => {
  fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/me", {
    headers: {
      Authorization: "Bearer " + Cookies.get("token"),
    },
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Failed to fetch user data");
      }
      return response.json();
    })
    .then((data) => {
      setValue("fullName", data.metadata.full_name);
      setValue("body", data.metadata.bio);
      
      // Parse the date properly
      if (data.metadata.birthday) {
        setValue("dateOfBirth", parseBackendDate(data.metadata.birthday));
      }

      setSubmitButtonText("Saved!")
      setTimeout(() => {
        setSubmitButtonText("Submit")
      }, 1000)
    })
    .catch((err) => {
      console.error("Fetch error:", err);
      setSubmitButtonText("Error Saving")
    });
}, [reloads]);

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          {/* BEGIN FORM */}
          <div className="max-w-5xl mx-auto p-4 w-full">
            <Card>
              <CardHeader>
                <CardTitle>User Settings</CardTitle>
                {/* <CardDescription>
                  Share a link with the community
                </CardDescription> */}
              </CardHeader>

              <CardContent className="space-y-6">
                {/* Calendar Field */}
                <div className="space-y-2">
                  <Label htmlFor="dateOfBirth">Date of Birth</Label>
                  <Controller
                    name="dateOfBirth"
                    control={control}
                    rules={{
                      required: "Date of birth is required",
                      validate: (value) => {
                        if (!value) return "Please select a date";
                        const selectedDate = new Date(value);
                        const today = new Date();
                        if (selectedDate > today) {
                          return "Date of birth cannot be in the future";
                        }
                        const age =
                          today.getFullYear() - selectedDate.getFullYear();
                        if (age > 120) {
                          return "Please enter a valid date of birth";
                        }
                        return true;
                      },
                    }}
                    render={({ field }) => (
                      <Popover>
                        <PopoverTrigger asChild>
                          <Button
                            variant="outline"
                            className={cn(
                              "w-full justify-start text-left font-normal",
                              !field.value && "text-muted-foreground"
                            )}
                          >
                            <CalendarIcon className="mr-2 h-4 w-4" />
                            {field.value ? (
                              format(field.value, "MM/dd/yyyy")
                            ) : (
                              <span>Pick a date</span>
                            )}
                          </Button>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-0" align="start">
                          <Calendar
                            mode="single"
                            selected={field.value}
                            onSelect={field.onChange}
                            disabled={(date) =>
                              date > new Date() || date < new Date("1900-01-01")
                            }
                            initialFocus
                          />
                        </PopoverContent>
                      </Popover>
                    )}
                  />
                  {errors.dateOfBirth && (
                    <p className="text-sm text-red-500">
                      {errors.dateOfBirth.message}
                    </p>
                  )}
                </div>

                {/* Full Name Field */}
                <div className="space-y-2">
                  <Label htmlFor="title">Full Name</Label>
                  <Input
                    id="fullName"
                    placeholder="Enter your full name"
                    {...register("fullName", {
                      required: "Full name is required",
                      minLength: {
                        value: 2,
                        message: "Name must be at least 2 characters",
                      },
                    })}
                  />
                  {errors.fullName && (
                    <p className="text-sm text-red-500">
                      {errors.fullName.message}
                    </p>
                  )}
                </div>

                {/* Body Field */}
                <div className="space-y-2">
                  <Label htmlFor="body">Bio</Label>
                  <Textarea
                    id="body"
                    placeholder="I'm a 25 year old developer from San Francisco..."
                    className="min-h-[120px]"
                    {...register("body", {
                      maxLength: {
                        value: 1000,
                        message: "Bio must be less than 1000 characters",
                      },
                    })}
                  />
                  {errors.body && (
                    <p className="text-sm text-red-500">
                      {errors.body.message}
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
                  disabled={isSubmitting || submitButtonText == "Saved!"}
                >
                  {submitButtonText}
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
