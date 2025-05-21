import NotificationsButton from './notifications-button'
import { fetchNotifications } from '@/app/lib/data/notifications/data'

export default async function NotificationsHandler() {
    // fetch notifications from the server
    const { newNotifications, notifications } = await fetchNotifications()

    return (
        <>
            <NotificationsButton newNotifications={newNotifications} notifications={notifications}/>
        </>
    )
}