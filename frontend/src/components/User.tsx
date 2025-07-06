import { useParams } from "react-router-dom";
import WebSidebar from "./fragments/WebSidebar";
import { SidebarInset, SidebarProvider } from "./ui/sidebar";
import { useEffect, useState } from "react";

type user = {
  email: string;
  joined: string;
  bio: string;
  birthday: string;
  full_name: string;
  username: string;
};

export default function User() {
  const { username } = useParams();
  const [pending, setPending] = useState<boolean>(true);
  const [currentUser, setCurrentUser] = useState<user>();

  useEffect(() => {
    fetch("http://localhost:3000/api/v1/user?username=" + username)
      .then((res) => {
        if (!res.ok) throw new Error("Network response was not ok");
        return res.json();
      })
      .then((data) => {
        const current: user = {
            bio: data.metadata.bio,
            birthday: data.metadata.birthday,
            email: data.email,
            full_name: data.metadata.full_name,
            username: data.username,
            joined: data.joined
        }

        setCurrentUser(current)
        setPending(false)
      })
      .catch((err) => {
        console.error("Fetch error:", err);
        setPending(false)
      });
  }, [username]);

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          <span>404</span>
          {!pending && <p>COMPLETE THIS LATER</p>}
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
