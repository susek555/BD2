'use client'

import { useSession } from "next-auth/react";

export default function Dashboard() {
  const { data: session, status } = useSession();

  if (status === "loading") {
    return <div>Loading...</div>;
  }

  if (!session) {
    return <div>No session</div>;
  }

  // Format the expiration date if it exists
  const expiryDate = session.user.accessTokenExpires ? new Date(session.user.accessTokenExpires) : null;

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold mb-4">Dashboard</h1>
      <div className="bg-gray-100 p-4 rounded">
        <p><strong>User:</strong> {session.user?.username || session.user?.email}</p>
        <p><strong>Email:</strong> {session.user?.email}</p>
        <p><strong>Session expires:</strong> {expiryDate?.toString()}</p>
        <p><strong>Time until expiry:</strong> {expiryDate ?
          `${Math.round((expiryDate.getTime() - Date.now()) / (1000 * 60))} minutes`
          : "Unknown"}
        </p>
      </div>
    </div>
  );
}
