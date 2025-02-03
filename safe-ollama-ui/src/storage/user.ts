import { atom } from "jotai";
import { atomWithStorage } from "jotai/utils";

export const userAtom = atomWithStorage(
  "safe-ollama-user",
  {
    userId: -1,
    username: "",
    role: "",
    token: "",
    expires: "",
  },
  undefined,
  {
    getOnInit: true,
  }
);

export const userRole = atom((get) => get(userAtom).role);
export const userToken = atom((get) => get(userAtom).token);
