import { BaseAccountButton } from "@/app/ui/(home)/base-account-buttons/base-account-button";
import { ArrowRightIcon, ArrowLeftIcon, UserCircleIcon } from '@heroicons/react/20/solid';
import Link from "next/link";

export default function AccountButtons() {
    const loggedin = false; // TODO: implement state

    return (
      <div className="flex flex-row space-x-2 px-2">
        {loggedin ? (
          <>
            <Link href="/login">
              <BaseAccountButton className="w-full">
                Log in <ArrowRightIcon className="ml-auto w-5 text-gray-50" />
              </BaseAccountButton>
            </Link>
            <Link href="/signup">
              <BaseAccountButton className="w-full">
                Sign up <ArrowRightIcon className="ml-auto w-5 text-gray-50" />
              </BaseAccountButton>
            </Link>
          </>
        ) : (
          <>
            <Link href="/logout">
            {/* TODO implement logout */}
              <BaseAccountButton className="w-full">
                Log out  <ArrowLeftIcon className="ml-auto w-5 text-gray-50" />
              </BaseAccountButton>
            </Link>
            <Link href="/account">
              <BaseAccountButton className="w-full">
                My Account  <UserCircleIcon className="ml-auto w-5 text-gray-50" />
              </BaseAccountButton>
            </Link>
          </>
        )}
      </div>
    );
}