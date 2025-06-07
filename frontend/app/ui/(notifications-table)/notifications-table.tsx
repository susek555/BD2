import { Notification } from "@/app/lib/definitions/notification"
import NotificationCard from "./notification-card"

export default function NotificationsTable({
  notifications,
  changeNumberOfUnread,
  updateNotificationSeenStatus
}: {
  notifications: Notification[],
  changeNumberOfUnread: (count: number) => void,
  updateNotificationSeenStatus: (id: string, seen: boolean) => void
}) {
  return (
    <div className="flex flex-col gap-2 px-5">
      {notifications.map((notification) => (
        <NotificationCard key={notification.id} notification={notification} changeNumberOfUnread={changeNumberOfUnread} updateNotificationSeenStatus={updateNotificationSeenStatus}/>
      ))}
    </div>
  )
}