import { Outlet, useNavigate } from "react-router";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { useEffect } from "react";
import { useAtomValue } from "jotai";
import { userToken } from "@/storage/user";

function Dashboard() {
  const navigate = useNavigate();
  const token = useAtomValue(userToken);
  useEffect(() => {
    validateToken();
  }, []);
  const validateToken = async () => {
    const res = await fetch("/api/auth/validateToken", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Token: token,
      },
    });
    if (!res.ok) {
      navigate("/login");
    }
  };
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <div className="flex flex-1">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}

export default Dashboard;
