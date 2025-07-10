"use client";
import Cookies from "js-cookie";
import { CirclePlus, Command, FileStack, ListFilterPlus, LogIn, Search, Settings, Trophy } from "lucide-react";
import {
  Sidebar,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import NavUser from "@/components/fragments/NavUser";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

type User = {
  name: string;
  email: string;
  avatar: string;
};

export default function WebSidebar() {
  const [user, setUser] = useState<User>();
  const [isAuth, setIsAuth] = useState(true);

  useEffect(() => {
    if (Cookies.get("token") == undefined) {
      setIsAuth(false);
      return;
    }

    fetch("http://localhost:3000/api/v1/me", {
      headers: {
        Authorization: "Bearer " + Cookies.get("token"),
      },
    })
      .then((res) => {
        if (!res.ok) {
          throw new Error("Network response was not ok");
        }
        return res.json();
      })
      .then((json) => {
        console.log(json);
        const u: User = {
          email: json.email,
          name: json.username,
          avatar: "/avatars/shadcn.jpg",
        };
        setUser(u);
      })
      .catch((err: Error) => {
        console.error(err);
        setIsAuth(false);
      });
  }, []);

  return (
    <Sidebar>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="/">
                <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                  <Command className="size-4" />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">HackerNews</span>
                  <span className="truncate text-xs">v0.0.1</span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup key={"Browse"}>
          <SidebarGroupLabel>{"Browse"}</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {/* LATEST BUTTON */}
              <SidebarMenuItem key={"latest"}>
                <SidebarMenuButton asChild>
                  <Link to="/">
                    <ListFilterPlus />
                    <span>Latest</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>

              {/* POPULAR/TOP BUTTON */}
              <SidebarMenuItem key={"popular"}>
                <SidebarMenuButton asChild>
                  <Link to="/top">
                    <Trophy />
                    <span>Popular</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>

              {/* SEARCH BUTTON */}
              <SidebarMenuItem key={"search"}>
                <SidebarMenuButton asChild>
                  <Link to="/search">
                    <Search />
                    <span>Search</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarGroup key={"Contribute"}>
          <SidebarGroupLabel>{"Contribute"}</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {/* NEW SUBMISSION BUTTON */}
              <SidebarMenuItem key={"new"}>
                <SidebarMenuButton asChild>
                  <Link to="/submit">
                    <CirclePlus />
                    <span>New Submission</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>

              {/* MY SUBMISSIONS BUTTON */}
              <SidebarMenuItem key={"my"}>
                <SidebarMenuButton asChild>
                  <Link to="/account/submissions">
                    <FileStack />
                    <span>My Submissions</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>

              {/* SETTINGS BUTTON */}
              <SidebarMenuItem key={"settings"}>
                <SidebarMenuButton asChild>
                  <Link to="/account/settings">
                    <Settings />
                    <span>Settings</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        {isAuth && user != undefined && (
          <NavUser avatar={user.avatar} email={user.email} name={user.name} />
        )}
        {!isAuth && (
          <SidebarMenuItem key={"login"}>
            <SidebarMenuButton asChild>
              <Link to={"/login"}>
                <LogIn />
                <span>Login</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        )}
      </SidebarFooter>
    </Sidebar>
  );
}
