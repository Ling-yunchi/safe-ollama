import { useAtom } from "jotai";
import { userAtom } from "../storage/user";
import { useNavigate } from "react-router";
import { useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { Label } from "@/components/ui/label";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { toast } from "sonner";

function Login() {
  const [user, setUser] = useAtom(userAtom);
  const navigate = useNavigate();

  useEffect(() => {
    if (user.userId !== -1 && new Date(user.expires) > new Date()) {
      navigate("/");
    }
  });

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <div className={cn("flex flex-col gap-6")}>
          <TooltipProvider delayDuration={0}>
            <Card>
              <CardHeader>
                <CardTitle className="text-2xl">Safe Ollama</CardTitle>
              </CardHeader>
              <CardContent>
                <form
                  onSubmit={async (e) => {
                    e.preventDefault();
                    const formData = new FormData(e.target as HTMLFormElement);
                    const data = {
                      username: formData.get("username"),
                      password: formData.get("password"),
                    };
                    const res = await fetch("/api/auth/login", {
                      method: "POST",
                      headers: {
                        "Content-Type": "application/json",
                      },
                      body: JSON.stringify(data),
                    });
                    if (res.ok) {
                      const data = await res.json();
                      setUser(data);
                      navigate("/");
                    } else {
                      const data = await res.json();
                      toast.error(data.error);
                    }
                  }}
                  className="flex flex-col gap-6"
                >
                  <div className="grid gap-2">
                    <Label htmlFor="username">Username</Label>
                    <Input id="username" name="username" required />
                  </div>
                  <div className="grid gap-2">
                    <div className="flex items-center">
                      <Label htmlFor="password">Password</Label>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <a className="ml-auto inline-block text-sm underline-offset-4 hover:underline">
                            Forgot your password?
                          </a>
                        </TooltipTrigger>
                        <TooltipContent>
                          Please contact with system administrator
                        </TooltipContent>
                      </Tooltip>
                    </div>
                    <Input
                      id="password"
                      name="password"
                      type="password"
                      required
                    />
                  </div>
                  <Button type="submit" className="w-full">
                    Login
                  </Button>
                </form>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <a className="mt-4 text-center text-sm  w-full inline-block underline-offset-4 hover:underline">
                      Don&apos;t have an account?{" "}
                    </a>
                  </TooltipTrigger>
                  <TooltipContent side="bottom">
                    Please contact with system administrator
                  </TooltipContent>
                </Tooltip>
              </CardContent>
            </Card>
          </TooltipProvider>
        </div>
      </div>
    </div>
  );
}

export default Login;
