"use client"

import { BaseAccountButton } from "@/app/ui/(home)/base-account-buttons/base-account-button";
import { ArrowLeftIcon } from '@heroicons/react/20/solid';
import { signOut } from "next-auth/react";
import { useRouter } from "next/navigation";
import { useState } from "react";

export function ClientLogoutButton() {
  const [isSigningOut, setIsSigningOut] = useState(false);
  const router = useRouter();

  const handleSignOut = async () => {
    if (isSigningOut) return;
    setIsSigningOut(true)
    try {
      const res = await fetch("/api/auth/logout", { method: "POST" });
      if (!res.ok) throw new Error("Logout failed");
      await signOut({ redirect: true });
      router.refresh();
    } catch (err) {
      console.error(err);
    } finally {
      setIsSigningOut(false);
    }
  };

  return (
    <BaseAccountButton onClick={handleSignOut}>
      Log out <ArrowLeftIcon className="ml-auto w-5 text-gray-50" />
    </BaseAccountButton>
  );
}
