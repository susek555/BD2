import LoginButtons from "@/app/ui/(topbar)/account-buttons";
import SearchField from "@/app/ui/(topbar)/search-field";
import Link from "next/link";
import { Suspense } from "react";
import { BaseAccountButton } from "@/app/ui/(topbar)/base-account-buttons/base-account-button";
import { PlusIcon } from "@heroicons/react/20/solid";
import { authConfig } from "@/app/lib/authConfig";
import { getServerSession } from "next-auth/next";
import HomeButton from "./home-button";


export async function TopBar() {
    const session = await getServerSession(authConfig);
    const loggedIn = !!session;

    return (
      <div className="flex w-full flex-row">
        <HomeButton />
        <div className="flex grow flex-col justify-between px-2 py-0 md:flex-row md:items-top md:space-x-2">
          <div className="flex grow justify-center">
            <Suspense>
              <SearchField placeholder="Search for cars..."/>
            </Suspense>
          </div>
          <div className="flex justify-end px-2">
            <div className="flex flex-row space-x-2">
              <LoginButtons />
              <Link href={loggedIn ? "/offer/add" : "/login"}>
                <BaseAccountButton className="hidden md:flex w-full">
                  Add Offer <PlusIcon className="ml-auto w-5 text-gray-50" />
                </BaseAccountButton>
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
}

export function OfferTopBar() {
    return (
      <div className="flex w-full flex-row">
        <Link
          className="mb-2 flex w-20 items-end justify-start rounded-br-lg bg-blue-600 md:w-80 md:h-12"
          href="/"
        >
            <div className="flex h-full w-full items-center justify-center rounded-md bg-blue-600 p-0 text-white">
            <p className="font-bold">HOME</p>
            </div>
        </Link>
        <div className="flex grow justify-end px-2">
          <div className="flex justify-end px-2">
            <LoginButtons />
          </div>
        </div>
      </div>
    );
}
