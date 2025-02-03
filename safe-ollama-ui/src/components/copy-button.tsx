import { useState } from "react";
import { Button } from "./ui/button";
import { Check, Copy } from "lucide-react";

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
        navigator.clipboard.writeText(value);
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
