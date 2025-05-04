import Link from "next/link";
import SearchField from "@/app/ui/[main]/search-field";
import LoginButtons from "@/app/ui/[main]/login-buttons";


export default function TopBar() {
    return (
      <div className="flex w-full flex-row">
        <Link
          className="mb-2 flex w-20 items-end justify-start rounded-br-lg bg-blue-600 p-2 md:w-60"
          href="/"
        >
          <div className="flex h-full w-full items-center justify-center rounded-md bg-blue-600 p-0 text-white">
            <p>HOME</p>
          </div>
        </Link>
        <div className="flex grow flex-col justify-between space-x-2 md:flex-row md:space-x-2 md:space-y-0">
          <SearchField />
          <LoginButtons />
          <div className="hidden w-auto h-full grow rounded-md bg-gray-50 md:block"></div>
        </div>
      </div>
    );
}