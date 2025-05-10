"use client";

import { signIn, signOut, useSession } from "next-auth/react";
import Link from "next/link";
import { useState } from "react";

const API_URL = process.env.NEXT_PUBLIC_API_URL

export default function Test() {
  const { data, status } = useSession();
  const isLoading = status === "loading";
  const [isSigningOut, setIsSigningOut] = useState(false);

  const handleSignOut = async () => {
    if (!data?.user?.refreshToken) {
      console.warn("No refresh token available for logout");
      signOut();
      return;
    }

    setIsSigningOut(true);
    console.log("Handle sign out");

    try {
      const response = await fetch(`${API_URL}/logout`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${data.user.accessToken}`
        },
        body: JSON.stringify({ refresh_token: data.user.refreshToken }),
      });

      console.log("Logout response:", response.status);

      if (!response.ok) {
        console.error(`Logout failed: ${response.status} ${response.statusText}`);
      }
    } catch (error) {
      console.error('Log out error:', error);
    } finally {
      signOut();
      setIsSigningOut(false);
    }
  };



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
              onClick={handleSignOut}
              disabled={isSigningOut}
              className={`px-4 py-2 mt-4 text-white ${isSigningOut ? "bg-gray-400" : "bg-red-600 hover:bg-red-700"
                } rounded-md transition-colors flex items-center gap-2`}
            >
              {isSigningOut ? (
                <>
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  Signing Out...
                </>
              ) : (
                "Sign Out"
              )}
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
