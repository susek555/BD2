'use client'

import { Notification } from "@/app/lib/definitions/notification";
import { EnvelopeIcon, EnvelopeOpenIcon } from "@heroicons/react/20/solid";
import { useState } from "react";

export default function NotificationCard({
  notification,
}: {
  notification: Notification
}) {
  const [isRead, setIsRead] = useState(notification.is_read);

  const handleReadChange = async () => {
    setIsRead(!isRead);

    // TODO send API request to update notification status
  }


  return (
    <div className={`flex flex-col gap-1.5 p-3 bg-white rounded-lg shadow-md border-3 ${isRead ? 'border-gray-300' : 'border-gray-800'} relative`}>
      <div className="absolute top-2 right-3" onClick={handleReadChange} role="button" tabIndex={0}>
        {isRead ?
          <EnvelopeOpenIcon className="h-6 w-6 text-gray-500 cursor-pointer" /> :
          <EnvelopeIcon className="h-6 w-6 text-gray-500 cursor-pointer" />
        }
      </div>
      <h3 className="text-lg font-semibold">{notification.title}</h3>
      <p className="text-gray-600">{notification.description}</p>
      <span className="text-sm text-gray-400">{notification.date}</span>
    </div>
  );
}