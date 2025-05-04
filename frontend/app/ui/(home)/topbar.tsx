import Link from "next/link";
import SearchField from "@/app/ui/(home)/search-field";
import LoginButtons from "@/app/ui/(home)/account-buttons";


export default function TopBar() {
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
        <div className="flex grow flex-col justify-between px-2 py-0 md:flex-row md:items-top md:space-x-2">
          <div className="flex grow justify-center">
            <SearchField placeholder="Search for cars..."/>
          </div>
          <div className="flex justify-end px-2">
            <LoginButtons />
          </div>
        </div>
      </div>
    );
}