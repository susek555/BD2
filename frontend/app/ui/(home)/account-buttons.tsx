import { authConfig } from "@/app/api/auth/[...nextauth]/route";
import { BaseAccountButton } from "@/app/ui/(home)/base-account-buttons/base-account-button";
import { ClientLogoutButton } from "@/app/ui/(home)/client-logout-button";
import { ArrowRightIcon, UserCircleIcon } from '@heroicons/react/20/solid';
import { getServerSession } from "next-auth/next";
import Link from "next/link";

export default async function LoginButtons() {
  const session = await getServerSession(authConfig);
  const loggedIn = !!session;

  return (
    <div className="flex flex-row space-x-2 px-2">
      {loggedIn ? (
        <>
          <ClientLogoutButton />
          <Link href="/account">
            <BaseAccountButton className="w-full">
              My Account <UserCircleIcon className="ml-auto w-5 text-gray-50" />
            </BaseAccountButton>
          </Link>
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
