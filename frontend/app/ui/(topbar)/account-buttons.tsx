import { authConfig } from "@/app/lib/authConfig";
import { BaseAccountButton } from "@/app/ui/(topbar)/base-account-buttons/base-account-button";
import { ClientLogoutButton } from "@/app/ui/(topbar)/client-logout-button";
import { ArrowRightIcon, UserCircleIcon } from '@heroicons/react/20/solid';
import { getServerSession } from "next-auth/next";
import Link from "next/link";
import NotificationsHandler from "./notifications/notifications-handler";
import { Suspense } from "react";
import { NotificationsButtonSkeleton } from "../skeletons";

export default async function LoginButtons() {
  const session = await getServerSession(authConfig);
  const loggedIn = !!session;

  return (
    <div className="flex flex-row flex-grow space-x-2 pl-2">
      {loggedIn ? (
        <>
          <ClientLogoutButton />
          <Link href="/account">
            <BaseAccountButton className="w-full">
              My Account <UserCircleIcon className="ml-auto w-5 text-gray-50" />
            </BaseAccountButton>
          </Link>
          <Suspense fallback={<NotificationsButtonSkeleton />}>
            <NotificationsHandler />
          </Suspense>
        </>
      ) : (
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
      )}
    </div>
  );
}
