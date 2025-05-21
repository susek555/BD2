import { Notification } from "@/app/lib/definitions/notification";

export default function NotificationCard({
  notification,
}: {
  notification: Notification
}) {
  return (
    <div className="flex flex-col gap-2 p-4 bg-white rounded-lg shadow-md">
      <h3 className="text-lg font-semibold">{notification.title}</h3>
      <p className="text-gray-600">{notification.description}</p>
      <span className="text-sm text-gray-400">{notification.date}</span>
    </div>
  );
}