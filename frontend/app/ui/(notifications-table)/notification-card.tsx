'use client'

import { markNotificationAs } from "@/app/lib/api/notifications/requests";
import { Notification } from "@/app/lib/definitions/notification";
import { EnvelopeIcon, EnvelopeOpenIcon } from "@heroicons/react/20/solid";
import { useState } from "react";

export default function NotificationCard({
  notification,
  changeNumberOfUnread
}: {
  notification: Notification
  changeNumberOfUnread?: (count: number) => void
}) {
  const [isRead, setIsRead] = useState(notification.is_read);

  const handleReadChange = async (input?: boolean) => {
    const newValue = input !== undefined ? input : !isRead;
    setIsRead(newValue !== undefined ? newValue : !isRead);

    try {
      markNotificationAs(notification.id, newValue !== undefined ? newValue : !isRead)
      changeNumberOfUnread?.(newValue ? -1 : 1);
    }
    catch (error) {
      alert(`Failed to update notification status ${error}`);
      setIsRead(!isRead);
    }
  }

  const handleClick = () => {
    handleReadChange(true);
    window.location.href = `/offer/${notification.offer_id}`
  }


  return (
    <div
      className={`flex flex-col gap-1.5 p-3 bg-white rounded-lg shadow-md border-3 ${
          isRead ? 'border-gray-300' : 'border-gray-800'
        } relative hover:bg-gray-100 transition-colors duration-200 ease-in-out cursor-pointer`}
      onClick={handleClick}
      role="button"
      tabIndex={0}
    >
      <div
        className="absolute top-2 right-3"
        onClick={(e) => {
          e.stopPropagation();
          handleReadChange();
        }}
        role="button"
        tabIndex={0}
      >
        {isRead ? (
          <EnvelopeOpenIcon className="h-6 w-6 text-gray-500 cursor-pointer hover:text-gray-700" />
        ) : (
          <EnvelopeIcon className="h-6 w-6 text-gray-500 cursor-pointer hover:text-gray-700" />
        )}
      </div>
      <h3 className="text-lg font-semibold">{notification.title}</h3>
      <p className="text-gray-600">{notification.description}</p>
      <span className="text-sm text-gray-400">{notification.date}</span>
    </div>
  );
}