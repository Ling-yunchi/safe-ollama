import * as React from "react";
import { BookOpen, Bot, ShieldHalf, SquareTerminal } from "lucide-react";

import { NavMain } from "@/components/nav-main";
import { NavUser } from "@/components/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { useAtomValue } from "jotai";
import { userAtom, userRole } from "@/storage/user";

const nav = [
  {
    title: "Home",
    url: "/",
    icon: SquareTerminal,
  },
  {
    title: "User",
    url: "/user",
    icon: Bot,
    role: "admin",
  },
  {
    title: "Token",
    url: "/token",
    icon: BookOpen,
  },
];

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const role = useAtomValue(userRole);
  const user = useAtomValue(userAtom);
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <NavUser user={user} />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={nav} role={role} />
      </SidebarContent>
      <SidebarRail />
      <SidebarFooter>
        <div className="flex items-center justify-center font-bold text-primary p-2">
          <ShieldHalf className="mr-2" /> Safe Ollama
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}
