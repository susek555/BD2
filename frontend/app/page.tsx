"use client";

import { signIn, signOut, useSession } from "next-auth/react";
import Link from "next/link";

export default function Home() {
  const { data, status } = useSession();
  const isLoading = status === "loading";

  return (
    <div className="flex items-center justify-center min-h-screen p-8 font-[family-name:var(--font-geist-sans)]">
      <div className="flex flex-col items-center gap-4">
        {isLoading ? (
          <div className="flex flex-col items-center gap-4">
            <div className="w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
            <p className="text-gray-600">Loading...</p>
          </div>
        ) : status === "authenticated" ? (
          <>
            <h1 className="text-2xl font-bold">Welcome, {data.user?.username || "User"}!</h1>
            <p className="text-gray-600">{data.user?.email}</p>
            <Link href="/dashboard">
              <button className="px-4 py-2 mt-4 text-white bg-blue-600 rounded-md hover:bg-blue-700 transition-colors">
                Dashboard
              </button>
            </Link>
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
