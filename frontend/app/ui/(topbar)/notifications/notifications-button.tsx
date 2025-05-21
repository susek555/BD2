import { BellAlertIcon } from "@heroicons/react/20/solid"
import { BaseAccountButton } from "@/app/ui/(topbar)/base-account-buttons/base-account-button"

export default async function NotificationsButton() {
    // fetch notifications from the server
    const newNotifications = 7

    function handleNotificationsButtonClick() {

    }

    return (
        <>
            <BaseAccountButton className="w-12 flex items-center justify-center relative">
                <BellAlertIcon className="w-5 text-gray-50" />
                {newNotifications > 0 && (
                    <div className="absolute -top-1 -right-1 bg-red-500 text-white rounded-full w-5 h-5 flex items-center justify-center text-xs font-bold">
                        {newNotifications > 99 ? '99+' : newNotifications}
                    </div>
                )}
            </BaseAccountButton>
        </>
    )
}