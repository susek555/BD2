import { BellAlertIcon } from "@heroicons/react/20/solid"
import { BaseAccountButton } from "./base-account-buttons/base-account-button"

export default function NotificationsButton() {
    

    return(
        <>
            <BaseAccountButton className="w-12 flex items-center justify-center">
                <BellAlertIcon className="w-5 text-gray-50" />
            </BaseAccountButton>
        </>
    )
}