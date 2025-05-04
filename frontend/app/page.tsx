"use client";

import { signIn, signOut, useSession } from "next-auth/react";

export default function Home() {
  const session = useSession();
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <div className="flex flex-col items-center gap-4">
        {session.status === "authenticated" ? (
          <>
            <h1 className="text-2xl font-bold">Welcome, {session.data.user?.name || "User"}!</h1>
            <p className="text-gray-600">{session.data.user?.email}</p>
            <button
              onClick={() => signOut()}
              className="px-4 py-2 mt-4 text-white bg-red-600 rounded-md hover:bg-red-700 transition-colors"
            >
              Sign Out
            </button>
          </>
        ) : (
          <>
            <h1 className="text-2xl font-bold">Welcome, Guest!</h1>
            <p className="text-gray-600">Please sign in to continue</p>
            <button
              onClick={() => signIn()}
              className="px-4 py-2 mt-4 text-white bg-blue-600 rounded-md hover:bg-blue-700 transition-colors"
            >
              Sign In
            </button>
          </>
        )}
      </div>
    </div>
  );
}
