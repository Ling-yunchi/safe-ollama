import { useState } from "react";
import { Button } from "./ui/button";
import { Check, Copy } from "lucide-react";
import copy from "copy-to-clipboard";

export function CopyButton({
  value,
  className,
  onCopy,
}: React.ComponentPropsWithoutRef<"button"> & {
  value: string;
  onCopy?: () => void;
}) {
  const [check, setCheck] = useState(false);

  return (
    <Button
      variant="ghost"
      size="icon"
      className={className}
      onClick={() => {
        try {
          navigator.clipboard.writeText(value);
        } catch {
          copy(value);
        }
        setCheck(true);
        onCopy?.();
        setTimeout(() => {
          setCheck(false);
        }, 1000);
      }}
    >
      {check ? <Check></Check> : <Copy></Copy>}
    </Button>
  );
}
