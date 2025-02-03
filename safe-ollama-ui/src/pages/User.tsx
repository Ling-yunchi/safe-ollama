import { Button } from "@/components/ui/button";
import { DialogHeader, DialogFooter } from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogTitle,
  DialogClose,
  DialogDescription,
} from "@/components/ui/dialog";
import { useAtomValue } from "jotai";
import { Edit, Trash } from "lucide-react";
import { useState, useEffect } from "react";
import { toast } from "sonner";
import { userAtom, userToken } from "@/storage/user";
import { cn } from "@/lib/utils";
import { useNavigate } from "react-router";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface User {
  id: number;
  username: string;
  role: string;
}

interface UserBean {
  id: number;
  username: string;
  role: string;
}

function User() {
  const navigate = useNavigate();
  const user = useAtomValue(userAtom);
  const token = useAtomValue(userToken);
  const [users, setUsers] = useState<User[]>([]);

  useEffect(() => {
    if (user.role !== "admin") {
      navigate("/");
    }
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    const response = await fetch("/api/user/", {
      headers: {
        "Content-Type": "application/json",
        Token: token,
      },
    });
    if (!response.ok) {
      console.log("Error fetching users");
      return;
    }
    const data = await response.json();
    console.log(data);
    setUsers(data);
  };

  const handleCreateUser = async (data: UserBean) => {
    const response = await fetch("/api/user/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Token: token,
      },
      body: JSON.stringify(data),
    });
    if (response.ok) {
      toast.success("user created");
      await fetchUsers();
    }
  };

  const handleUpdateUser = async (data: UserBean) => {
    const response = await fetch(`/api/user/${data.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Token: token,
      },
      body: JSON.stringify(data),
    });
    if (response.ok) {
      toast.success("user updated");
      await fetchUsers();
    }
  };

  const handleDeleteUser = async (userId: number) => {
    const response = await fetch(`/api/user/${userId}`, {
      method: "DELETE",
      headers: {
        Token: token,
      },
    });
    if (response.ok) {
      setUsers(users.filter((user) => user.id !== userId));
    } else {
      const data = await response.json();
      toast.error(`delete user failed: ${data.error}`);
    }
  };

  return (
    <div className="p-4 h-full w-full bg-white space-y-4">
      <div className="flex justify-end items-center">
        <UserDialog
          user={null}
          onCreate={handleCreateUser}
          onEdit={handleUpdateUser}
        />
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-24">ID</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Role</TableHead>
            <TableHead>Operation</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {users.map((user) => (
            <TableRow key={user.id}>
              <TableCell className="font-medium">{user.id}</TableCell>
              <TableCell>{user.username}</TableCell>
              <TableCell
                className={cn(
                  "text-left",
                  user.role === "admin" ? "text-red-600" : "text-gray-500"
                )}
              >
                {user.role}
              </TableCell>
              <TableCell className="justify-end gap-2 flex">
                <UserDialog
                  user={user}
                  onCreate={handleUpdateUser}
                  onEdit={handleUpdateUser}
                />
                <DeleteDialog user={user} onDelete={handleDeleteUser} />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

const formSchema = z.object({
  username: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
  role: z.string(),
  password: z.string().min(6, {
    message: "Password must be at least 6 characters.",
  }),
});

const UserDialog = ({
  user,
  onCreate,
  onEdit,
}: {
  user: User | null;
  onCreate: (data: UserBean) => void;
  onEdit: (data: UserBean) => void;
}) => {
  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: user ? user.username : "",
      role: user ? user.role : "user",
      password: "",
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const userBean: UserBean = {
      ...data,
      id: user?.id || 0,
    };
    if (user === null) {
      onCreate(userBean);
    } else {
      onEdit(userBean);
    }

    setOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {user === null ? (
          <Button variant="default">Create User</Button>
        ) : (
          <Button variant="outline" size="icon">
            <Edit />
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {user === null
              ? "Create User"
              : `Edit User ${user.username}@${user.id}`}
          </DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input disabled={user !== null} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="role"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Role</FormLabel>
                  <FormControl>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <SelectTrigger className="w-[180px]">
                        <SelectValue placeholder="Select Role" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          <SelectItem value="admin" disabled>
                            Admin
                          </SelectItem>
                          <SelectItem value="user">User</SelectItem>
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password" {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
            <div className=" w-full flex gap-2 justify-end">
              <Button type="submit">Submit</Button>
              <DialogClose asChild>
                <Button variant="secondary">Cancel</Button>
              </DialogClose>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};

const DeleteDialog = ({
  user,
  onDelete,
}: {
  user: User;
  onDelete: (id: number) => void;
}) => {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="destructive" size="icon">
          <Trash />
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete user</DialogTitle>
        </DialogHeader>
        <DialogDescription>
          Are you sure you want to delete this user?
          <br />
          user id: {user.id}
          <br />
          username: {user.username}
        </DialogDescription>
        <DialogFooter>
          <Button variant="destructive" onClick={() => onDelete(user.id)}>
            Delete
          </Button>
          <DialogClose asChild>
            <Button variant="default">Cancel</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default User;
