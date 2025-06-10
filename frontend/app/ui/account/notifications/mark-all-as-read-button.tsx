'use client'

import { markAllNotificationsAsSeen } from "@/app/lib/api/notifications/requests";

export default function MarkAllAsReadButton() {

  const handleClick = async () => {
    try{
      await markAllNotificationsAsSeen();
      window.location.reload();
    } catch {
      alert("Failed to mark all notifications as seen")
    }
  }

    return (
        <button
            className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600"
            onClick={() => handleClick()}
        >
            Mark all as seen
        </button>
    )
}