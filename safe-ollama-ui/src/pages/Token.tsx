import { useAtomValue } from "jotai";
import { useState, useEffect } from "react";
import { userToken } from "../storage/user";
import { toast } from "sonner";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { CopyButton } from "@/components/copy-button";
import {
  DialogHeader,
  DialogFooter,
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogTitle,
  DialogDescription,
  DialogClose,
} from "@/components/ui/dialog";
import { Trash } from "lucide-react";

interface Token {
  id: number;
  name: string;
  token: string;
  createdAt: number;
}

function Token() {
  const token = useAtomValue(userToken);
  const [tokens, setTokens] = useState<Token[]>([]);

  useEffect(() => {
    fetchTokens();
  }, []);

  const fetchTokens = async () => {
    const response = await fetch("/api/token/", {
      headers: {
        "Content-Type": "application/json",
        Token: token,
      },
    });
    if (!response.ok) {
      console.log("Error fetching tokens");
      return;
    }
    const data = await response.json();
    console.log(data);

    setTokens(data);
  };

  const handleCreateToken = async (data: { name: string }) => {
    const response = await fetch("/api/token/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Token: token,
      },
      body: JSON.stringify(data),
    });
    if (response.ok) {
      toast.success("Token created");
      await fetchTokens();
    }
  };

  const handleDeleteToken = async (tokenId: number) => {
    const response = await fetch(`/api/token/${tokenId}`, {
      method: "DELETE",
      headers: {
        Token: token,
      },
    });
    if (response.ok) {
      setTokens(tokens.filter((token) => token.id !== tokenId));
    }
  };

  return (
    <div className="p-4 h-full w-full bg-white space-y-4">
      <div className="flex justify-end items-center">
        <CreateDialog onCreate={handleCreateToken} />
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-24">ID</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Token</TableHead>
            <TableHead>Creat Time</TableHead>
            <TableHead>Operation</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {tokens.map((token) => (
            <TableRow key={token.id}>
              <TableCell className="font-medium">{token.id}</TableCell>
              <TableCell>{token.name}</TableCell>
              <TableCell className="flex justify-between items-center">
                {token.token} <CopyButton value={token.token} />
              </TableCell>
              <TableCell className="text-left">
                {new Date(token.createdAt).toLocaleString()}
              </TableCell>
              <TableCell className="text-right">
                <DeleteDialog token={token} onDelete={handleDeleteToken} />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

const CreateDialog = ({
  onCreate,
}: {
  onCreate: (data: { name: string }) => void;
}) => {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="default">New Token</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>New Token</DialogTitle>
        </DialogHeader>
        <Label htmlFor="name">Token Name</Label>
        <Input
          id="name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="w-64"
        />
        <DialogFooter>
          <Button
            type="submit"
            onClick={() => {
              if (!name) {
                toast.warning("Token name is required");
                return;
              }
              onCreate({ name: name });
              setName("");
              setOpen(false);
            }}
          >
            Create
          </Button>
          <DialogClose asChild>
            <Button variant="secondary">Cancel</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

const DeleteDialog = ({
  token,
  onDelete,
}: {
  token: Token;
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
          <DialogTitle>Delete Token</DialogTitle>
        </DialogHeader>
        <DialogDescription>
          Are you sure you want to delete this token?
          <br />
          token id: {token.id}
          <br />
          token name: {token.name}
        </DialogDescription>
        <DialogFooter>
          <Button variant="destructive" onClick={() => onDelete(token.id)}>
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

export default Token;
